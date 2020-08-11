package hw04_lru_cache //nolint:golint,stylecheck

import (
	golist "container/list"
	"testing"

	"github.com/stretchr/testify/require"
)

func listToSlice(l List) []int {
	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}

	return elems
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := listToSlice(l)
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestListExtra(t *testing.T) {
	t.Run("push front and remove back", func(t *testing.T) {
		l := NewList()

		l.PushFront(20) // [20]

		require.Equal(t, l.Len(), 1)
		require.Equal(t, l.Front(), l.Back())

		l.PushFront(10) // [10, 20]

		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value.(int), 10)
		require.Equal(t, l.Back().Value.(int), 20)

		{
			elems := listToSlice(l)
			require.Equal(t, []int{10, 20}, elems)
		}

		l.Remove(l.Back()) // [10]

		require.Equal(t, l.Len(), 1)
		require.Equal(t, l.Front(), l.Back())
		require.Equal(t, l.Front().Value.(int), 10)

		l.Remove(l.Back()) // []

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("push back and remove front", func(t *testing.T) {
		l := NewList()

		l.PushBack(10) // [10]

		require.Equal(t, l.Len(), 1)
		require.Equal(t, l.Front(), l.Back())

		l.PushBack(20) // [10, 20]

		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value.(int), 10)
		require.Equal(t, l.Back().Value.(int), 20)

		{
			elems := listToSlice(l)
			require.Equal(t, []int{10, 20}, elems)
		}

		l.Remove(l.Front()) // [20]

		require.Equal(t, l.Len(), 1)
		require.Equal(t, l.Front(), l.Back())
		require.Equal(t, l.Front().Value.(int), 20)

		l.Remove(l.Front()) // []

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("multy type test", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)      // [10]
		l.PushBack("hello") // [10, "hello"]

		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value.(int), 10)
		require.Equal(t, l.Back().Value.(string), "hello")
	})

	t.Run("move to front", func(t *testing.T) {
		l := NewList()

		l.PushBack(10) // [10]
		l.PushBack(20) // [10, 20]

		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value.(int), 10)
		require.Equal(t, l.Back().Value.(int), 20)

		{
			valBeforeMove := l.Back()
			l.MoveToFront(valBeforeMove) // [20, 10]
			valAferMove := l.Front()

			require.Equal(t, l.Len(), 2)
			require.Equal(t, l.Front().Value.(int), 20)
			require.Equal(t, l.Back().Value.(int), 10)
			require.Equal(t, valBeforeMove, valAferMove)
		}
		{
			middleItem := l.Back()    // take element 10 of [20, 10]
			l.PushBack(5)             // [20, 10, 5]
			l.MoveToFront(middleItem) // [10 20 5]

			require.Equal(t, l.Len(), 3)
			require.Equal(t, l.Front().Value.(int), 10) // check fitst
			require.Equal(t, l.Back().Value.(int), 5)   // check last

			// check element's links
			elems := listToSlice(l)
			require.Equal(t, []int{10, 20, 5}, elems)

		}
	})
}

// BenchmarkListPushFront/my_own_list-8         	12398272	       113 ns/op	      40 B/op	       2 allocs/op
// BenchmarkListPushFront/std_list-8            	 8667722	       145 ns/op	      56 B/op	       2 allocs/op
func BenchmarkListPushFront(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewList()
		for i := 0; i < t.N; i++ {
			l.PushFront(i)
		}
	})
	t.Run("std list", func(t *testing.B) {
		l := golist.New()
		for i := 0; i < t.N; i++ {
			l.PushFront(i)
		}
	})
}

// BenchmarkListNext/my_own_list-8         	1000000000	         0.00296 ns/op	       0 B/op	       0 allocs/op
// BenchmarkListNext/std_list-8            	1000000000	         0.00485 ns/op	       0 B/op	       0 allocs/op
func BenchmarkListNext(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewList()
		for i := 0; i < 1000000; i++ {
			l.PushFront(i)
		}

		t.ResetTimer()
		for i := l.Front(); i != nil; i = i.Next {
		}
	})
	t.Run("std list", func(t *testing.B) {
		l := golist.New()
		for i := 0; i < 1000000; i++ {
			l.PushFront(i)
		}

		t.ResetTimer()
		for i := l.Front(); i != nil; i = i.Next() {
		}
	})
}

// BenchmarkListRemoveLast/my_own_list-8         	1000000000	         0.00535 ns/op	       0 B/op	       0 allocs/op
// BenchmarkListRemoveLast/std_list-8            	1000000000	         0.000000 ns/op	       0 B/op	       0 allocs/op
func BenchmarkListRemoveLast(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewList()
		for i := 0; i < 1000000; i++ {
			l.PushFront(i)
		}

		t.ResetTimer()
		for i := l.Back(); i != nil; i = i.Prev {
			l.Remove(i)
		}
	})
	t.Run("std list", func(t *testing.B) {
		l := golist.New()
		for i := 0; i < 1000000; i++ {
			l.PushFront(i)
		}

		t.ResetTimer()
		for i := l.Back(); i != nil; i = i.Prev() {
			l.Remove(i)
		}
	})
}

// BenchmarkListMoveToFront/my_own_list-8         	299759928	         3.97 ns/op	       0 B/op	       0 allocs/op
// BenchmarkListMoveToFront/std_list-8            	199056926	         5.91 ns/op	       0 B/op	       0 allocs/op
func BenchmarkListMoveToFront(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewList()
		for i := 0; i < 1000; i++ {
			l.PushFront(i)
		}

		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			l.MoveToFront(l.Back())
		}
	})
	t.Run("std list", func(t *testing.B) {
		l := golist.New()
		for i := 0; i < 1000; i++ {
			l.PushFront(i)
		}

		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			l.MoveToFront(l.Back())
		}
	})
}
