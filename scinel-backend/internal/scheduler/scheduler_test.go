package scheduler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// testLogger records log calls including their fields. It is mutex-guarded
// because the scheduler logs from every job goroutine concurrently and the
// suite runs under -race.
type testLogger struct {
	mu     sync.Mutex
	infos  []string
	warns  []string
	errors []string
}

func (l *testLogger) record(dst *[]string, msg string, fields []interface{}) {
	entry := msg
	for i := 0; i+1 < len(fields); i += 2 {
		entry += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	*dst = append(*dst, entry)
}

func (l *testLogger) Info(msg string, fields ...interface{}) {
	l.record(&l.infos, msg, fields)
}

func (l *testLogger) Warn(msg string, fields ...interface{}) {
	l.record(&l.warns, msg, fields)
}

func (l *testLogger) Error(msg string, fields ...interface{}) {
	l.record(&l.errors, msg, fields)
}

func (l *testLogger) snapshot() (infos, warns, errs []string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.infos...),
		append([]string(nil), l.warns...),
		append([]string(nil), l.errors...)
}

func (l *testLogger) contains(list []string, want string) bool {
	for _, m := range list {
		if strings.Contains(m, want) {
			return true
		}
	}
	return false
}

func (l *testLogger) hasError(want string) bool {
	_, _, errs := l.snapshot()
	return l.contains(errs, want)
}

func (l *testLogger) hasInfo(want string) bool {
	infos, _, _ := l.snapshot()
	return l.contains(infos, want)
}

func (l *testLogger) hasWarn(want string) bool {
	_, warns, _ := l.snapshot()
	return l.contains(warns, want)
}

// await fails the test if ch does not receive within timeout. Channel
// synchronisation is used everywhere instead of time.Sleep so a slow CI box
// cannot turn into a flaky assertion.
func await(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
	}
}

// shortJobInterval keeps the suite fast while still exercising the
// reset-timer path more than once.
const shortJobInterval = 15 * time.Millisecond

func TestNewSchedulerStartsEmpty(t *testing.T) {
	s := New(&testLogger{})
	if s == nil {
		t.Fatal("New returned nil")
	}
	if got := len(s.Jobs()); got != 0 {
		t.Fatalf("expected 0 registered jobs, got %d", got)
	}
}

func TestNewSchedulerAcceptsNilLogger(t *testing.T) {
	s := New(nil)
	s.Register(Job{Name: "job", Interval: time.Hour, Run: func(context.Context) error { return nil }})
	if got := len(s.Jobs()); got != 1 {
		t.Fatalf("expected 1 registered job, got %d", got)
	}
	// Start/Stop must not panic on the nil logger path.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	s.Stop()
}

func TestRegisterRejectsInvalidJobs(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	s.Register(Job{Name: "", Interval: time.Hour, Run: func(context.Context) error { return nil }})
	s.Register(Job{Name: "no-run", Interval: time.Hour, Run: nil})
	s.Register(Job{Name: "zero-interval", Interval: 0, Run: func(context.Context) error { return nil }})
	s.Register(Job{Name: "negative-interval", Interval: -time.Second, Run: func(context.Context) error { return nil }})

	if got := len(s.Jobs()); got != 0 {
		t.Fatalf("expected all invalid jobs to be rejected, got %d registered", got)
	}
	for _, want := range []string{
		"empty name",
		"nil Run function",
		"non-positive interval",
	} {
		if !log.hasError(want) {
			t.Errorf("expected rejection log containing %q, got errors=%v", want, mustErrors(log))
		}
	}
}

func TestRegisterClampsNegativeInitialDelay(t *testing.T) {
	log := &testLogger{}
	s := New(log)
	s.Register(Job{
		Name:         "clamped",
		Interval:     time.Hour,
		InitialDelay: -5 * time.Minute,
		Run:          func(context.Context) error { return nil },
	})

	jobs := s.Jobs()
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].InitialDelay != 0 {
		t.Errorf("expected initial delay clamped to 0, got %v", jobs[0].InitialDelay)
	}
	if !log.hasWarn("negative initial delay") {
		t.Errorf("expected warn about negative initial delay, got warns=%v", mustWarns(log))
	}
}

