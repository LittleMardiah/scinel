package repository

import (
"context"
"database/sql"
"fmt"

"github.com/google/uuid"
"github.com/LittleMardiah/scinel/internal/database"
"github.com/LittleMardiah/scinel/internal/model"
)

// UserRepository defines user database operations
type UserRepository interface {
CreateUser(ctx context.Context, email, passwordHash, fullName, role string) (*model.User, error)
GetUserByEmail(ctx context.Context, email string) (*model.User, error)
GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error
UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error
UpdateLastLogin(ctx context.Context, id uuid.UUID) error
ListUsers(ctx context.Context, limit, offset int) ([]model.User, error)
CountUsers(ctx context.Context) (int, error)
}

// userRepository implements UserRepository
type userRepository struct {
db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) UserRepository {
return &userRepository{db: db}
}

// CreateUser inserts a new user
func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash, fullName, role string) (*model.User, error) {
query := `
INSERT INTO users (email, password_hash, full_name, role, status)
VALUES ($1, $2, $3, $4, 'active')
RETURNING id, email, password_hash, full_name, role, status, last_login_at, created_at, updated_at
`

user := &model.User{}
err := r.db.QueryRowContext(ctx, query, email, passwordHash, fullName, role).Scan(
&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
&user.Role, &user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
)
if err != nil {
if err == sql.ErrNoRows {
return nil, fmt.Errorf("user not created")
}
return nil, fmt.Errorf("failed to create user: %w", err)
}

return user, nil
}

// GetUserByEmail retrieves user by email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
query := `
SELECT id, email, password_hash, full_name, role, status, last_login_at, created_at, updated_at
FROM users WHERE email = $1
`

user := &model.User{}
err := r.db.QueryRowContext(ctx, query, email).Scan(
&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
&user.Role, &user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
)
if err != nil {
if err == sql.ErrNoRows {
return nil, fmt.Errorf("user not found")
}
return nil, fmt.Errorf("failed to get user: %w", err)
}

return user, nil
}

// GetUserByID retrieves user by ID
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
query := `
SELECT id, email, password_hash, full_name, role, status, last_login_at, created_at, updated_at
FROM users WHERE id = $1
`

user := &model.User{}
err := r.db.QueryRowContext(ctx, query, id).Scan(
&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
&user.Role, &user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
)
if err != nil {
if err == sql.ErrNoRows {
return nil, fmt.Errorf("user not found")
}
return nil, fmt.Errorf("failed to get user: %w", err)
}

return user, nil
}

// UpdateUserRole updates user's role
func (r *userRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error {
query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
result, err := r.db.ExecContext(ctx, query, role, id)
if err != nil {
return fmt.Errorf("failed to update user role: %w", err)
}

rowsAffected, err := result.RowsAffected()
if err != nil {
return fmt.Errorf("failed to get rows affected: %w", err)
}

if rowsAffected == 0 {
return fmt.Errorf("user not found")
}

return nil
}

// UpdateUserStatus updates user's status
func (r *userRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error {
query := `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`
result, err := r.db.ExecContext(ctx, query, status, id)
if err != nil {
return fmt.Errorf("failed to update user status: %w", err)
}

rowsAffected, err := result.RowsAffected()
if err != nil {
return fmt.Errorf("failed to get rows affected: %w", err)
}

if rowsAffected == 0 {
return fmt.Errorf("user not found")
}

return nil
}

// UpdateLastLogin updates last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
query := `UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`
_, err := r.db.ExecContext(ctx, query, id)
return err
}

// ListUsers retrieves users with pagination
func (r *userRepository) ListUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
query := `
SELECT id, email, password_hash, full_name, role, status, last_login_at, created_at, updated_at
FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
`

rows, err := r.db.QueryContext(ctx, query, limit, offset)
if err != nil {
return nil, fmt.Errorf("failed to list users: %w", err)
}
defer rows.Close()

var users []model.User
for rows.Next() {
var u model.User
if err = rows.Scan(
&u.ID, &u.Email, &u.PasswordHash, &u.FullName,
&u.Role, &u.Status, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
); err != nil {
return nil, fmt.Errorf("failed to scan user: %w", err)
}
users = append(users, u)
}

return users, rows.Err()
}

// CountUsers returns total user count
func (r *userRepository) CountUsers(ctx context.Context) (int, error) {
var count int
query := `SELECT COUNT(*) FROM users`
err := r.db.QueryRowContext(ctx, query).Scan(&count)
return count, err
}