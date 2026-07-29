package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
)

// WebhookRepository defines webhook data access operations
type WebhookRepository interface {
    // Create subscription
    CreateSubscription(ctx context.Context, sub *model.WebhookSubscription) (*model.WebhookSubscription, error)

    // Get subscriptions for user
    GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.WebhookSubscription, error)

    // Get subscriptions by event type (for publishing)
    GetSubscriptionsByEventType(ctx context.Context, eventType string) ([]model.WebhookSubscription, error)

    // Get subscription by ID
    GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*model.WebhookSubscription, error)

    // Update subscription
    UpdateSubscription(ctx context.Context, sub *model.WebhookSubscription) error

    // Delete subscription
    DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error

    // Record delivery attempt
    CreateDelivery(ctx context.Context, delivery *model.WebhookDelivery) (*model.WebhookDelivery, error)

    // Get pending deliveries
    GetPendingDeliveries(ctx context.Context, limit int) ([]model.WebhookDelivery, error)

    // Update delivery status
    UpdateDeliveryStatus(ctx context.Context, deliveryID uuid.UUID, status model.DeliveryStatus, httpStatus *int, errMsg *string) error

    // Get delivery history
    GetDeliveryHistory(ctx context.Context, subscriptionID uuid.UUID, limit, offset int) ([]model.WebhookDelivery, int64, error)
}