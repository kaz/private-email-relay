package assign

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kaz/private-email-relay/internal/storage"
)

type (
	DefaultStrategy struct {
		domain string
		store  storage.Storage
	}
)

var (
	domain = os.Getenv("MG_DOMAIN")
)

func IsDefaultStrategyAvailable() error {
	if domain == "" {
		return fmt.Errorf("MG_DOMAIN is missing")
	}
	return nil
}

func NewDefaultStrategy(store storage.Storage) Strategy {
	return &DefaultStrategy{domain, store}
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

	assigned := fmt.Sprintf("%s@%s", randomString(4), s.domain)
	if err := s.store.Set(ctx, edom, assigned); err != nil {
		return "", fmt.Errorf("failed to put value to storage: %w", err)
	}

	return assigned, nil
}
