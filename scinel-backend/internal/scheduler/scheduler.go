// Package scheduler provides a minimal, dependency-free background job
// runner. It exists because two service methods had no active call site
// (TD-N6 / TD-N13): ForecastService.GenerateForecastBatch and
// WebhookPublisher.ProcessDeliveries.
//
// Design constraints (kept deliberately boring, R4 - zero risk):
//   - Native stdlib only. No cron library, no external scheduler.
//   - One goroutine per registered job, driven by a resettable timer so a
//     job can have an initial delay (e.g. "next 02:00 UTC") without burning
//     CPU in a sleep loop.
//   - A panicking or failing job never takes down the process: the panic is
//     recovered and the job loop keeps its own schedule.
//   - Stop cancels the derived context and blocks until every job goroutine
//     has returned, so shutdown cannot leak goroutines.
package scheduler

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// Logger is the shared logging contract used by the scheduler. It mirrors
// service.ServiceLogger (internal/service/forecast_service.go) so that
// *logger.Logger (pkg/logger) and existing test mocks both satisfy it without
// this package importing the service layer.
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// Job names of the background jobs wired in cmd/api. They are exported so the
// names have a single source of truth: they appear in structured logs, so a
// silent rename would break log queries.
const (
	// JobNameForecastBatch is the daily forecast generation batch
	// (ForecastService.GenerateForecastBatch). Run at 02:00 UTC.
	JobNameForecastBatch = "forecast-batch"
	// JobNameWebhookWorker polls pending webhook deliveries and drives
	// retry / dead-letter handling (WebhookPublisher.ProcessDeliveries).
	JobNameWebhookWorker = "webhook-worker"
)

// Job is a single scheduled unit of work.
//
// Interval must be > 0. InitialDelay is optional (0 means "run as soon as
// the scheduler starts") and is what lets a daily job be aligned to a wall
// clock hour, e.g. InitialDelay: scheduler.NextDailyAt(2).
type Job struct {
	Name         string
	Interval     time.Duration
	InitialDelay time.Duration
	Run          func(context.Context) error
}

// Scheduler owns a set of jobs and runs each of them on its own goroutine.
// The zero value is not usable; construct one with New.
type Scheduler struct {
	logger Logger

	mu      sync.Mutex
	jobs    []Job
	started bool
	cancel  context.CancelFunc

	wg sync.WaitGroup
}

// New creates an empty Scheduler. log may be nil, in which case logging is
// dropped (useful for tests that do not assert on log output).
func New(log Logger) *Scheduler {
	return &Scheduler{logger: log}
}

// Register adds a job. Invalid jobs (blank name, nil Run, non-positive
// Interval) are rejected with an error log instead of being started, because
// a zero Interval would otherwise turn into a hot loop.
func (s *Scheduler) Register(job Job) {
	switch {
	case job.Name == "":
		s.error("Scheduler rejected job: empty name", "interval", job.Interval)
		return
	case job.Run == nil:
		s.error("Scheduler rejected job: nil Run function", "job", job.Name)
		return
	case job.Interval <= 0:
		s.error("Scheduler rejected job: non-positive interval",
			"job", job.Name, "interval", job.Interval)
		return
	}
	if job.InitialDelay < 0 {
		s.warn("Scheduler clamped negative initial delay to 0",
			"job", job.Name, "initial_delay", job.InitialDelay)
		job.InitialDelay = 0
	}

	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		s.error("Scheduler rejected job registered after Start", "job", job.Name)
		return
	}
	s.jobs = append(s.jobs, job)
	count := len(s.jobs)
	s.mu.Unlock()

	s.info("Job registered", "job", job.Name,
		"interval", job.Interval, "initial_delay", job.InitialDelay,
		"registered_jobs", count)
}

// Jobs returns a copy of the registered jobs (callers cannot mutate the
// scheduler's internal slice).
func (s *Scheduler) Jobs() []Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Job, len(s.jobs))
	copy(out, s.jobs)
	return out
}

// Start launches one goroutine per registered job. It is safe to call more
// than once; subsequent calls are logged and ignored so jobs are never
// duplicated.
//
// The job contexts are children of ctx and are cancelled by Stop, so a
// caller can also abort everything by cancelling ctx.
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		s.warn("Scheduler Start called more than once, ignoring")
		return
	}
	s.started = true
	jobs := make([]Job, len(s.jobs))
	copy(jobs, s.jobs)
	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.mu.Unlock()

	for i := range jobs {
		s.wg.Add(1)
		go s.run(runCtx, jobs[i])
	}

	s.info("Scheduler started", "jobs", len(jobs))
}

// Stop cancels every job context and blocks until all job goroutines have
// returned. It is safe to call repeatedly and safe to call when Start was
// never called.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	started := s.started
	cancel := s.cancel
	s.mu.Unlock()

	if !started {
		s.info("Scheduler stop requested while not started, nothing to do")
		return
	}

	if cancel != nil {
		cancel()
	}
	s.wg.Wait()
	s.info("Scheduler stopped")
}

// run drives a single job until its context is done. A resettable timer is
// used instead of a Ticker so the first execution can be delayed.
func (s *Scheduler) run(ctx context.Context, job Job) {
	defer s.wg.Done()

	timer := time.NewTimer(job.InitialDelay)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			s.info("Job stopped", "job", job.Name)
			return
		case <-timer.C:
			s.execute(ctx, job)
			timer.Reset(job.Interval)
		}
	}
}

// execute runs the job body with panic recovery. A panicking job is logged
// (message + stack) and its schedule continues; an erroring job is logged
// and likewise keeps its schedule.
func (s *Scheduler) execute(ctx context.Context, job Job) {
	defer func() {
		if r := recover(); r != nil {
			s.error("Job panicked, recovered by scheduler",
				"job", job.Name,
				"panic", fmt.Sprintf("%v", r),
				"stack", string(debug.Stack()))
		}
	}()

	start := time.Now()
	err := job.Run(ctx)
	if err != nil {
		// A cancelled context during shutdown is expected, not a fault.
		if ctx.Err() != nil {
			s.info("Job aborted by shutdown", "job", job.Name, "error", err.Error())
			return
		}
		s.error("Job failed", "job", job.Name, "error", err.Error(),
			"duration_ms", time.Since(start).Milliseconds())
		return
	}
	s.info("Job completed", "job", job.Name,
		"duration_ms", time.Since(start).Milliseconds())
}

// NextDailyAt returns the duration from now until the next occurrence of
// hourUTC:00:00 UTC. The result is always in (0, 24h], so it is safe to use
// as a Job.InitialDelay for a daily job.
//
// Clock-hour alignment is stable across restarts because UTC has no DST.
func NextDailyAt(hourUTC int) time.Duration {
	if hourUTC < 0 || hourUTC > 23 {
		// Guard against a misconfigured constant rather than panicking on
		// time.Date normalisation.
		hourUTC = ((hourUTC % 24) + 24) % 24
	}
	now := time.Now().UTC()
	next := time.Date(now.Year(), now.Month(), now.Day(), hourUTC, 0, 0, 0, time.UTC)
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	return next.Sub(now)
}

func (s *Scheduler) info(msg string, fields ...interface{}) {
	if s.logger != nil {
		s.logger.Info(msg, fields...)
	}
}

func (s *Scheduler) warn(msg string, fields ...interface{}) {
	if s.logger != nil {
		s.logger.Warn(msg, fields...)
	}
}

func (s *Scheduler) error(msg string, fields ...interface{}) {
	if s.logger != nil {
		s.logger.Error(msg, fields...)
	}
}
