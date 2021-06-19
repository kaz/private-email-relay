package assign

import "context"

type (
	Strategy interface {
		Assign(ctx context.Context, url string) (string, error)
		UnassignByUrl(ctx context.Context, addr string) error
		UnassignByAddr(ctx context.Context, addr string) error
	}
)
