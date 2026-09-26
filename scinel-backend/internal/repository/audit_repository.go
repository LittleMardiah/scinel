package repository

import (
	"context"

	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/google/uuid"
)

// AuditRepository defines audit trail data access operations.
//
// The writer is deliberately append-only and non-blocking for the caller: a
// failed audit insert must never roll back the business transaction it
// describes (see service.AuditService.Log).
type AuditRepository interface {
	// Create appends a single audit entry.
	Create(ctx context.Context, entry *model.AuditLog) error

	// GetByUser retrieves audit entries for one user, newest first.
	// user_id is joined with a LEFT JOIN so entries survive user deletion
	// (ON DELETE SET NULL) and still return an empty user_email in that case.
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.AuditLog, error)
}
