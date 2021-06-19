package assign_test

import (
	"context"
	"testing"

	"github.com/kaz/private-email-relay/internal/assign"
	"github.com/kaz/private-email-relay/internal/router"
	"github.com/kaz/private-email-relay/internal/storage"
	"github.com/stretchr/testify/assert"
)

var (
	ctx   = context.Background()
	store = storage.NewMemoryStorage()
	route = router.NewMockRouter()
)

func implements(t *testing.T) map[string]assign.Strategy {
	var err error
	impl := map[string]assign.Strategy{}

	impl["default"], err = assign.NewDefaultStrategy(store, route)
	if !assert.NoError(t, err) {
		delete(impl, "default")
	}

	return impl
}

func TestAssignDifferentSite(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testAssignDifferentSite(t, impl)
		})
	}
}
func testAssignDifferentSite(t *testing.T, s assign.Strategy) {
	urls := []string{
		"https://www.youtube.com/watch?v=mZ0sJQC8qkE",
		"https://github.com/kaz/private-email-relay",
	}
	addrs := make([]string, len(urls))

	var err error

	addrs[0], err = s.Assign(ctx, urls[0])
	assert.NoError(t, err)

	addrs[1], err = s.Assign(ctx, urls[1])
	assert.NoError(t, err)

	assert.NotEqual(t, addrs[0], addrs[1])
}

func TestAssignExactlySameSite(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testAssignExactlySameSite(t, impl)
		})
	}
}

func testAssignExactlySameSite(t *testing.T, s assign.Strategy) {
	urls := []string{
		"https://www.youtube.com/watch?v=mZ0sJQC8qkE",
		"https://www.youtube.com/watch?v=i-b1lfCWGmc",
	}
	addrs := make([]string, len(urls))

	var err error

	addrs[0], err = s.Assign(ctx, urls[0])
	assert.NoError(t, err)

	addrs[1], err = s.Assign(ctx, urls[1])
	assert.NoError(t, err)

	assert.Equal(t, addrs[0], addrs[1])
}

func TestAssignEffectivelySameSite(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testAssignEffectivelySameSite(t, impl)
		})
	}
}
func testAssignEffectivelySameSite(t *testing.T, s assign.Strategy) {
	urls := []string{
		"https://www.youtube.com/watch?v=mZ0sJQC8qkE",
		"https://music.youtube.com/channel/UCuCfKSM0_23RRXxQGYTVJlw",
	}
	addrs := make([]string, len(urls))

	var err error

	addrs[0], err = s.Assign(ctx, urls[0])
	assert.NoError(t, err)

	addrs[1], err = s.Assign(ctx, urls[1])
	assert.NoError(t, err)

	assert.Equal(t, addrs[0], addrs[1])
}

func TestAssignConfusingDifferentSite(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testAssignConfusingDifferentSite(t, impl)
		})
	}
}
func testAssignConfusingDifferentSite(t *testing.T, s assign.Strategy) {
	urls := []string{
		"https://kaz.github.io",
		"https://sekai67.github.io",
	}
	addrs := make([]string, len(urls))

	var err error

	addrs[0], err = s.Assign(ctx, urls[0])
	assert.NoError(t, err)

	addrs[1], err = s.Assign(ctx, urls[1])
	assert.NoError(t, err)

	assert.NotEqual(t, addrs[0], addrs[1])
}

func TestUnassign(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnassign(t, impl)
		})
	}
}
func testUnassign(t *testing.T, s assign.Strategy) {
	addr, err := s.Assign(ctx, "https://kaz.github.io")
	assert.NoError(t, err)

	err = s.Unassign(ctx, addr)
	assert.NoError(t, err)
}

func TestUnassignUndefined(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnassignUndefined(t, impl)
		})
	}
}
func testUnassignUndefined(t *testing.T, s assign.Strategy) {
	err := s.Unassign(ctx, "unassigned@test.test")
	assert.Error(t, err)
}
