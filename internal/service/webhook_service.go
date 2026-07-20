package service

import (
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

// Logger interface (placeholder)
type WebhookLogger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
}

// WebhookPublisher implements WebhookPublisherInterface
type WebhookPublisher struct {
    repo   repository.WebhookRepository
    logger WebhookLogger
    client *http.Client
}

// NewWebhookPublisher creates a new webhook publisher
func NewWebhookPublisher(repo repository.WebhookRepository, logger WebhookLogger) *WebhookPublisher {
    return &WebhookPublisher{
        repo:   repo,
        logger: logger,
        client: &http.Client{Timeout: 10 * time.Second},
    }
}

// PublishAsync implements WebhookPublisherInterface.PublishAsync
func (wp *WebhookPublisher) PublishAsync(ctx context.Context, eventType string, data interface{}) {
    go func() {
        wp.publish(ctx, eventType, data)
    }()
}

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
    payload, _ := json.Marshal(event)

    for _, sub := range subs {
        if !sub.IsActive {
            continue
        }
        delivery := &model.WebhookDelivery{
            SubscriptionID: sub.ID,
            EventID:        eventID,
            Payload:        payload,
            Status:         model.DeliveryPending,
            AttemptCount:   0,
        }
        _, _ = wp.repo.CreateDelivery(ctx, delivery)
    }

    wp.logger.Info("Webhook published", map[string]interface{}{
        "event_type":    eventType,
        "event_id":      eventID,
        "subscriptions": len(subs),
    })
}

// ProcessDeliveries runs retry logic for pending deliveries
func (wp *WebhookPublisher) ProcessDeliveries(ctx context.Context) error {
    deliveries, err := wp.repo.GetPendingDeliveries(ctx, 100)
    if err != nil {
        return fmt.Errorf("failed to get pending deliveries: %w", err)
    }

    for _, delivery := range deliveries {
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
        _ = wp.attemptDelivery(ctx, &delivery, sub)
    }
    return nil
}

func (wp *WebhookPublisher) attemptDelivery(ctx context.Context, delivery *model.WebhookDelivery, sub *model.WebhookSubscription) error {
    delivery.AttemptCount++
    signature := wp.signPayload(delivery.Payload, sub.SecretKey)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.WebhookURL, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Webhook-Event", delivery.EventID)
    req.Header.Set("X-Webhook-Signature", "sha256="+signature)
    req.Header.Set("X-Webhook-Timestamp", time.Now().Format(time.RFC3339))

    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    resp, err := wp.client.Do(req)
    if err != nil {
        wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryPending, nil, ptrString(err.Error()))
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliverySuccess, &resp.StatusCode, nil)
        return nil
    }

    if delivery.AttemptCount >= sub.MaxRetries {
        errMsg := fmt.Sprintf("HTTP %d after %d attempts", resp.StatusCode, delivery.AttemptCount)
        wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryDeadLetter, &resp.StatusCode, &errMsg)
        wp.logger.Warn("Webhook delivery failed (max retries)", map[string]interface{}{
            "delivery_id":     delivery.ID,
            "subscription_id": delivery.SubscriptionID,
            "attempts":        delivery.AttemptCount,
            "http_status":     resp.StatusCode,
        })
        return fmt.Errorf("max retries exceeded")
    }
    nextRetryTime := wp.calculateRetryTime(delivery.AttemptCount, sub.RetryInterval)
    errMsg := fmt.Sprintf("HTTP %d, retry at %v", resp.StatusCode, nextRetryTime)
    wp.repo.UpdateDeliveryStatus(ctx, delivery.ID, model.DeliveryPending, &resp.StatusCode, &errMsg)
    return fmt.Errorf("failed with status %d, will retry", resp.StatusCode)
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

func ptrString(s string) *string {
    return &s
}