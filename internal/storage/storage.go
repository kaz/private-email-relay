package storage

import (
	"context"
	"fmt"
	"time"
)

type (
	Storage interface {
		// returns ErrorUndefinedKey
		Get(ctx context.Context, key string) (value string, err error)
		// returns ErrorDuplicatedKey, ErrorDuplicatedValue
		Set(ctx context.Context, key, value string, expires time.Time) (err error)
		// returns ErrorUndefinedKey
		UnsetByKey(ctx context.Context, key string) (deletedValue string, err error)
		// returns ErrorUndefinedValue
		UnsetByValue(ctx context.Context, value string) (deletedValue string, err error)
		// returns [Nothing]
		UnsetExpired(ctx context.Context, until time.Time) (deletedValues []string, err error)
	}
)

var (
	ErrorUndefinedKey    = fmt.Errorf("undefined key")
	ErrorUndefinedValue  = fmt.Errorf("undefined value")
	ErrorDuplicatedKey   = fmt.Errorf("duplicated key")
	ErrorDuplicatedValue = fmt.Errorf("duplicated value")

	NeverExpire = time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)
)
