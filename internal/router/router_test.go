package router_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kaz/private-email-relay/internal/router"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
)

func implements(t *testing.T) map[string]router.Router {
	var err error
	impl := map[string]router.Router{
		"mock": router.NewMockRouter(),
	}

	impl["mailgun"], err = router.NewMailgunRouter()
	if !assert.NoError(t, err) {
		delete(impl, "mailgun")
	}

	return impl
}

func TestSetAndUnset(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testSetAndUnset(t, impl)
		})
	}
}
func testSetAndUnset(t *testing.T, r router.Router) {
	from := "testSetAndUnset@test.test"
	to := "recipient@test.test"

	// Run 2 times to confirm an entry is successfully deleted
	for i := 0; i < 2; i++ {
		err := r.Set(ctx, from, to)
		assert.NoError(t, err)

		err = r.Unset(ctx, from)
		assert.NoError(t, err)
	}
}

func TestSetDuplicated(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testSetDuplicated(t, impl)
		})
	}
}
func testSetDuplicated(t *testing.T, r router.Router) {
	from := "testSetDuplicated@test.test"
	to := "recipient@test.test"

	err := r.Set(ctx, from, to)
	assert.NoError(t, err)

	err = r.Set(ctx, from, to)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, router.ErrorDuplicated))

	// cleanup
	err = r.Unset(ctx, from)
	assert.NoError(t, err)
}

func TestUnsetUndefined(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnsetUndefined(t, impl)
		})
	}
}
func testUnsetUndefined(t *testing.T, r router.Router) {
	from := "testUnsetUndefined@test.test"

	err := r.Unset(ctx, from)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, router.ErrorUndefined))
}
