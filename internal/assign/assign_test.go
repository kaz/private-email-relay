package assign_test

import (
	"context"
	"testing"

	"github.com/kaz/private-email-relay/internal/assign"
	"github.com/kaz/private-email-relay/internal/storage"
	"github.com/stretchr/testify/assert"
)

var (
	ctx   = context.Background()
	store = storage.NewMemoryStorage()
)

func TestAssignDifferentSite(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		testAssignDifferentSite(t, assign.NewDefaultStrategy(store))
	})
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
	t.Run("Default", func(t *testing.T) {
		testAssignExactlySameSite(t, assign.NewDefaultStrategy(store))
	})
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
	t.Run("Default", func(t *testing.T) {
		testAssignEffectivelySameSite(t, assign.NewDefaultStrategy(store))
	})
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
	t.Run("Default", func(t *testing.T) {
		testAssignConfusingDifferentSite(t, assign.NewDefaultStrategy(store))
	})
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