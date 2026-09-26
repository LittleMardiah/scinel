package service

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/LittleMardiah/scinel/internal/model"
	"github.com/LittleMardiah/scinel/internal/repository"
	"github.com/google/uuid"
)

// auditUserAgentMaxLength matches VARCHAR(500) on audit_log.user_agent.
const auditUserAgentMaxLength = 500

// AuditService writes append-only audit trail entries.
//
// Every write is best-effort by design: an audit failure must never turn a
// successful business operation into an error for the caller, and must never
// roll back the transaction the entry describes. Failures are surfaced as
// Warn-level logs so they stay visible in monitoring.
type AuditService struct {
	repo   repository.AuditRepository
	logger ServiceLogger
}

// NewAuditService creates a new audit service.
func NewAuditService(repo repository.AuditRepository, logger ServiceLogger) *AuditService {
	return &AuditService{
		repo:   repo,
		logger: logger,
	}
}

// Log persists one audit entry. It is intentionally non-fatal: a nil service
// or a repository error is logged and swallowed.
func (s *AuditService) Log(ctx context.Context, entry *model.AuditLog) {
	if s == nil || s.repo == nil {
		return
	}
	if entry == nil {
		s.warn("Audit log skipped: nil entry", nil)
		return
	}

	s.applyDefaults(entry)

	if err := s.repo.Create(ctx, entry); err != nil {
		s.warn("Failed to persist audit log", map[string]interface{}{
			"action": entry.Action,
			"error":  err.Error(),
		})
	}
}

// GetByUser exposes the read side of the audit trail.
func (s *AuditService) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.AuditLog, error) {
	if s == nil || s.repo == nil {
		return nil, nil
	}
	return s.repo.GetByUser(ctx, userID, limit, offset)
}

// applyDefaults fills the fields the caller is allowed to leave blank and
// normalises the ones the database constrains.
func (s *AuditService) applyDefaults(entry *model.AuditLog) {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if entry.Action == "" {
		entry.Action = "UNKNOWN"
	}
	entry.ResourceType = trimToNil(entry.ResourceType)
	entry.ResourceID = trimToNil(entry.ResourceID)
	entry.IPAddress = trimToNil(entry.IPAddress)
	entry.UserAgent = truncateToNil(entry.UserAgent, auditUserAgentMaxLength)
}

// warn routes a message to the service logger, tolerating a nil logger so the
// service stays usable in tests that pass no logger.
func (s *AuditService) warn(msg string, fields map[string]interface{}) {
	if s.logger == nil {
		return
	}
	if len(fields) == 0 {
		s.logger.Warn(msg)
		return
	}
	flat := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		flat = append(flat, k, v)
	}
	s.logger.Warn(msg, flat...)
}

func trimToNil(v *string) *string {
    if v == nil {
        return nil
    }
    trimmed := strings.TrimSpace(*v)
    if trimmed == "" {
        return nil
    }
    return &trimmed
}

// nullableString returns nil for an empty value so nullable audit_log columns
// are stored as SQL NULL rather than an empty string.
func nullableString(v string) *string {
    if strings.TrimSpace(v) == "" {
        return nil
    }
    return &v
}

func truncateToNil(v *string, max int) *string {
	trimmed := trimToNil(v)
	if trimmed == nil {
		return nil
	}
	if len(*trimmed) > max {
		shortened := (*trimmed)[:max]
		return &shortened
	}
	return trimmed
}

// ==================== REQUEST HELPERS ====================

// ExtractIP resolves the caller IP address, preferring proxy headers.
// Returns "" when nothing usable is present so the column is stored as NULL.
func ExtractIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	candidates := []string{
		r.Header.Get("X-Forwarded-For"),
		r.Header.Get("X-Real-IP"),
	}
	for _, candidate := range candidates {
		if ip := firstForwardedIP(candidate); ip != "" {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return host
}

// firstForwardedIP extracts the left-most entry of an X-Forwarded-For chain.
func firstForwardedIP(header string) string {
	if strings.TrimSpace(header) == "" {
		return ""
	}
	first := strings.TrimSpace(strings.Split(header, ",")[0])
	if first == "" {
		return ""
	}
	if ip := net.ParseIP(first); ip != nil {
		return ip.String()
	}
	return ""
}

// ExtractUserAgent returns the raw User-Agent header, or "" when absent.
func ExtractUserAgent(r *http.Request) string {
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.Header.Get("User-Agent"))
}
