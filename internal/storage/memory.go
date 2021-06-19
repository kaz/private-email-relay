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
	return &MemoryStorage{data: sync.Map{}}
}

func (s *MemoryStorage) find(needle string) *string {
	var ret *string
	s.data.Range(func(key, value interface{}) bool {
		if value.(string) == needle {
			keyStr := key.(string)
			ret = &keyStr
		}
		return true
	})
	return ret
}

func (s *MemoryStorage) Get(ctx context.Context, key string) (string, error) {
	val, ok := s.data.Load(key)
	if !ok {
		return "", fmt.Errorf("%w: key=%v", ErrorUndefinedKey, key)
	}
	return val.(string), nil
}

func (s *MemoryStorage) Set(ctx context.Context, key, value string) error {
	if s.find(value) != nil {
		return fmt.Errorf("%w: value=%v", ErrorDuplicatedValue, value)
	}

	if _, loaded := s.data.LoadOrStore(key, value); loaded {
		return fmt.Errorf("%w: key=%v", ErrorDuplicatedKey, key)
	}
	return nil
}

func (s *MemoryStorage) UnsetByKey(ctx context.Context, key string) error {
	if _, loaded := s.data.LoadAndDelete(key); !loaded {
		return fmt.Errorf("%w: key=%v", ErrorUndefinedKey, key)
	}
	return nil
}

func (s *MemoryStorage) UnsetByValue(ctx context.Context, value string) error {
	keyRef := s.find(value)
	if keyRef == nil {
		return fmt.Errorf("%w: value=%v", ErrorUndefinedValue, value)
	}

	key := *keyRef
	if _, loaded := s.data.LoadAndDelete(key); !loaded {
		return fmt.Errorf("%w: key=%v", ErrorUndefinedKey, key)
	}
	return nil
}
