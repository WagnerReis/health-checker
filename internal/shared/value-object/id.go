package valueobjects

import "github.com/google/uuid"

type ID struct {
	value uuid.UUID
}

func NewID(value uuid.UUID) *ID {
	if value == uuid.Nil {
		value = uuid.New()
	}
	return &ID{value: value}
}

func (id *ID) Value() uuid.UUID {
	return id.value
}

func (id *ID) String() uuid.UUID {
	return id.value
}
