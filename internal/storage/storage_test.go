package storage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kaz/private-email-relay/internal/storage"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
)

func TestSetAndGet(t *testing.T) {
	t.Run("Memory", func(t *testing.T) {
		testSetAndGet(t, storage.NewMemoryStorage())
	})
	t.Run("Firestore", func(t *testing.T) {
		s, err := storage.NewFirestoreStorage(ctx)
		assert.NoError(t, err)

		testSetAndGet(t, s)
	})
}

func testSetAndGet(t *testing.T, s storage.Storage) {
	key := "testSetAndGet.test"
	value := "test@test.test"

	err := s.Set(ctx, key, value)
	assert.NoError(t, err)

	got, err := s.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestSetAndUnset(t *testing.T) {
	t.Run("Memory", func(t *testing.T) {
		testSetAndUnset(t, storage.NewMemoryStorage())
	})
	t.Run("Firestore", func(t *testing.T) {
		s, err := storage.NewFirestoreStorage(ctx)
		assert.NoError(t, err)

		testSetAndUnset(t, s)
	})
}

func testSetAndUnset(t *testing.T, s storage.Storage) {
	key := "testSetAndUnset.test"
	value := "test@test.test"

	err := s.Set(ctx, key, value)
	assert.NoError(t, err)

	err = s.Unset(ctx, key)
	assert.NoError(t, err)

	_, err = s.Get(ctx, key)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorNotFound))
}
