package storage

import (
	"context"
	"fmt"
)

type (
	Storage interface {
		// returns ErrorUndefinedKey
		Get(ctx context.Context, key string) (string, error)
		// returns ErrorDuplicatedKey, ErrorDuplicatedValue
		Set(ctx context.Context, key, value string) error
		// returns ErrorUndefinedKey
		UnsetByKey(ctx context.Context, key string) error
		// returns ErrorUndefinedValue
		UnsetByValue(ctx context.Context, value string) error
	}
)

var (
	ErrorUndefinedKey    = fmt.Errorf("undefined key")
	ErrorUndefinedValue  = fmt.Errorf("undefined value")
	ErrorDuplicatedKey   = fmt.Errorf("duplicated key")
	ErrorDuplicatedValue = fmt.Errorf("duplicated value")
)
