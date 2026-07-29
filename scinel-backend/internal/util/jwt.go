package util

import (
"fmt"
"time"

"github.com/golang-jwt/jwt/v5"
"github.com/google/uuid"
)

// AccessTokenClaims represents access token claims
type AccessTokenClaims struct {
UserID uuid.UUID `json:"user_id"`
Email  string    `json:"email"`
Role   string    `json:"role"`
jwt.RegisteredClaims
}

// RefreshTokenClaims represents refresh token claims
type RefreshTokenClaims struct {
UserID uuid.UUID `json:"user_id"`
jwt.RegisteredClaims
}

// JWTManager manages JWT token generation and verification
type JWTManager struct {
secretKey     string
accessExpiry  time.Duration
refreshExpiry time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
return &JWTManager{
secretKey:     secretKey,
accessExpiry:  accessExpiry,
refreshExpiry: refreshExpiry,
}
}

// GenerateAccessToken generates a new access token
func (jm *JWTManager) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
now := time.Now().UTC()
claims := AccessTokenClaims{
UserID: userID,
Email:  email,
Role:   role,
RegisteredClaims: jwt.RegisteredClaims{
IssuedAt:  jwt.NewNumericDate(now),
ExpiresAt: jwt.NewNumericDate(now.Add(jm.accessExpiry)),
},
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(jm.secretKey))
if err != nil {
return "", fmt.Errorf("failed to sign token: %w", err)
}

return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token
func (jm *JWTManager) GenerateRefreshToken(userID uuid.UUID) (string, error) {
now := time.Now().UTC()
claims := RefreshTokenClaims{
UserID: userID,
RegisteredClaims: jwt.RegisteredClaims{
IssuedAt:  jwt.NewNumericDate(now),
ExpiresAt: jwt.NewNumericDate(now.Add(jm.refreshExpiry)),
},
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(jm.secretKey))
if err != nil {
return "", fmt.Errorf("failed to sign token: %w", err)
}

return tokenString, nil
}

// VerifyAccessToken verifies and parses access token
func (jm *JWTManager) VerifyAccessToken(tokenString string) (*AccessTokenClaims, error) {
claims := &AccessTokenClaims{}
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
}
return []byte(jm.secretKey), nil
})

if err != nil {
return nil, fmt.Errorf("failed to parse token: %w", err)
}

if !token.Valid {
return nil, fmt.Errorf("invalid token")
}

return claims, nil
}

// VerifyRefreshToken verifies and parses refresh token
func (jm *JWTManager) VerifyRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
claims := &RefreshTokenClaims{}
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
}
return []byte(jm.secretKey), nil
})

if err != nil {
return nil, fmt.Errorf("failed to parse token: %w", err)
}

if !token.Valid {
return nil, fmt.Errorf("invalid token")
}

return claims, nil
}