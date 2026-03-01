package criptography

type FakeHasher struct {
	ErrOnHash error
}

func NewFakeHasher() *FakeHasher {
	return &FakeHasher{}
}

func (h *FakeHasher) Hash(password string) (*string, error) {
	if h.ErrOnHash != nil {
		return nil, h.ErrOnHash
	}
	hashedPassword := password + "-hashed"
	return &hashedPassword, nil
}

func (h *FakeHasher) Compare(password, hash string) bool {
	return password+"-hashed" == hash
}
