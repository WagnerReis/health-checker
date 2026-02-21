package entities

import (
	domainerrors "health-checker/internal/domain/errors"
	valueobjects "health-checker/internal/shared/value-object"
	"testing"

	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_NewUser(t *testing.T) {
	email, err := valueobjects.NewEmail("john.doe@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}
	user, err := NewUser(uuid.Nil, "John Doe", email.String(), "password", nil)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john.doe@example.com", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Nil(t, user.RefreshToken)
}

func TestUser_NewUser_ErrorWhenEmailIsInvalid(t *testing.T) {
	user, err := NewUser(uuid.Nil, "John Doe", "invalid-email", "password", nil)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "invalid email")
}

func TestUser_NewUser_ErrorWhenNameIsEmpty(t *testing.T) {
	user, err := NewUser(uuid.Nil, "", "john.doe@example.com", "password", nil)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, domainerrors.ErrUserNameRequired.Error())
}

func TestUser_NewUser_ErrorWhenPasswordIsEmpty(t *testing.T) {
	user, err := NewUser(uuid.Nil, "John Doe", "john.doe@example.com", "", nil)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, domainerrors.ErrUserPasswordRequired.Error())
}

func TestUser_NewUser_ErrorWhenPasswordIsLessThan8Characters(t *testing.T) {
	user, err := NewUser(uuid.Nil, "John Doe", "john.doe@example.com", "pass", nil)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, domainerrors.ErrUserPasswordTooShort.Error())
}
