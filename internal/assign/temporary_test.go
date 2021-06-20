package assign_test

import (
	"testing"
	"time"

	"github.com/kaz/private-email-relay/internal/assign"
	"github.com/stretchr/testify/assert"
)

func TestUnassignExpired(t *testing.T) {
	testUnassignExpired(t, implements["temporary"].(*assign.TemporaryStrategy))
}
func testUnassignExpired(t *testing.T, s *assign.TemporaryStrategy) {
	urls := []string{
		"http://testUnassignExpired-0.test",
		"http://testUnassignExpired-1.test",
		"http://testUnassignExpired-2.test",
		"http://testUnassignExpired-3.test",
	}

	for _, url := range urls {
		_, err := s.Assign(ctx, url)
		assert.NoError(t, err)
	}

	count, err := s.UnassignExpired(ctx, time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// may delete entries created by other test
	count, err = s.UnassignExpired(ctx, time.Now().Add(48*time.Hour))
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, len(urls))
}
