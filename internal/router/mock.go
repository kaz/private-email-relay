package router

import (
	"context"
	"fmt"
	"sync"
)

type (
	MockRouter struct {
		data sync.Map
	}
)

func NewMockRouter() Router {
	return &MockRouter{}
}

func (r *MockRouter) Set(ctx context.Context, from, to string) error {
	if _, loaded := r.data.LoadOrStore(from, to); loaded {
		return fmt.Errorf("%w: %v", ErrorDuplicated, from)
	}
	return nil
}
func (r *MockRouter) Unset(ctx context.Context, from string) error {
	if _, loaded := r.data.LoadAndDelete(from); !loaded {
		return fmt.Errorf("%w: %v", ErrorUndefined, from)
	}
	return nil
}
