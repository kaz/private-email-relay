package assign

import (
	"context"
)

type (
	Strategy interface {
		Assign(ctx context.Context, url string) (assignedAddr string, err error)
		Unassign(ctx context.Context, url string) error
		UnassignByAddr(ctx context.Context, addr string) error
	}
)
