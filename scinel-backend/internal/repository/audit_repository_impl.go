package repository

import (
	"context"
	"fmt"

	"github.com/LittleMardiah/scinel/internal/database"
	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/google/uuid"
)

// auditRepository implements AuditRepository interface
type auditRepository struct {
	db *database.DB
}

// NewAuditRepository creates a new audit repository instance
func NewAuditRepository(db *database.DB) AuditRepository {
	return &auditRepository{db: db}
}

// --- IMPLEMENTATIONS ---

// Create appends a single audit entry.
//
// JSONB handling follows the existing repository convention
// (webhook_repository_impl.go:173): details is written as text and read back
// with an explicit ::text cast, because lib/pq has no native JSONB driver
// value. Nullable columns are passed as typed nil pointers so PostgreSQL stores
// SQL NULL rather than the string "NULL" or an empty string.
func (r *auditRepository) Create(ctx context.Context, entry *model.AuditLog) error {
	if entry == nil {
		return fmt.Errorf("audit entry must not be nil")
	}

	var details interface{}
	if len(entry.Details) > 0 {
		details = string(entry.Details)
	}

	query := `
        INSERT INTO audit_log (
            id, user_id, action, resource_type, resource_id,
            details, ip_address, user_agent, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	if _, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.UserID, entry.Action, entry.ResourceType,
		entry.ResourceID, details, entry.IPAddress, entry.UserAgent,
		entry.CreatedAt,
	); err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

// GetByUser retrieves audit entries for one user, newest first.
//
// LEFT JOIN (not INNER JOIN) is deliberate: audit_log.user_id uses
// ON DELETE SET NULL, so the row disappears from this listing when its user is
// deleted, but the join must never drop rows that still have a live user_id.
// Scanning user_email into a *string therefore yields nil, not an error, for
// the "system event" shape where user_id is NULL.
func (r *auditRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.AuditLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `
        SELECT
            a.id, a.user_id, a.action, a.resource_type, a.resource_id,
            a.details::text, a.ip_address, a.user_agent, a.created_at,
            u.email
        FROM audit_log a
        LEFT JOIN users u ON a.user_id = u.id
        WHERE a.user_id = $1
        ORDER BY a.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit log by user: %w", err)
	}
	defer rows.Close()

	var items []model.AuditLog
	for rows.Next() {
		var entry model.AuditLog
		var details []byte
		var userEmail *string

		if err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.Action,
			&entry.ResourceType, &entry.ResourceID,
			&details, &entry.IPAddress, &entry.UserAgent, &entry.CreatedAt,
			&userEmail,
		); err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if len(details) > 0 {
			entry.Details = append([]byte(nil), details...)
		}
		if userEmail != nil {
			entry.UserEmail = *userEmail
		}
		items = append(items, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}
