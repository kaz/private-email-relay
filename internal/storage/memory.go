package storage

import (
	"context"
	"fmt"
	"sync"
)

type (
	MemoryStorage struct {
		data sync.Map
	}
)

func NewMemoryStorage() Storage {
	return &MemoryStorage{sync.Map{}}
}

func (s *MemoryStorage) Get(ctx context.Context, key string) (string, error) {
	val, ok := s.data.Load(key)
	if !ok {
		return "", fmt.Errorf("%w: %v", ErrorNotFound, key)
	}
	return val.(string), nil
}

func (s *MemoryStorage) Set(ctx context.Context, key, value string) error {
	s.data.Store(key, value)
	return nil
}

func (s *MemoryStorage) Unset(ctx context.Context, key string) error {
	s.data.Delete(key)
	return nil
}
