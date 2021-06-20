package assign_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kaz/private-email-relay/internal/assign"
	"github.com/kaz/private-email-relay/internal/router"
	"github.com/kaz/private-email-relay/internal/storage"
	"github.com/stretchr/testify/assert"
)

var (
	ctx        = context.Background()
	implements = map[string]assign.Strategy{}
)

func TestMain(m *testing.M) {
	var err error

	store := storage.NewMemoryStorage()
	route := router.NewMockRouter()

	implements["default"], err = assign.NewDefaultStrategy(store, route)
	if err != nil {
		fmt.Printf("[[WARNING]] skip default: %v", err)
		delete(implements, "default")
	}

	m.Run()
}

func TestAssignDifferentSite(t *testing.T) {
	for name, impl := range implements {
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
	for name, impl := range implements {
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
	for name, impl := range implements {
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
	for name, impl := range implements {
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

func TestAssignTemporary(t *testing.T) {
	for name, impl := range implements {
		t.Run(name, func(t *testing.T) {
			testAssignTemporary(t, impl)
		})
	}
}
func testAssignTemporary(t *testing.T, s assign.Strategy) {
	url := "http://testAssignTemporary.test"

	addr1, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	addr2, err := s.AssignTemporary(ctx, url)
	assert.NoError(t, err)

	assert.NotEqual(t, addr1, addr2)
}

func TestUnassign(t *testing.T) {
	for name, impl := range implements {
		t.Run(name, func(t *testing.T) {
			testUnassign(t, impl)
		})
	}
}
func testUnassign(t *testing.T, s assign.Strategy) {
	url := "http://testUnassign.test"

	addr1, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	err = s.Unassign(ctx, url)
	assert.NoError(t, err)

	addr2, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	assert.NotEqual(t, addr1, addr2)
}

func TestUnassignByAddr(t *testing.T) {
	for name, impl := range implements {
		t.Run(name, func(t *testing.T) {
			testUnassignByAddr(t, impl)
		})
	}
}
func testUnassignByAddr(t *testing.T, s assign.Strategy) {
	url := "http://testUnassignByAddr.test"

	addr1, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	err = s.UnassignByAddr(ctx, addr1)
	assert.NoError(t, err)

	addr2, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	assert.NotEqual(t, addr1, addr2)
}

func TestUnassignTemporary(t *testing.T) {
	for name, impl := range implements {
		t.Run(name, func(t *testing.T) {
			testUnassignTemporary(t, impl)
		})
	}
}
func testUnassignTemporary(t *testing.T, s assign.Strategy) {
	url := "http://testUnassignTemporary.test"

	addr1a, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	addr1t, err := s.AssignTemporary(ctx, url)
	assert.NoError(t, err)

	err = s.Unassign(ctx, url)
	assert.NoError(t, err)

	addr2a, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	addr2t, err := s.AssignTemporary(ctx, url)
	assert.NoError(t, err)

	err = s.UnassignTemporary(ctx, url)
	assert.NoError(t, err)

	addr3a, err := s.Assign(ctx, url)
	assert.NoError(t, err)

	addr3t, err := s.AssignTemporary(ctx, url)
	assert.NoError(t, err)

	assert.NotEqual(t, addr1a, addr2a)
	assert.Equal(t, addr1t, addr2t)

	assert.Equal(t, addr2a, addr3a)
	assert.NotEqual(t, addr2t, addr3t)
}

func TestUnassignExpired(t *testing.T) {
	for name, impl := range implements {
		t.Run(name, func(t *testing.T) {
			testUnassignExpired(t, impl)
		})
	}
}
func testUnassignExpired(t *testing.T, s assign.Strategy) {
	urls := []string{
		"http://testUnassignExpired-0.test",
		"http://testUnassignExpired-1.test",
		"http://testUnassignExpired-2.test",
		"http://testUnassignExpired-3.test",
	}

	for _, url := range urls {
		_, err := s.AssignTemporary(ctx, url)
		assert.NoError(t, err)
	}

	count, err := s.UnassignExpired(ctx, time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// may delete entries created by other test
	count, err = s.UnassignExpired(ctx, time.Now().Add(365*24*time.Hour))
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, len(urls))
}

func TestUnassignUndefined(t *testing.T) {
	for name, impl := range implements {
		t.Run(name, func(t *testing.T) {
			testUnassignUndefined(t, impl)
		})
	}
}
func testUnassignUndefined(t *testing.T, s assign.Strategy) {
	url := "http://testUnassignUndefined.test"

	err := s.Unassign(ctx, url)
	assert.Error(t, err)
}

func TestUnassignByAddrUndefined(t *testing.T) {
	for name, impl := range implements {
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
	for name, impl := range implements {
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
	for name, impl := range implements {
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
