package hw04_lru_cache //nolint:golint,stylecheck

import (
	golist "container/list"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSyncList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewSyncList(NewList())

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewSyncList(NewList())

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, l.Len(), 3)

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, l.Len(), 2)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := listToSlice(l)
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestSyncListExtra(t *testing.T) {
	t.Run("push front and remove back", func(t *testing.T) {
		l := NewSyncList(NewList())

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
		l := NewSyncList(NewList())

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
}

func TestSyncListConcurrencyPushBack(t *testing.T) {
	t.Run("concurrency push", func(t *testing.T) {
		l := NewSyncList(NewList())

		data := make([]int, 1000000)
		for i := range data {
			data[i] = i
		}

		fpush := func(wg *sync.WaitGroup, l List, data []int) {
			defer wg.Done()
			for _, v := range data {
				l.PushBack(v)
			}
		}

		wg := sync.WaitGroup{}
		threadsCnt := 5
		batchSize := len(data) / threadsCnt
		for i := 0; i < threadsCnt; i++ {
			wg.Add(1)
			go fpush(&wg, l, data[batchSize*i:batchSize*(i+1)])
		}

		wg.Wait()

		listSlice := listToSlice(l)
		sort.Sort(sort.IntSlice(listSlice))

		require.Equal(t, data, listSlice)
	})
}

// BenchmarkListPushFront/my_own_list-8 (no-sync)	  	12398272	       113 ns/op	      40 B/op	       2 allocs/op
// BenchmarkSyncListPushFront/my_own_list-8         	10443217	       114 ns/op	      40 B/op	       2 allocs/op
// BenchmarkSyncListPushFront/std_list-8            	 9038499	       158 ns/op	      56 B/op	       2 allocs/op
func BenchmarkSyncListPushFront(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewSyncList(NewList())
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

// BenchmarkListNext/my_own_list-8 (no-sync)   	1000000000	         0.00296 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSyncListNext/my_own_list-8         	1000000000	         0.00300 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSyncListNext/std_list-8            	1000000000	         0.00478 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSyncListNext(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewSyncList(NewList())
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

// BenchmarkListRemoveLast/my_own_list-8 (no-sync)     	1000000000	         0.00535 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSyncListRemoveLast/my_own_list-8         	1000000000	         0.0353 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSyncListRemoveLast/std_list-8            	1000000000	         0.000001 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSyncListRemoveLast(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewSyncList(NewList())
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

// BenchmarkListMoveToFront/my_own_list-8 (no-sync)        	299759928	         3.97 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSyncListMoveToFront/my_own_list-8         	12725692	        93.2 ns/op	      48 B/op	       1 allocs/op
// BenchmarkSyncListMoveToFront/std_list-8            	191637933	         6.04 ns/op	       0 B/op	       0 allocs/op

// Не понимаю откуда появился 1 allocs/op ???
func BenchmarkSyncListMoveToFront(t *testing.B) {
	t.Run("my own list", func(t *testing.B) {
		l := NewSyncList(NewList())
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
