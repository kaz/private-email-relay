package assign

import "context"

type (
	Strategy interface {
		Assign(ctx context.Context, url string) (string, error)
	}
)
