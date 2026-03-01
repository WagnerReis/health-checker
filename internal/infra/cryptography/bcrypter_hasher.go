package cryptography

import "golang.org/x/crypto/bcrypt"

type BcrypterHasher struct{}

func NewBcrypterHasher() *BcrypterHasher {
	return &BcrypterHasher{}
}

func (h *BcrypterHasher) Hash(password string) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPasswordString := string(hashedPassword)
	return &hashedPasswordString, nil
}

func (h *BcrypterHasher) Compare(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
