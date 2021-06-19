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

func TestSetAndUnset(t *testing.T) {
	t.Run("Mock", func(t *testing.T) {
		testSetAndUnset(t, router.NewMockRouter())
	})
	t.Run("Mailgun", func(t *testing.T) {
		r, err := router.NewMailgunRouter()
		assert.NoError(t, err)

		testSetAndUnset(t, r)
	})
}

func testSetAndUnset(t *testing.T, r router.Router) {
	from := "testSetAndUnset@test.test"
	to := "dummy@test.test"

	// Run 2 times to confirm an entry is successfully deleted
	for i := 0; i < 2; i++ {
		err := r.Set(ctx, from, to)
		assert.NoError(t, err)

		err = r.Unset(ctx, from)
		assert.NoError(t, err)
	}
}

func TestDuplicateSet(t *testing.T) {
	t.Run("Mock", func(t *testing.T) {
		testDuplicateSet(t, router.NewMockRouter())
	})
	t.Run("Mailgun", func(t *testing.T) {
		r, err := router.NewMailgunRouter()
		assert.NoError(t, err)

		testDuplicateSet(t, r)
	})
}

func testDuplicateSet(t *testing.T, r router.Router) {
	from := "testDuplicateSet@test.test"
	to := "dummy@test.test"

	err := r.Set(ctx, from, to)
	assert.NoError(t, err)

	err = r.Set(ctx, from, to)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, router.ErrorDuplicated))

	// cleanup
	err = r.Unset(ctx, from)
	assert.NoError(t, err)
}

func TestUnsetNonexistent(t *testing.T) {
	t.Run("Mock", func(t *testing.T) {
		testUnsetNonexistent(t, router.NewMockRouter())
	})
	t.Run("Mailgun", func(t *testing.T) {
		r, err := router.NewMailgunRouter()
		assert.NoError(t, err)

		testUnsetNonexistent(t, r)
	})
}

func testUnsetNonexistent(t *testing.T, r router.Router) {
	from := "testUnsetNonexistent@test.test"

	err := r.Unset(ctx, from)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, router.ErrorUnsetNonexistent))
}
