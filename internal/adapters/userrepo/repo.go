package userrepo

import (
	"context"
	"errors"
	"homework9/internal/users"
	"sync"
)

type RepositoryMap struct {
	repo   map[int64]users.User
	lastId int64
	mx     *sync.RWMutex
}

func New() *RepositoryMap {
	return &RepositoryMap{repo: make(map[int64]users.User), lastId: -1, mx: &sync.RWMutex{}}
}

var ErrNotFound = errors.New("not found")

func (r *RepositoryMap) GetUserByID(ctx context.Context, id int64) (*users.User, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	user, ok := r.repo[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *RepositoryMap) AddUser(ctx context.Context, u users.User) (int64, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.lastId++
	id := r.lastId
	r.repo[id] = u
	return id, nil
}
func (r *RepositoryMap) UpdateByID(ctx context.Context, id int64, u users.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.repo[id]
	if !ok {
		return errors.New("ad not found")
	}
	r.repo[id] = u
	return nil
}

func (r *RepositoryMap) DeleteUser(ctx context.Context, id int64) {
	r.mx.Lock()
	defer r.mx.Unlock()
	delete(r.repo, id)
}
