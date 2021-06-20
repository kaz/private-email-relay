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
	DefaultStrategy struct {
		emailDomain   string
		recipientAddr string

		store storage.Storage
		route router.Router
	}
)

func NewDefaultStrategy(store storage.Storage, route router.Router) (Strategy, error) {
	strategy := &DefaultStrategy{
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

func (s *DefaultStrategy) tmpKey(key string) string {
	return fmt.Sprintf("#tmp#%s", key)
}

func (s *DefaultStrategy) assignByKey(ctx context.Context, key string, randLen int, expires time.Time) (string, error) {
	val, err := s.store.Get(ctx, key)
	if err != nil && !errors.Is(err, storage.ErrorUndefinedKey) {
		return "", fmt.Errorf("failed to get value from storage: %w", err)
	}
	if val != "" {
		return val, nil
	}

	assignedAddr := fmt.Sprintf("%s@%s", randomString(randLen), s.emailDomain)
	if err := s.store.Set(ctx, key, assignedAddr, expires); err != nil {
		return "", fmt.Errorf("failed to write to storage: %w", err)
	}

	if err := s.route.Set(ctx, assignedAddr, s.recipientAddr); err != nil {
		return "", fmt.Errorf("failed to create route: %w", err)
	}

	return assignedAddr, nil
}
func (s *DefaultStrategy) Assign(ctx context.Context, url string) (string, error) {
	key, err := effectiveDomain(url)
	if err != nil {
		return "", fmt.Errorf("error occurred while finding effective domain: %w", err)
	}
	return s.assignByKey(ctx, key, 4, storage.NeverExpire)
}
func (s *DefaultStrategy) AssignTemporary(ctx context.Context, url string) (string, error) {
	key, err := effectiveDomain(url)
	if err != nil {
		return "", fmt.Errorf("error occurred while finding effective domain: %w", err)
	}
	return s.assignByKey(ctx, s.tmpKey(key), 10, time.Now().Add(7*24*time.Hour))
}

func (s *DefaultStrategy) unassignByKey(ctx context.Context, key string) error {
	addr, err := s.store.UnsetByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to determine address: %w", err)
	}

	if err := s.route.Unset(ctx, addr); err != nil {
		return fmt.Errorf("failed to remove route: %w", err)
	}
	return nil
}
func (s *DefaultStrategy) Unassign(ctx context.Context, url string) error {
	key, err := effectiveDomain(url)
	if err != nil {
		return fmt.Errorf("error occurred while finding effective domain: %w", err)
	}
	return s.unassignByKey(ctx, key)
}
func (s *DefaultStrategy) UnassignTemporary(ctx context.Context, url string) error {
	key, err := effectiveDomain(url)
	if err != nil {
		return fmt.Errorf("error occurred while finding effective domain: %w", err)
	}
	return s.unassignByKey(ctx, s.tmpKey(key))
}

func (s *DefaultStrategy) UnassignByAddr(ctx context.Context, addr string) error {
	if _, err := s.store.UnsetByValue(ctx, addr); err != nil {
		return fmt.Errorf("failed to delete from storage: %w", err)
	}

	if err := s.route.Unset(ctx, addr); err != nil {
		return fmt.Errorf("failed to remove route: %w", err)
	}
	return nil
}

func (s *DefaultStrategy) UnassignExpired(ctx context.Context, until time.Time) (int, error) {
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
