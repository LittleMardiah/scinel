package repository

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/model"
)

type webhookRepository struct {
    db *database.DB
}

func NewWebhookRepository(db *database.DB) WebhookRepository {
    return &webhookRepository{db: db}
}

// GetSubscriptions retrieves all webhook subscriptions for a user
func (r *webhookRepository) GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.WebhookSubscription, error) {
    query := `
        SELECT id, user_id, event_type, webhook_url, secret_key,
               is_active, max_retries, retry_interval,
               created_at, updated_at
        FROM webhook_subscriptions
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query webhook subscriptions: %w", err)
    }
    defer rows.Close()

    var subs []model.WebhookSubscription
    for rows.Next() {
        var s model.WebhookSubscription
        err := rows.Scan(
            &s.ID, &s.UserID, &s.EventType, &s.WebhookURL, &s.SecretKey,
            &s.IsActive, &s.MaxRetries, &s.RetryInterval,
            &s.CreatedAt, &s.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan webhook subscription: %w", err)
        }
        subs = append(subs, s)
    }
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    return subs, nil
}

// GetSubscriptionsByEventType retrieves active subscriptions for a specific event type
func (r *webhookRepository) GetSubscriptionsByEventType(ctx context.Context, eventType string) ([]model.WebhookSubscription, error) {
    query := `
        SELECT id, user_id, event_type, webhook_url, secret_key,
               is_active, max_retries, retry_interval,
               created_at, updated_at
        FROM webhook_subscriptions
        WHERE event_type = $1 AND is_active = true
    `
    rows, err := r.db.QueryContext(ctx, query, eventType)
    if err != nil {
        return nil, fmt.Errorf("failed to query webhook subscriptions by event: %w", err)
    }
    defer rows.Close()

    var subs []model.WebhookSubscription
    for rows.Next() {
        var s model.WebhookSubscription
        err := rows.Scan(
            &s.ID, &s.UserID, &s.EventType, &s.WebhookURL, &s.SecretKey,
            &s.IsActive, &s.MaxRetries, &s.RetryInterval,
            &s.CreatedAt, &s.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan webhook subscription: %w", err)
        }
        subs = append(subs, s)
    }
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    return subs, nil
}

// CreateSubscription inserts a new webhook subscription
func (r *webhookRepository) CreateSubscription(ctx context.Context, sub *model.WebhookSubscription) (*model.WebhookSubscription, error) {
    query := `
        INSERT INTO webhook_subscriptions (id, user_id, event_type, webhook_url, secret_key,
                                           is_active, max_retries, retry_interval, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, user_id, event_type, webhook_url, secret_key,
                  is_active, max_retries, retry_interval,
                  created_at, updated_at
    `
    row := r.db.QueryRowContext(ctx, query,
        sub.ID, sub.UserID, sub.EventType, sub.WebhookURL, sub.SecretKey,
        sub.IsActive, sub.MaxRetries, sub.RetryInterval, sub.CreatedAt, sub.UpdatedAt,
    )
    err := row.Scan(
        &sub.ID, &sub.UserID, &sub.EventType, &sub.WebhookURL, &sub.SecretKey,
        &sub.IsActive, &sub.MaxRetries, &sub.RetryInterval,
        &sub.CreatedAt, &sub.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create webhook subscription: %w", err)
    }
    return sub, nil
}

// UpdateSubscription updates a webhook subscription
func (r *webhookRepository) UpdateSubscription(ctx context.Context, sub *model.WebhookSubscription) error {
    query := `
        UPDATE webhook_subscriptions
        SET event_type = $1, webhook_url = $2, is_active = $3,
            max_retries = $4, retry_interval = $5,
            updated_at = NOW()
        WHERE id = $6 AND user_id = $7
    `
    result, err := r.db.ExecContext(ctx, query,
        sub.EventType, sub.WebhookURL, sub.IsActive,
        sub.MaxRetries, sub.RetryInterval,
        sub.ID, sub.UserID,
    )
    if err != nil {
        return fmt.Errorf("failed to update webhook subscription: %w", err)
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("webhook subscription not found or not owned by user")
    }
    return nil
}

// DeleteSubscription deletes a webhook subscription
func (r *webhookRepository) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
    query := `DELETE FROM webhook_subscriptions WHERE id = $1`
    result, err := r.db.ExecContext(ctx, query, subscriptionID)
    if err != nil {
        return fmt.Errorf("failed to delete webhook subscription: %w", err)
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("webhook subscription not found")
    }
    return nil
}

// CreateDelivery records a webhook delivery attempt
func (r *webhookRepository) CreateDelivery(ctx context.Context, delivery *model.WebhookDelivery) (*model.WebhookDelivery, error) {
    // TODO: implement when needed
    return nil, fmt.Errorf("CreateDelivery not implemented yet")
}

// GetPendingDeliveries retrieves pending deliveries for retry
func (r *webhookRepository) GetPendingDeliveries(ctx context.Context, limit int) ([]model.WebhookDelivery, error) {
    // TODO: implement when needed
    return nil, fmt.Errorf("GetPendingDeliveries not implemented yet")
}

// UpdateDeliveryStatus updates delivery status
func (r *webhookRepository) UpdateDeliveryStatus(ctx context.Context, deliveryID uuid.UUID, status model.DeliveryStatus, httpStatus *int, errMsg *string) error {
    // TODO: implement when needed
    return fmt.Errorf("UpdateDeliveryStatus not implemented yet")
}

// GetDeliveryHistory retrieves delivery history for a subscription
func (r *webhookRepository) GetDeliveryHistory(ctx context.Context, subscriptionID uuid.UUID, limit, offset int) ([]model.WebhookDelivery, int64, error) {
    // TODO: implement when needed
    return nil, 0, fmt.Errorf("GetDeliveryHistory not implemented yet")
}

// GetSubscriptionByID retrieves a webhook subscription by ID
func (r *webhookRepository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*model.WebhookSubscription, error) {
    query := `
        SELECT id, user_id, event_type, webhook_url, secret_key,
               is_active, max_retries, retry_interval,
               created_at, updated_at
        FROM webhook_subscriptions
        WHERE id = $1
    `
    var s model.WebhookSubscription
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &s.ID, &s.UserID, &s.EventType, &s.WebhookURL, &s.SecretKey,
        &s.IsActive, &s.MaxRetries, &s.RetryInterval,
        &s.CreatedAt, &s.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get webhook subscription: %w", err)
    }
    return &s, nil
}
