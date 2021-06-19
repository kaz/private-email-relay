package storage_test

import (
	"errors"
	"testing"

	"github.com/kaz/private-email-relay/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	t.Run("Firestore", func(t *testing.T) {
		s, err := storage.NewFirestoreStorage()
		assert.NoError(t, err)

		testSetAndGet(t, s)
	})
}

func testSetAndGet(t *testing.T, s storage.Storage) {
	key := "testSetAndGet.test"
	value := "test@test.test"

	err := s.Set(key, value)
	assert.NoError(t, err)

	got, err := s.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestSetAndUnset(t *testing.T) {
	t.Run("Firestore", func(t *testing.T) {
		s, err := storage.NewFirestoreStorage()
		assert.NoError(t, err)

		testSetAndUnset(t, s)
	})
}

func testSetAndUnset(t *testing.T, s storage.Storage) {
	key := "testSetAndUnset.test"
	value := "test@test.test"

	err := s.Set(key, value)
	assert.NoError(t, err)

	err = s.Unset(key)
	assert.NoError(t, err)

	_, err = s.Get(key)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, storage.ErrorNotFound))
}