func TestRegisterAfterStartIsRejected(t *testing.T) {
	log := &testLogger{}
	s := New(log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Start(ctx)
	s.Register(Job{Name: "late", Interval: time.Hour, Run: func(context.Context) error { return nil }})

	if got := len(s.Jobs()); got != 0 {
		t.Fatalf("expected job registered after Start to be rejected, got %d", got)
	}
	if !log.hasError("registered after Start") {
		t.Errorf("expected rejection log for late registration, got errors=%v", mustErrors(log))
	}
	s.Stop()
}

func TestJobsReturnsDefensiveCopy(t *testing.T) {
	s := New(&testLogger{})
	s.Register(Job{Name: "a", Interval: time.Hour, Run: func(context.Context) error { return nil }})

	jobs := s.Jobs()
	jobs[0].Name = "mutated"

	if s.Jobs()[0].Name != "a" {
		t.Fatal("Jobs() must return a copy; internal slice was mutated")
	}
}

func TestStartRunsJobImmediatelyAndOnInterval(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	var calls int32
	first := make(chan struct{}, 8)
	s.Register(Job{
		Name:     "ticker-job",
		Interval: shortJobInterval,
		Run: func(context.Context) error {
			atomic.AddInt32(&calls, 1)
			select {
			case first <- struct{}{}:
			default:
			}
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	defer s.Stop()

	// No InitialDelay => must run at once, then repeat on the interval.
	await(t, first, "first immediate run")
	await(t, first, "second run on interval")
	await(t, first, "third run on interval")

	if n := atomic.LoadInt32(&calls); n < 3 {
		t.Errorf("expected at least 3 runs, got %d", n)
	}
	if !log.hasInfo("Job completed") {
		t.Errorf("expected Job completed log, got infos=%v", mustInfos(log))
	}
	if !log.hasInfo("Scheduler started") {
		t.Error("expected Scheduler started log")
	}
}

func TestInitialDelayDefersFirstRun(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	var calls int32
	ran := make(chan struct{}, 4)
	s.Register(Job{
		Name:         "delayed-job",
		Interval:     shortJobInterval,
		InitialDelay: 60 * time.Millisecond,
		Run: func(context.Context) error {
			atomic.AddInt32(&calls, 1)
			select {
			case ran <- struct{}{}:
			default:
			}
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	defer s.Stop()

	// Nothing must run before InitialDelay elapses: a non-blocking read
	// right after Start is the assertion (no sleep, no flake).
	if n := atomic.LoadInt32(&calls); n != 0 {
		t.Fatalf("job ran %d time(s) before its initial delay elapsed", n)
	}
	select {
	case <-ran:
		t.Fatal("job ran before initial delay elapsed")
	default:
	}

	await(t, ran, "run after initial delay")
}

func TestJobErrorIsLoggedAndScheduleContinues(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	var calls int32
	attempts := make(chan struct{}, 8)
	s.Register(Job{
		Name:     "flaky-job",
		Interval: shortJobInterval,
		Run: func(context.Context) error {
			n := atomic.AddInt32(&calls, 1)
			select {
			case attempts <- struct{}{}:
			default:
			}
			if n == 1 {
				return errors.New("boom")
			}
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	defer s.Stop()

	// First attempt fails, second attempt must still happen: a failing job
	// must not stop its own schedule.
	await(t, attempts, "failing attempt")
	await(t, attempts, "retry attempt after error")

	if !log.hasError("Job failed") {
		t.Errorf("expected Job failed log, got errors=%v", mustErrors(log))
	}
	if !log.hasError("boom") {
		t.Error("expected the underlying error to be logged")
	}
	if n := atomic.LoadInt32(&calls); n < 2 {
		t.Errorf("expected scheduler to keep running the job, got %d call(s)", n)
	}
}

func TestPanicIsRecoveredAndOtherJobsSurvive(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	var healthyCalls int32
	healthy := make(chan struct{}, 16)
	s.Register(Job{
		Name:     "panicking-job",
		Interval: shortJobInterval,
		Run: func(context.Context) error {
			panic("scheduler must survive this")
		},
	})
	s.Register(Job{
		Name:     "healthy-job",
		Interval: shortJobInterval,
		Run: func(context.Context) error {
			atomic.AddInt32(&healthyCalls, 1)
			select {
			case healthy <- struct{}{}:
			default:
			}
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	defer s.Stop()

	// The panicking job keeps panicking on every tick; the healthy job must
	// keep running. If recovery failed, the whole test binary would die.
	await(t, healthy, "healthy job run alongside panicking job")
	await(t, healthy, "second healthy job run after panics")
	await(t, healthy, "third healthy job run after panics")

	if n := atomic.LoadInt32(&healthyCalls); n < 3 {
		t.Errorf("expected healthy job to keep running, got %d call(s)", n)
	}
	if !log.hasError("Job panicked, recovered by scheduler") {
		t.Errorf("expected panic recovery log, got errors=%v", mustErrors(log))
	}
	if !log.hasError("scheduler must survive this") {
		t.Error("expected the panic value to be logged")
	}
}

func TestStopWaitsForJobsAndIsIdempotent(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	var calls int32
	finished := make(chan struct{}, 32)
	s.Register(Job{
		Name:     "stop-job",
		Interval: time.Hour,
		Run: func(context.Context) error {
			atomic.AddInt32(&calls, 1)
			finished <- struct{}{}
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	await(t, finished, "job run after Start")

	s.Stop()
	// Stop returned => the goroutine observed cancellation and exited, so
	// the job cannot fire again. Confirm no further run is possible.
	select {
	case <-finished:
		t.Fatal("job ran after Stop returned")
	case <-time.After(50 * time.Millisecond):
	}

	// Repeated Stop must be a no-op, not a panic or a hang.
	s.Stop()
	s.Stop()

	if !log.hasInfo("Scheduler stopped") {
		t.Errorf("expected Scheduler stopped log, got infos=%v", mustInfos(log))
	}
	if !log.hasInfo("Job stopped") {
		t.Errorf("expected per-job stop log, got infos=%v", mustInfos(log))
	}
}

func TestStopWithoutStartIsSafe(t *testing.T) {
	log := &testLogger{}
	s := New(log)
	s.Register(Job{Name: "never-started", Interval: time.Hour, Run: func(context.Context) error { return nil }})

	s.Stop()

	if !log.hasInfo("not started") {
		t.Errorf("expected not-started log, got infos=%v", mustInfos(log))
	}
}

func TestContextCancelStopsAllJobs(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	var calls int32
	ran := make(chan struct{}, 64)
	for _, name := range []string{"job-a", "job-b", "job-c"} {
		s.Register(Job{
			Name:     name,
			Interval: time.Hour,
			Run: func(context.Context) error {
				atomic.AddInt32(&calls, 1)
				select {
				case ran <- struct{}{}:
				default:
				}
				return nil
			},
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.Start(ctx)

	await(t, ran, "first job run")
	await(t, ran, "second job run")
	await(t, ran, "third job run")

	cancel()
	// Stop must return promptly now that ctx is cancelled: it is waiting on
	// the job goroutines only.
	stopped := make(chan struct{})
	go func() {
		s.Stop()
		close(stopped)
	}()
	await(t, stopped, "Stop to return after ctx cancel")

	before := atomic.LoadInt32(&calls)
	select {
	case <-ran:
		t.Fatal("a job ran after ctx cancellation")
	case <-time.After(50 * time.Millisecond):
	}
	if after := atomic.LoadInt32(&calls); after != before {
		t.Errorf("job call count changed after shutdown: %d -> %d", before, after)
	}
}

func TestLongRunningJobReceivesCancelledContext(t *testing.T) {
	log := &testLogger{}
	s := New(log)

	entered := make(chan struct{})
	observed := make(chan struct{}, 1)
	s.Register(Job{
		Name:     "long-job",
		Interval: time.Hour,
		Run: func(ctx context.Context) error {
			close(entered)
			<-ctx.Done() // blocks until Stop cancels the job context
			observed <- struct{}{}
			return ctx.Err()
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	await(t, entered, "long job to start")

	s.Stop()

	select {
	case <-observed:
	case <-time.After(5 * time.Second):
		t.Fatal("long-running job never observed context cancellation")
	}
	// A job aborted by shutdown must not be reported as a failure.
	if log.hasError("Job failed") {
		t.Errorf("shutdown abort must not log Job failed, got errors=%v", mustErrors(log))
	}
	if !log.hasInfo("Job aborted by shutdown") {
		t.Errorf("expected shutdown abort log, got infos=%v", mustInfos(log))
	}
}

func TestNextDailyAtLandsOnRequestedUTCHour(t *testing.T) {
	for hour := 0; hour < 24; hour++ {
		before := time.Now().UTC()
		d := NextDailyAt(hour)
		after := time.Now().UTC()

		if d <= 0 {
			t.Fatalf("hour %d: expected positive delay, got %v", hour, d)
		}
		if d > 24*time.Hour {
			t.Fatalf("hour %d: expected delay <= 24h, got %v", hour, d)
		}

		// The target instant must be hour:00:00 UTC. NextDailyAt samples its
		// own "now" internally, so reconstructing the target as before+d can
		// land a few hundred nanoseconds on either side of the boundary (and
		// on the neighbouring hour when the two samples straddle midnight).
		// Accept a candidate hour boundary within a 1s tolerance.
		target := before.Add(d)
		base := target.Truncate(time.Hour)
		aligned := false
		for _, cand := range []time.Time{base, base.Add(time.Hour)} {
			if cand.Hour() != hour {
				continue
			}
			if drift := target.Sub(cand); drift > -time.Second && drift < time.Second {
				aligned = true
				break
			}
		}
		if !aligned {
			t.Fatalf("hour %d: reconstructed target %s is not within 1s of %02d:00:00 UTC",
				hour, target, hour)
		}
		if target.Before(before) || target.After(after.Add(24*time.Hour)) {
			t.Fatalf("hour %d: target %s outside expected window", hour, target)
		}
	}
}

func TestNextDailyAtNeverReturnsNonPositiveDelay(t *testing.T) {
	// Regression guard for the forecast-batch schedule: a 24h InitialDelay of
	// 0 or negative would make the job fire in a hot loop. NextDailyAt reads
	// the wall clock internally, so this asserts the invariant over the real
	// clock rather than pretending to inject a fixed "now" at the boundary.
	for _, hour := range []int{0, 2, 12, 23} {
		d := NextDailyAt(hour)
		if d <= 0 {
			t.Fatalf("hour %d: expected positive delay, got %v", hour, d)
		}
		if d > 24*time.Hour {
			t.Fatalf("hour %d: expected delay <= 24h, got %v", hour, d)
		}
	}
}

func TestNextDailyAtNormalisesOutOfRangeHour(t *testing.T) {
	for _, hour := range []int{-1, 24, 25, 49} {
		if d := NextDailyAt(hour); d <= 0 || d > 24*time.Hour {
			t.Fatalf("hour %d: expected delay in (0,24h], got %v", hour, d)
		}
	}
}

func TestJobNameConstants(t *testing.T) {
	// These strings end up in structured logs, so a rename would silently
	// break log queries. Pin them.
	if JobNameForecastBatch != "forecast-batch" {
		t.Errorf("JobNameForecastBatch = %q, want %q", JobNameForecastBatch, "forecast-batch")
	}
	if JobNameWebhookWorker != "webhook-worker" {
		t.Errorf("JobNameWebhookWorker = %q, want %q", JobNameWebhookWorker, "webhook-worker")
	}
}

func mustInfos(l *testLogger) []string {
	infos, _, _ := l.snapshot()
	return infos
}

func mustWarns(l *testLogger) []string {
	_, warns, _ := l.snapshot()
	return warns
}

func mustErrors(l *testLogger) []string {
	_, _, errs := l.snapshot()
	return errs
}
