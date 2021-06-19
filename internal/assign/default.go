package assign

import (
	"context"
	"errors"
	"fmt"
	"os"

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

func (s *DefaultStrategy) Assign(ctx context.Context, url string) (string, error) {
	edom, err := effectiveDomain(url)
	if err != nil {
		return "", fmt.Errorf("effective domain error: %w", err)
	}

	val, err := s.store.Get(ctx, edom)
	if err != nil && !errors.Is(err, storage.ErrorNotFound) {
		return "", fmt.Errorf("failed to get value from storage: %w", err)
	}
	if val != "" {
		return val, nil
	}

	assignedAddr := fmt.Sprintf("%s@%s", randomString(4), s.emailDomain)
	if err := s.store.Set(ctx, edom, assignedAddr); err != nil {
		return "", fmt.Errorf("failed to put value to storage: %w", err)
	}

	if err := s.route.Set(ctx, assignedAddr, s.recipientAddr); err != nil {
		return "", fmt.Errorf("failed to create route: %w", err)
	}

	return assignedAddr, nil
}
