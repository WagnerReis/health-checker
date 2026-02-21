package cryptography

import (
	"crypto/sha256"
	"encoding/hex"
)

type SHA256Hash struct{}

func NewSHA256Hash() *SHA256Hash {
	return &SHA256Hash{}
}

func (h *SHA256Hash) Hash(input []byte) string {
	return hex.EncodeToString(sha256.New().Sum(input))
}
