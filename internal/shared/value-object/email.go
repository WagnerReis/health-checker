package valueobjects

import (
	"fmt"
	"net/mail"
)

type Email string

func NewEmail(value string) (*Email, error) {
	_, err := mail.ParseAddress(value)
	if err != nil {
		return nil, fmt.Errorf("invalid email")
	}
	email := Email(value)
	return &email, nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Value() string {
	return string(e)
}
