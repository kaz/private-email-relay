package router

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"
)

type (
	MailgunRouter struct {
		client *mailgun.MailgunImpl
	}
)

func NewMailgunRouter() (Router, error) {
	client, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to create Mailgun client: %w", err)
	}
	return &MailgunRouter{client}, nil
}

func (r *MailgunRouter) createExpression(from string) string {
	return fmt.Sprintf("match_recipient(\"%s\")", from)
}
func (r *MailgunRouter) createRoute(from, to string) mailgun.Route {
	return mailgun.Route{
		Expression: r.createExpression(from),
		Actions: []string{
			fmt.Sprintf("forward(\"%s\")", to),
			"stop()",
		},
		Priority: 8000,
	}
}

func (r *MailgunRouter) findRoute(ctx context.Context, from string) (*mailgun.Route, error) {
	expression := r.createExpression(from)

	iter := r.client.ListRoutes(nil)
	results := []mailgun.Route{}

	for iter.Next(ctx, &results) {
		for _, route := range results {
			if route.Expression == expression {
				return &route, nil
			}
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred while finding route: %w", err)
	}
	return nil, nil
}

func (r *MailgunRouter) Set(ctx context.Context, from, to string) error {
	route, err := r.findRoute(ctx, from)
	if err != nil {
		return fmt.Errorf("failed to find route: %w", err)
	}
	if route != nil {
		return fmt.Errorf("%w: %v", ErrorDuplicated, from)
	}

	if _, err := r.client.CreateRoute(ctx, r.createRoute(from, to)); err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}
	return nil
}
func (r *MailgunRouter) Unset(ctx context.Context, from string) error {
	route, err := r.findRoute(ctx, from)
	if err != nil {
		return fmt.Errorf("failed to find route: %w", err)
	}
	if route == nil {
		return fmt.Errorf("%w: %v", ErrorUndefined, from)
	}

	if err := r.client.DeleteRoute(ctx, route.Id); err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	return nil
}
