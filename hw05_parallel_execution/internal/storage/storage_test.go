package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("next from empty", func(t *testing.T) {
		store := New(nil)

		require.Nil(t, store.Next())
	})

	t.Run("not empty", func(t *testing.T) {
		f1 := func() error {
			return nil
		}

		testingError := errors.New("testing error")
		f2 := func() error {
			return testingError
		}
		data := []Task{f1, f2}

		store := New(data)

		{
			next := store.Next()
			require.NotNil(t, next)
			require.NoError(t, next())
		}

		{
			next := store.Next()
			require.NotNil(t, next)
			require.Equal(t, testingError, next())
		}

		{
			next := store.Next()
			require.Nil(t, next)
		}
	})
}
