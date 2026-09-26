package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/LittleMardiah/scinel/internal/repository"
)

// WebhookPublisher implements WebhookPublisherInterface and handles event
// publishing, retry with exponential backoff, and dead-letter queueing.
type WebhookPublisher struct {
    repo   repository.WebhookRepository
    logger ServiceLogger
    client *http.Client
    auditSvc *AuditService
    // callbackFn, when set, is invoked instead of a real HTTP call. Used by tests.
    callbackFn func(req *http.Request) (*http.Response, error)
}

// NewWebhookPublisher creates a new webhook publisher.
//
// auditSvc is variadic so existing call sites keep compiling; with no
// AuditService the WEBHOOK_PUBLISHED audit entry is simply not written.
func NewWebhookPublisher(repo repository.WebhookRepository, logger ServiceLogger, auditSvc ...*AuditService) *WebhookPublisher {
    wp := &WebhookPublisher{
        repo:   repo,
        logger: logger,
        client: &http.Client{Timeout: 10 * time.Second},
    }
    for _, a := range auditSvc {
        if a != nil {
            wp.auditSvc = a
        }
    }
    return wp
}

// PublishAsync implements WebhookPublisherInterface.PublishAsync.
func (wp *WebhookPublisher) PublishAsync(ctx context.Context, eventType string, data interface{}) {
	go func() {
		wp.publish(ctx, eventType, data)
	}()
}

