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

func TestUnassignByUrl(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnassignByUrl(t, impl)
		})
	}
}
func testUnassignByUrl(t *testing.T, s assign.Strategy) {
	url := "http://testUnassignByUrl.test"

	_, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	err = s.UnassignByUrl(ctx, url)
	assert.NoError(t, err)
}

func TestUnassignByAddr(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnassignByAddr(t, impl)
		})
	}
}
func testUnassignByAddr(t *testing.T, s assign.Strategy) {
	url := "http://testUnassignByAddr.test"

	addr, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	err = s.UnassignByAddr(ctx, addr)
	assert.NoError(t, err)
}

func TestUnassignByUrlUndefined(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnassignByUrlUndefined(t, impl)
		})
	}
}
func testUnassignByUrlUndefined(t *testing.T, s assign.Strategy) {
	url := "http://testUnassignByUrlUndefined.test"

	err := s.UnassignByUrl(ctx, url)
	assert.Error(t, err)
}

func TestUnassignByAddrUndefined(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnassignByAddrUndefined(t, impl)
		})
	}
}
func testUnassignByAddrUndefined(t *testing.T, s assign.Strategy) {
	addr := "testUnassignByAddrUndefined@test.test"

	err := s.UnassignByAddr(ctx, addr)
	assert.Error(t, err)
}

func TestAssignEdgeSuccess(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testAssignEdgeSuccess(t, impl)
		})
	}
}
func testAssignEdgeSuccess(t *testing.T, s assign.Strategy) {
	urls := []string{
		"//testAssignEdgeSuccess.test",
		"http://testAssignEdgeSuccess.test:8080",
		"http://127.0.0.1",
		"http://127.0.0.1:8080",
	}

	for _, url := range urls {
		_, err := s.Assign(ctx, url)
		assert.NoError(t, err)
	}
}

func TestAssignEdgeFail(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testAssignEdgeFail(t, impl)
		})
	}
}
func testAssignEdgeFail(t *testing.T, s assign.Strategy) {
	urls := []string{
		"/testAssignEdgeFail.test",
		"http://[::1]",
		"http://[::1]:8080",
	}

	for _, url := range urls {
		_, err := s.Assign(ctx, url)
		assert.Error(t, err)
	}
}
