package valueobjects

import "github.com/gofrs/uuid"

type ID struct {
	value uuid.UUID
}

func NewID(value uuid.UUID) *ID {
	if value == uuid.Nil {
		id, _ := uuid.NewV4()
		value = id
	}
	return &ID{value: value}
}

func (id *ID) Value() uuid.UUID {
	return id.value
}

func (id *ID) String() uuid.UUID {
	return id.value
}
