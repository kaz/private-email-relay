package storage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kaz/private-email-relay/internal/storage"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
)

func implements(t *testing.T) map[string]storage.Storage {
	var err error
	impl := map[string]storage.Storage{
		"memory": storage.NewMemoryStorage(),
	}

	impl["firestore"], err = storage.NewFirestoreStorage(ctx)
	if !assert.NoError(t, err) {
		delete(impl, "firestore")
	}

	return impl
}

func TestSetAndGet(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testSetAndGet(t, impl)
		})
	}
}
func testSetAndGet(t *testing.T, s storage.Storage) {
	key := "testSetAndGet.test"
	value := "testSetAndGet@test.test"

	err := s.Set(ctx, key, value, storage.NeverExpire)
	assert.NoError(t, err)

	got, err := s.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)

	// cleanup
	_, err = s.UnsetByKey(ctx, key)
	assert.NoError(t, err)
}

func TestUnsetByKey(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnsetByKey(t, impl)
		})
	}
}
func testUnsetByKey(t *testing.T, s storage.Storage) {
	key := "testUnsetByKey.test"
	value := "testUnsetByKey@test.test"

	err := s.Set(ctx, key, value, storage.NeverExpire)
	assert.NoError(t, err)

	deleted, err := s.UnsetByKey(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, deleted)

	_, err = s.Get(ctx, key)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorUndefinedKey))
}

func TestUnsetByValue(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnsetByValue(t, impl)
		})
	}
}
func testUnsetByValue(t *testing.T, s storage.Storage) {
	key := "testUnsetByKey.test"
	value := "testUnsetByValue@test.test"

	err := s.Set(ctx, key, value, storage.NeverExpire)
	assert.NoError(t, err)

	deleted, err := s.UnsetByValue(ctx, value)
	assert.NoError(t, err)
	assert.Equal(t, value, deleted)

	_, err = s.Get(ctx, key)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorUndefinedKey))
}

func TestUnsetExpired(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnsetExpired(t, impl)
		})
	}
}
func testUnsetExpired(t *testing.T, s storage.Storage) {
	now := time.Now()

	testCases := []struct {
		key     string
		value   string
		expires time.Time
		deleted bool
	}{
		{
			"testUnsetExpired-0.test",
			"testUnsetExpired-0@test.test",
			now.Add(-24 * time.Hour),
			true,
		},
		{
			"testUnsetExpired-1.test",
			"testUnsetExpired-1@test.test",
			now.Add(-1 * time.Hour),
			true,
		},
		{
			"testUnsetExpired-2.test",
			"testUnsetExpired-2@test.test",
			now,
			false,
		},
		{
			"testUnsetExpired-3.test",
			"testUnsetExpired-3@test.test",
			now.Add(1 * time.Hour),
			false,
		},
		{
			"testUnsetExpired-4.test",
			"testUnsetExpired-4@test.test",
			now.Add(24 * time.Hour),
			false,
		},
		{
			"testUnsetExpired-5.test",
			"testUnsetExpired-5@test.test",
			storage.NeverExpire,
			false,
		},
	}

	expectDeleted := []string{}
	for _, testCase := range testCases {
		err := s.Set(ctx, testCase.key, testCase.value, testCase.expires)
		assert.NoError(t, err)

		if testCase.deleted {
			expectDeleted = append(expectDeleted, testCase.value)
		}
	}

	actualDeleted, err := s.UnsetExpired(ctx, now)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectDeleted, actualDeleted)

	for _, testCase := range testCases {
		if testCase.deleted {
			_, err = s.Get(ctx, testCase.key)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, storage.ErrorUndefinedKey))
		} else {
			got, err := s.Get(ctx, testCase.key)
			assert.NoError(t, err)
			assert.Equal(t, testCase.value, got)

			// cleanup
			_, err = s.UnsetByKey(ctx, testCase.key)
			assert.NoError(t, err)
		}
	}
}

func TestGetUndefinedKey(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testGetUndefinedKey(t, impl)
		})
	}
}
func testGetUndefinedKey(t *testing.T, s storage.Storage) {
	key := "testGetUndefinedKey.test"

	_, err := s.Get(ctx, key)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorUndefinedKey))
}

func TestSetDuplicatedKey(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testSetDuplicatedKey(t, impl)
		})
	}
}
func testSetDuplicatedKey(t *testing.T, s storage.Storage) {
	key := "testSetDuplicatedKey.test"
	values := []string{
		"testSetDuplicatedKey-0@test.test",
		"testSetDuplicatedKey-1@test.test",
	}

	err := s.Set(ctx, key, values[0], storage.NeverExpire)
	assert.NoError(t, err)

	err = s.Set(ctx, key, values[1], storage.NeverExpire)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorDuplicatedKey))

	// cleanup
	_, err = s.UnsetByKey(ctx, key)
	assert.NoError(t, err)
}

func TestSetDuplicatedValue(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testSetDuplicatedValue(t, impl)
		})
	}
}
func testSetDuplicatedValue(t *testing.T, s storage.Storage) {
	keys := []string{
		"testSetDuplicatedValue-0.test",
		"testSetDuplicatedValue-1.test",
	}
	value := "testSetDuplicatedValue@test.test"

	err := s.Set(ctx, keys[0], value, storage.NeverExpire)
	assert.NoError(t, err)

	err = s.Set(ctx, keys[1], value, storage.NeverExpire)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorDuplicatedValue))

	// cleanup
	_, err = s.UnsetByKey(ctx, keys[0])
	assert.NoError(t, err)
}

func TestUnsetUndefinedKey(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnsetUndefinedKey(t, impl)
		})
	}
}
func testUnsetUndefinedKey(t *testing.T, s storage.Storage) {
	key := "testUnsetUndefinedKey.test"

	_, err := s.UnsetByKey(ctx, key)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorUndefinedKey))
}

func TestUnsetUndefinedValue(t *testing.T) {
	for name, impl := range implements(t) {
		t.Run(name, func(t *testing.T) {
			testUnsetUndefinedValue(t, impl)
		})
	}
}
func testUnsetUndefinedValue(t *testing.T, s storage.Storage) {
	value := "testUnsetUndefinedValue@test.test"

	_, err := s.UnsetByValue(ctx, value)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorUndefinedValue))
}
