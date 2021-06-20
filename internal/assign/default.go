package assign

import (
	"context"
	"fmt"

	"github.com/kaz/private-email-relay/internal/router"
	"github.com/kaz/private-email-relay/internal/storage"
)

type (
	DefaultStrategy struct {
		*baseStrategy
	}
)

func NewDefaultStrategy(store storage.Storage, route router.Router) (Strategy, error) {
	base, err := newBaseStrategy(store, route)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base strategy: %w", err)
	}
	return &DefaultStrategy{base}, nil
}

func (s *DefaultStrategy) keyProducerFactory(url string) producer {
	return func() (string, error) {
		key, err := effectiveDomain(url)
		if err != nil {
			return "", fmt.Errorf("error occurred while finding effective domain: %w", err)
		}
		return key, nil
	}
}

func (s *DefaultStrategy) Assign(ctx context.Context, url string) (string, error) {
	return s.assignByKey(ctx, s.keyProducerFactory(url), s.addressProducerFactory("", 4), storage.NeverExpire)
}

func (s *DefaultStrategy) Unassign(ctx context.Context, url string) error {
	return s.unassignByKey(ctx, s.keyProducerFactory(url))
}

func (s *DefaultStrategy) UnassignByAddr(ctx context.Context, addr string) error {
	return s.unassignByAddr(ctx, addr)
}
