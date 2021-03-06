package storage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type (
	MemoryStorage struct {
		data map[string]memoryStorageEntry
		mu   sync.RWMutex
	}
	memoryStorageEntry struct {
		value   string
		expires time.Time
	}
)

func NewMemoryStorage() Storage {
	return &MemoryStorage{
		data: map[string]memoryStorageEntry{},
		mu:   sync.RWMutex{},
	}
}

func (s *MemoryStorage) find(needle string) (string, bool) {
	for key, entry := range s.data {
		if entry.value == needle {
			return key, true
		}
	}
	return "", false
}

func (s *MemoryStorage) Get(ctx context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("%w: key=%v", ErrorUndefinedKey, key)
	}
	return entry.value, nil
}

func (s *MemoryStorage) Set(ctx context.Context, key, value string, expires time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; ok {
		return fmt.Errorf("%w: key=%v", ErrorDuplicatedKey, key)
	}
	if _, ok := s.find(value); ok {
		return fmt.Errorf("%w: value=%v", ErrorDuplicatedValue, value)
	}

	s.data[key] = memoryStorageEntry{value, expires}
	return nil
}

func (s *MemoryStorage) UnsetByKey(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("%w: key=%v", ErrorUndefinedKey, key)
	}

	delete(s.data, key)
	return entry.value, nil
}

func (s *MemoryStorage) UnsetByValue(ctx context.Context, value string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key, ok := s.find(value)
	if !ok {
		return "", fmt.Errorf("%w: value=%v", ErrorUndefinedValue, value)
	}

	delete(s.data, key)
	return value, nil
}

func (s *MemoryStorage) UnsetExpired(ctx context.Context, until time.Time) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valuesExpired := []string{}

	for key, entry := range s.data {
		if until.After(entry.expires) {
			valuesExpired = append(valuesExpired, entry.value)
			delete(s.data, key)
		}
	}
	return valuesExpired, nil
}
