package cryptography

type SHA256Hash interface {
	Hash([]byte) string
}
