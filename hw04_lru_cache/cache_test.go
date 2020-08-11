package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("repeat set", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("aaa", 100)
		require.True(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)
	})

	t.Run("change set", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		wasInCache = c.Set("aaa", 200)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 200, val)
	})

	t.Run("set moves to first", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100) // [aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [bbb, aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("aaa", 150) // aaa become first, [aaa, bbb]
		require.True(t, wasInCache)

		wasInCache = c.Set("ccc", 300) // ccc overwrite bbb and come the first, [ccc, aaa]
		require.False(t, wasInCache)

		val, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("get moves to first", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100) // [aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [bbb, aaa]
		require.False(t, wasInCache)

		val, ok := c.Get("aaa") // aaa become first, [aaa, bbb]
		require.True(t, ok)
		require.Equal(t, 100, val)

		wasInCache = c.Set("ccc", 300) // ccc overwrite bbb and come the first [ccc, aaa]
		require.False(t, wasInCache)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("add overflowed item again", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100) // [aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [bbb, aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300) // ccc overwrite aaa and come the first [ccc, bbb]
		require.False(t, wasInCache)

		wasInCache = c.Set("aaa", 100) // [aaa, ccc]
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("ccc") // aaa become first, [aaa, bbb]
		require.True(t, ok)
		require.Equal(t, 300, val)
	})
}

func TestClear(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		c.Clear()

		val, ok = c.Get("aaa")
		require.False(t, ok)

		val, ok = c.Get("bbb")
		require.False(t, ok)

		wasInCache = c.Set("aaa", 100)
		require.False(t, wasInCache)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