// publish queues a delivery for every active subscription of eventType.
func (wp *WebhookPublisher) publish(ctx context.Context, eventType string, data interface{}) {
	subs, err := wp.repo.GetSubscriptionsByEventType(ctx, eventType)
	if err != nil {
		wp.logger.Error("Failed to get subscriptions", map[string]interface{}{
			"event_type": eventType,
			"error":      err.Error(),
		})
		return
	}

	if len(subs) == 0 {
		return
	}

	eventID := uuid.New().String()
	event := model.WebhookEvent{
		EventID:   eventID,
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		wp.logger.Error("Failed to marshal webhook payload", map[string]interface{}{
			"event_type": eventType,
			"error":      err.Error(),
		})
		return
	}

	queued := 0
	for _, sub := range subs {
		if !sub.IsActive {
			continue
		}
		now := time.Now().UTC()
		delivery := &model.WebhookDelivery{
			ID:               uuid.New(),
			SubscriptionID:   sub.ID,
			EventID:          eventID,
			Payload:          payload,
			Status:           model.DeliveryPending,
			AttemptCount:     0,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if _, err := wp.repo.CreateDelivery(ctx, delivery); err != nil {
			wp.logger.Error("Failed to enqueue webhook delivery", map[string]interface{}{
				"event_type": eventType, "subscription_id": sub.ID, "error": err.Error(),
			})
			continue
		}
		queued++
	}

    wp.logger.Info("Webhook published", map[string]interface{}{
        "event_type":    eventType,
        "event_id":      eventID,
        "subscriptions": queued,
    })

    // Audit: WEBHOOK_PUBLISHED. System-initiated (no acting user), so user_id
    // stays NULL. Non-fatal by construction.
    if raw, err := json.Marshal(map[string]interface{}{
        "event_id":      eventID,
        "event_type":    eventType,
        "subscriptions": queued,
    }); err == nil {
        wp.auditSvc.Log(ctx, &model.AuditLog{
            Action:       model.AuditActionWebhookPublished,
            ResourceType: nullableString(model.AuditResourceWebhookEvent),
            ResourceID:   nullableString(eventID),
            Details:      raw,
        })
    }
}

// ProcessDeliveries runs retry logic for pending deliveries whose retry time
// has arrived.
func (wp *WebhookPublisher) ProcessDeliveries(ctx context.Context) error {
	deliveries, err := wp.repo.GetPendingDeliveries(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to get pending deliveries: %w", err)
	}

	for i := range deliveries {
		delivery := deliveries[i]
		if delivery.ScheduledRetryAt != nil && delivery.ScheduledRetryAt.After(time.Now()) {
			continue
		}
		sub, err := wp.repo.GetSubscriptionByID(ctx, delivery.SubscriptionID)
		if err != nil || sub == nil {
			wp.logger.Error("Failed to get subscription for delivery", map[string]interface{}{
				"delivery_id":     delivery.ID,
				"subscription_id": delivery.SubscriptionID,
				"error":           err,
			})
			continue
		}
		if err := wp.attemptDelivery(ctx, &delivery, sub); err != nil {
			wp.logger.Warn("Webhook delivery attempt failed", map[string]interface{}{
				"delivery_id": delivery.ID, "subscription_id": sub.ID,
				"attempt": delivery.AttemptCount, "error": err.Error(),
			})
		}
	}
	return nil
}

// attemptDelivery POSTs the payload to the subscriber with HMAC signature and
// applies exponential-backoff retry / dead-letter logic.
func (wp *WebhookPublisher) attemptDelivery(ctx context.Context, delivery *model.WebhookDelivery, sub *model.WebhookSubscription) error {
	delivery.AttemptCount++
	signature := wp.signPayload(delivery.Payload, sub.SecretKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.WebhookURL, bytes.NewReader(delivery.Payload))
	if err != nil {
		msg := "failed to create request: " + err.Error()
		return wp.handleSendFailure(ctx, delivery, sub, 0, msg)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", delivery.EventID)
	req.Header.Set("X-Webhook-Signature", "sha256="+signature)
	req.Header.Set("X-Webhook-Timestamp", time.Now().Format(time.RFC3339))

	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var resp *http.Response
	if wp.callbackFn != nil {
		resp, err = wp.callbackFn(req.Clone(callCtx))
	} else {
		resp, err = wp.client.Do(req.WithContext(callCtx))
	}

	wp.logger.Info("Webhook delivery attempt", map[string]interface{}{
		"delivery_id": delivery.ID, "subscription_id": sub.ID,
		"attempt": delivery.AttemptCount, "event_id": delivery.EventID,
		"url": sub.WebhookURL,
	})

	if err != nil {
		msg := "network error: " + err.Error()
		return wp.handleSendFailure(ctx, delivery, sub, 0, msg)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliverySuccess, &resp.StatusCode, nil); err != nil {
			return fmt.Errorf("failed to update delivery status: %w", err)
		}
		wp.logger.Info("Webhook delivered", map[string]interface{}{
			"delivery_id": delivery.ID, "http_status": resp.StatusCode, "attempt": delivery.AttemptCount,
		})
		return nil
	}

	msg := fmt.Sprintf("HTTP %d", resp.StatusCode)
	return wp.handleSendFailure(ctx, delivery, sub, resp.StatusCode, msg)
}

// handleSendFailure applies dead-letter or schedule-retry semantics after a
// failed delivery attempt.
func (wp *WebhookPublisher) handleSendFailure(ctx context.Context, delivery *model.WebhookDelivery, sub *model.WebhookSubscription, httpStatus int, msg string) error {
	if delivery.AttemptCount >= sub.MaxRetries {
		errMsg := fmt.Sprintf("%s after %d attempts", msg, delivery.AttemptCount)
		var statusPtr *int
		if httpStatus != 0 {
			statusPtr = &httpStatus
		}
		if err := wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryDeadLetter, statusPtr, &errMsg); err != nil {
			return fmt.Errorf("failed to move delivery to dead-letter queue: %w", err)
		}
		wp.logger.Warn("Webhook delivery moved to dead-letter queue", map[string]interface{}{
			"delivery_id": delivery.ID, "subscription_id": delivery.SubscriptionID,
			"attempts": delivery.AttemptCount, "error": msg,
		})
		return fmt.Errorf("max retries exceeded: %s", errMsg)
	}

	nextRetryTime := wp.calculateRetryTime(delivery.AttemptCount, sub.RetryInterval)
	var statusPtr *int
	if httpStatus != 0 {
		statusPtr = &httpStatus
	}
	if err := wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryPending, statusPtr, &msg); err != nil {
		return fmt.Errorf("failed to schedule retry: %w", err)
	}
	wp.logger.Warn("Webhook delivery scheduled for retry", map[string]interface{}{
		"delivery_id": delivery.ID, "attempt": delivery.AttemptCount,
		"retry_at": nextRetryTime.Format(time.RFC3339), "error": msg,
	})
	return fmt.Errorf("failed with %s, retry at %v", msg, nextRetryTime.Format(time.RFC3339))
}

func (wp *WebhookPublisher) signPayload(payload []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (wp *WebhookPublisher) calculateRetryTime(attemptCount, retryInterval int) time.Time {
	backoffSeconds := 1 << uint(attemptCount-1)
	if backoffSeconds < retryInterval {
		backoffSeconds = retryInterval
	}
	if backoffSeconds > 3600 {
		backoffSeconds = 3600
	}
	return time.Now().Add(time.Duration(backoffSeconds) * time.Second)
}
