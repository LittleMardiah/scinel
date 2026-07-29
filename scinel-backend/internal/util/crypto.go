package util

import (
    "crypto/sha256"
    "encoding/hex"
)

// HashSHA256 computes SHA256 hash of a string (for API key hashing)
func HashSHA256(input string) string {
    sum := sha256.Sum256([]byte(input))
    return hex.EncodeToString(sum[:])
}