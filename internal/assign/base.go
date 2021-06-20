package assign

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kaz/private-email-relay/internal/router"
	"github.com/kaz/private-email-relay/internal/storage"
)

type (
	baseStrategy struct {
		emailDomain   string
		recipientAddr string

		store storage.Storage
		route router.Router
	}

	producer func() (string, error)
)

func newBaseStrategy(store storage.Storage, route router.Router) (*baseStrategy, error) {
	strategy := &baseStrategy{
		store: store,
		route: route,
	}

	strategy.emailDomain = os.Getenv("MG_DOMAIN")
	if strategy.emailDomain == "" {
		return nil, fmt.Errorf("MG_DOMAIN is missing")
	}

	strategy.recipientAddr = os.Getenv("RECIPIENT")
	if strategy.recipientAddr == "" {
		return nil, fmt.Errorf("RECIPIENT is missing")
	}

	return strategy, nil
}

func (s *baseStrategy) addressProducerFactory(prefix string, randLen int) producer {
	return func() (string, error) {
		return fmt.Sprintf("%s%s@%s", prefix, randomString(randLen), s.emailDomain), nil
	}
}

func (s *baseStrategy) assignByKey(ctx context.Context, keyProd producer, addrProd producer, expires time.Time) (string, error) {
	key, err := keyProd()
	if err != nil {
		return "", fmt.Errorf("failed to produce key: %w", err)
	}

	val, err := s.store.Get(ctx, key)
	if err != nil && !errors.Is(err, storage.ErrorUndefinedKey) {
		return "", fmt.Errorf("failed to get value from storage: %w", err)
	}
	if val != "" {
		return val, nil
	}

	addr, err := addrProd()
	if err != nil {
		return "", fmt.Errorf("failed to produce address: %w", err)
	}

	if err := s.store.Set(ctx, key, addr, expires); err != nil {
		return "", fmt.Errorf("failed to write to storage: %w", err)
	}
	if err := s.route.Set(ctx, addr, s.recipientAddr); err != nil {
		return "", fmt.Errorf("failed to create route: %w", err)
	}

	return addr, nil
}
func (s *baseStrategy) unassignByKey(ctx context.Context, keyProd producer) error {
	key, err := keyProd()
	if err != nil {
		return fmt.Errorf("failed to produce key: %w", err)
	}

	addr, err := s.store.UnsetByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to determine address: %w", err)
	}

	if err := s.route.Unset(ctx, addr); err != nil {
		return fmt.Errorf("failed to remove route: %w", err)
	}
	return nil
}
func (s *baseStrategy) unassignByAddr(ctx context.Context, addr string) error {
	if _, err := s.store.UnsetByValue(ctx, addr); err != nil {
		return fmt.Errorf("failed to delete from storage: %w", err)
	}

	if err := s.route.Unset(ctx, addr); err != nil {
		return fmt.Errorf("failed to remove route: %w", err)
	}
	return nil
}
