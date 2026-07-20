package model

import (
    "time"

    "github.com/google/uuid"
)

// WebhookSubscription represents a webhook event subscription
type WebhookSubscription struct {
    ID             uuid.UUID `db:"id"`
    UserID         uuid.UUID `db:"user_id"`
    EventType      string    `db:"event_type"` // e.g., "inventory.changed", "forecast.generated", "anomaly.detected"
    WebhookURL     string    `db:"webhook_url"`
    SecretKey      string    `db:"secret_key"` // For HMAC-SHA256 signing
    IsActive       bool      `db:"is_active"`
    MaxRetries     int       `db:"max_retries"`
    RetryInterval  int       `db:"retry_interval"` // Seconds between retries
    CreatedAt      time.Time `db:"created_at"`
    UpdatedAt      time.Time `db:"updated_at"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
    ID               uuid.UUID      `db:"id"`
    SubscriptionID   uuid.UUID      `db:"subscription_id"`
    EventID          string         `db:"event_id"` // Unique event identifier (for idempotency)
    Payload          []byte         `db:"payload"`  // JSONB: Event data
    Status           DeliveryStatus `db:"status"`
    HTTPStatusCode   *int           `db:"http_status_code"`
    ErrorMessage     *string        `db:"error_message"`
    AttemptCount     int            `db:"attempt_count"`
    LastAttemptAt    *time.Time     `db:"last_attempt_at"`
    ScheduledRetryAt *time.Time     `db:"scheduled_retry_at"`
    CreatedAt        time.Time      `db:"created_at"`
    UpdatedAt        time.Time      `db:"updated_at"`
}

// DeliveryStatus represents webhook delivery status
type DeliveryStatus string

const (
    DeliveryPending    DeliveryStatus = "pending"     // Waiting to be sent
    DeliverySuccess    DeliveryStatus = "success"     // Successfully delivered
    DeliveryFailed     DeliveryStatus = "failed"      // Failed after max retries
    DeliveryDeadLetter DeliveryStatus = "dead_letter" // Moved to dead-letter queue
)

// WebhookEvent represents an event payload
type WebhookEvent struct {
    EventID   string      `json:"event_id"`   // Unique identifier
    EventType string      `json:"event_type"` // Type of event
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data"`
}