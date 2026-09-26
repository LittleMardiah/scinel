package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Audit action constants — the closed set of events the application writes.
const (
	AuditActionLogin            = "LOGIN"
	AuditActionLogout           = "LOGOUT"
	AuditActionInventoryChanged = "INVENTORY_CHANGED"
	AuditActionAnomalyDetected  = "ANOMALY_DETECTED"
	AuditActionWebhookPublished = "WEBHOOK_PUBLISHED"
)

// Audit resource types used with AuditActionInventoryChanged / AnomalyDetected.
const (
	AuditResourceTransaction  = "transaction"
	AuditResourceInventory    = "inventory"
	AuditResourceAnomaly      = "anomaly"
	AuditResourceWebhookEvent = "webhook_event"
)

// AuditLog mirrors table audit_log created by migration
// 000009_create_audit_log_table.up.sql.
//
// Nullable columns are represented as pointers so a genuinely absent value can
// be told apart from a zero value. Details is JSONB and therefore optional.
type AuditLog struct {
	ID           uuid.UUID       `db:"id"`
	UserID       *uuid.UUID      `db:"user_id"` // NULL for system/unauthenticated events
	Action       string          `db:"action"`
	ResourceType *string         `db:"resource_type"` // NULL when the event has no resource
	ResourceID   *string         `db:"resource_id"`   // NULL when the event has no resource
	Details      json.RawMessage `db:"details"`       // NULL when there is no payload
	IPAddress    *string         `db:"ip_address"`    // NULL when the client IP is unknown
	UserAgent    *string         `db:"user_agent"`    // NULL when the header is absent
	CreatedAt    time.Time       `db:"created_at"`

	// UserEmail is not an audit_log column: it is filled from the LEFT JOIN
	// against users performed by AuditRepository.GetByUser. Empty when the
	// user row is gone (ON DELETE SET NULL) or was never set.
	UserEmail string `db:"user_email"`
}
