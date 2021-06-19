package router

import (
	"context"
	"fmt"
)

type (
	Router interface {
		// retuns ErrorDuplicated
		Set(ctx context.Context, from, to string) error
		// returns ErrorUndefined
		Unset(ctx context.Context, from string) error
	}
)

var (
	ErrorDuplicated = fmt.Errorf("duplicated")
	ErrorUndefined  = fmt.Errorf("undefined")
)
