package storage

import (
	"context"
	"fmt"
)

type (
	Storage interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key, value string) error
		Unset(ctx context.Context, key string) error
	}
)

var (
	ErrorNotFound = fmt.Errorf("no such entry")
)
