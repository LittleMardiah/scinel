package model

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidRole(t *testing.T) {
	assert.True(t, ValidRole("admin"))
	assert.False(t, ValidRole("superadmin"))
}

func TestValidStatus(t *testing.T) {
	assert.True(t, ValidStatus("active"))
	assert.False(t, ValidStatus("deleted"))
}

func TestRoleHierarchy(t *testing.T) {
	assert.True(t, RoleHierarchy("admin", "viewer"))
	assert.False(t, RoleHierarchy("analyst", "manager"))
}

func TestUserToResponse(t *testing.T) {
	now := time.Now().UTC()
	user := &User{
		ID: uuid.New(),
		Email: "test@example.com",
		PasswordHash: "hash",
		FullName: "Test",
		Role: "admin",
		Status: "active",
		CreatedAt: now,
	}
	resp := user.ToResponse()
	assert.Equal(t, "test@example.com", resp.Email)

}
