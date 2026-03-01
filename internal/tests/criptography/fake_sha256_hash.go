package criptography

type FakeSHA256Hash struct {
	ErrOnHash error
}

func NewFakeSHA256Hash() *FakeSHA256Hash {
	return &FakeSHA256Hash{}
}

func (h *FakeSHA256Hash) Hash(input []byte) string {
	return "hash"
}
