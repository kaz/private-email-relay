package router

import (
	"context"
	"fmt"
)

type (
	Router interface {
		Set(ctx context.Context, from, to string) error
		Unset(ctx context.Context, from string) error
	}
)

var (
	ErrorDuplicated       = fmt.Errorf("duplicated")
	ErrorUnsetNonexistent = fmt.Errorf("no such entry")
)
