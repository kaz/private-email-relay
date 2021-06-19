package assign

import "context"

type (
	Strategy interface {
		Assign(ctx context.Context, url string) (string, error)
		Unassign(ctx context.Context, addr string) error
	}
)
