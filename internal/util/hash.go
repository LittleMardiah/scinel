package util

import (
"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string, cost int) (string, error) {
if cost == 0 {
cost = 12 // Default cost from .env
}
hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
if err != nil {
return "", err
}
return string(hash), nil
}

// VerifyPassword compares password with hash
func VerifyPassword(password, hash string) error {
return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}