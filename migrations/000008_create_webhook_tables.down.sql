DROP INDEX IF EXISTS idx_webhook_deliveries_event;
DROP INDEX IF EXISTS idx_webhook_deliveries_status;
DROP INDEX IF EXISTS idx_webhook_deliveries_subscription;
DROP TABLE IF EXISTS webhook_deliveries;
DROP INDEX IF EXISTS idx_webhook_subscriptions_active;
DROP INDEX IF EXISTS idx_webhook_subscriptions_event;
DROP INDEX IF EXISTS idx_webhook_subscriptions_user;
DROP TABLE IF EXISTS webhook_subscriptions;