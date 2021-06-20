package router_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/kaz/private-email-relay/internal/router"
	"github.com/stretchr/testify/assert"
)

var (
	ctx        = context.Background()
	implements = map[string]router.Router{}
)

func TestMain(m *testing.M) {
	var err error

	implements["mock"] = router.NewMockRouter()

	implements["mailgun"], err = router.NewMailgunRouter()
	if err != nil {
		fmt.Printf("[[WARNING]] skip mailgun: %v", err)
		delete(implements, "mailgun")
	}

	m.Run()
}

func TestSetAndUnset(t *testing.T) {
	for name, impl := range implements {
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
	for name, impl := range implements {
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
	for name, impl := range implements {
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
