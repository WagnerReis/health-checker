package repository

import (
	"context"
	entities "health-checker/internal/domain/entity"
	"health-checker/internal/domain/errors"
	domainerrors "health-checker/internal/domain/errors"
	"sync"

	"github.com/google/uuid"
)

type UserRepositoryInMemory struct {
	users       map[uuid.UUID]*entities.User
	mu          sync.Mutex
	ErrOnCreate error
	ErrOnFind   error
	ErrOnUpdate error
}

func NewUserRepositoryInMemory() *UserRepositoryInMemory {
	return &UserRepositoryInMemory{
		users: make(map[uuid.UUID]*entities.User),
		mu:    sync.Mutex{},
	}
}

func (r *UserRepositoryInMemory) Create(ctx context.Context, user *entities.User) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *UserRepositoryInMemory) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	if !ok {
		return nil, domainerrors.ErrUserNotFound
	}
	return user, nil
}

func (r *UserRepositoryInMemory) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.ErrUserNotFound
}

func (r *UserRepositoryInMemory) Update(ctx context.Context, user *entities.User) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return domainerrors.ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}
