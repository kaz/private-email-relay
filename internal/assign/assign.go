package assign

import (
	"context"
	"time"
)

type (
	Strategy interface {
		Assign(ctx context.Context, url string) (assignedAddr string, err error)
		AssignTemporary(ctx context.Context, url string) (assignedAddr string, err error)
		Unassign(ctx context.Context, url string) error
		UnassignTemporary(ctx context.Context, url string) error
		UnassignByAddr(ctx context.Context, addr string) error
		UnassignExpired(ctx context.Context, until time.Time) (unassignedCount int, err error)
	}
)
