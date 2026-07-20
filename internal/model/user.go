package model

import (
"time"

"github.com/google/uuid"
)

// User represents a system user
type User struct {
ID                uuid.UUID
Email             string
PasswordHash      string // Never expose in JSON
FullName          string
Role              string    // admin, manager, analyst, viewer
Status            string    // active, suspended, pending_email_verification
LastLoginAt       *time.Time
CreatedAt         time.Time
UpdatedAt         time.Time
}

// UserResponse is the JSON representation of User (without sensitive fields)
type UserResponse struct {
ID            uuid.UUID  `json:"id"`
Email         string     `json:"email"`
FullName      string     `json:"full_name"`
Role          string     `json:"role"`
Status        string     `json:"status"`
LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
CreatedAt     time.Time  `json:"created_at"`
}

// ToResponse converts User to UserResponse (hiding password)
func (u *User) ToResponse() UserResponse {
return UserResponse{
ID:          u.ID,
Email:       u.Email,
FullName:    u.FullName,
Role:        u.Role,
Status:      u.Status,
LastLoginAt: u.LastLoginAt,
CreatedAt:   u.CreatedAt,
}
}

// ValidRole checks if role is valid
func ValidRole(role string) bool {
validRoles := map[string]bool{
"admin":   true,
"manager": true,
"analyst": true,
"viewer":  true,
}
return validRoles[role]
}

// ValidStatus checks if status is valid
func ValidStatus(status string) bool {
validStatuses := map[string]bool{
"active":                      true,
"suspended":                   true,
"pending_email_verification": true,
}
return validStatuses[status]
}

// RoleHierarchy returns true if role1 has higher priority than role2
func RoleHierarchy(role1, role2 string) bool {
hierarchy := map[string]int{
"admin":   4,
"manager": 3,
"analyst": 2,
"viewer":  1,
}
return hierarchy[role1] >= hierarchy[role2]
}
