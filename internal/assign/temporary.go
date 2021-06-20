package assign

import (
	"context"
	"fmt"
	"time"

	"github.com/kaz/private-email-relay/internal/router"
	"github.com/kaz/private-email-relay/internal/storage"
)

type (
	TemporaryStrategy struct {
		*baseStrategy

		deadline deadline
	}

	deadline func() time.Time
)

func NewTemporaryStrategy(store storage.Storage, route router.Router, deadline deadline) (Strategy, error) {
	base, err := newBaseStrategy(store, route)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize base strategy: %w", err)
	}
	return &TemporaryStrategy{base, deadline}, nil
}

func (s *TemporaryStrategy) keyProducerFactory(url string) producer {
	return func() (string, error) {
		key, err := effectiveDomain(url)
		if err != nil {
			return "", fmt.Errorf("error occurred while finding effective domain: %w", err)
		}
		return fmt.Sprintf("temp#%s", key), nil
	}
}

func (s *TemporaryStrategy) Assign(ctx context.Context, url string) (string, error) {
	return s.assignByKey(ctx, s.keyProducerFactory(url), s.addressProducerFactory("t-", 6), s.deadline())
}

func (s *TemporaryStrategy) Unassign(ctx context.Context, url string) error {
	return s.unassignByKey(ctx, s.keyProducerFactory(url))
}

func (s *TemporaryStrategy) UnassignByAddr(ctx context.Context, addr string) error {
	return s.unassignByAddr(ctx, addr)
}

func (s *TemporaryStrategy) UnassignExpired(ctx context.Context, until time.Time) (int, error) {
	deletedAddrs, err := s.store.UnsetExpired(ctx, until)
	if err != nil {
		return 0, fmt.Errorf("failed to delete from storage: %w", err)
	}

	for _, addr := range deletedAddrs {
		if err := s.route.Unset(ctx, addr); err != nil {
			return 0, fmt.Errorf("failed to remove route: %w", err)
		}
	}

	return len(deletedAddrs), nil
}
