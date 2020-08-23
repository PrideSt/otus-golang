package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("next from empty", func(t *testing.T) {
		data := []Task{}

		stor := New(data)

		require.Nil(t, stor.Next())
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

		stor := New(data)

		{
			next := stor.Next()
			require.NotNil(t, next)
			require.Nil(t, next())
		}

		{
			next := stor.Next()
			require.NotNil(t, next)
			require.Equal(t, testingError, next())
		}

		{
			next := stor.Next()
			require.Nil(t, next)
		}
	})
}
