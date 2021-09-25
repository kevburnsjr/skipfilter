package skipfilter

import (
	"fmt"
	"testing"
)

func TestSkipFilter(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var sf *SkipFilter
		t.Run("success", func(t *testing.T) {
			test := func(value interface{}, filter interface{}) bool {
				return true
			}
			for i, n := range []int{-1, 0, 10} {
				t.Run(fmt.Sprintf("size %d", n), func(t *testing.T) {
					sf = New(test, n)
					if sf == nil {
						t.Fatalf("New returned nil (%d)", i)
					}
				})
			}
		})
	})
	modTest := func(value interface{}, filter interface{}) bool {
		// value passes filter if value is multiple of filter
		return value.(int)%filter.(int) == 0
	}
	t.Run("Add", func(t *testing.T) {
		sf := New(modTest, 10)
		t.Run("success", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				sf.Add(i)
			}
		})
	})
	t.Run("Remove", func(t *testing.T) {
		sf := New(modTest, 10)
		t.Run("success", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				sf.Add(i)
			}
			sf.Remove(0)
			sf.Remove(11)
		})
	})
	t.Run("MatchAny", func(t *testing.T) {
		sf := New(modTest, 10)
		for i := 0; i < 10; i++ {
			sf.Add(i)
		}
		t.Run("success", func(t *testing.T) {
			res := sf.MatchAny(1)
			if len(res) != 10 {
				t.Fatalf("Expected 10 results, received (%d)", len(res))
			}
			res = sf.MatchAny(2)
			if len(res) != 5 {
				t.Fatalf("Expected 5 results, received (%d)", len(res))
			}
			res = sf.MatchAny(2)
			if len(res) != 5 {
				t.Fatalf("Expected 5 results, received (%d)", len(res))
			}
		})
		t.Run("removal", func(t *testing.T) {
			sf.Remove(0)
			res := sf.MatchAny(1)
			if len(res) != 9 {
				t.Fatalf("Expected 9 results, received (%d)", len(res))
			}
			res = sf.MatchAny(2)
			if len(res) != 4 {
				t.Fatalf("Expected 4 results, received (%d)", len(res))
			}
		})
		t.Run("remove all", func(t *testing.T) {
			for i := 1; i < 10; i++ {
				sf.Remove(i)
			}
			res := sf.MatchAny(1)
			if len(res) != 0 {
				t.Fatalf("Expected 0 results, received (%d)", len(res))
			}
			res = sf.MatchAny(2)
			if len(res) != 0 {
				t.Fatalf("Expected 0 results, received (%d)", len(res))
			}
		})
	})
	t.Run("Walk", func(t *testing.T) {
		sf := New(modTest, 10)
		for i := 0; i < 10; i++ {
			sf.Add(i)
		}
		t.Run("success", func(t *testing.T) {
			var n uint64
			id := sf.Walk(5, func(i interface{}) bool {
				n++
				return n < 5
			})
			if n != 5 {
				t.Fatalf("Expected WalkAll to execute callback 5 times, received (%d)", n)
			}
			if id != 10 {
				t.Fatalf(`Expected WalkAll to return cursor "10", received (%d)`, id)
			}
		})
		t.Run("removal", func(t *testing.T) {
			sf.Remove(0)
			var n uint64
			id := sf.Walk(5, func(i interface{}) bool {
				n++
				return n < 5
			})
			if n != 5 {
				t.Fatalf("Expected WalkAll to execute callback 5 times, received (%d)", n)
			}
			if id != 10 {
				t.Fatalf(`Expected WalkAll to return cursor "10", received (%d)`, id)
			}
		})
	})
	t.Run("Walk all", func(t *testing.T) {
		sf := New(modTest, 10)
		for i := 0; i < 10; i++ {
			sf.Add(i)
		}
		t.Run("success", func(t *testing.T) {
			var n uint64
			id := sf.Walk(0, func(i interface{}) bool {
				n++
				return true
			})
			if n != 10 {
				t.Fatalf("Expected WalkAll to execute callback 10 times, received (%d)", n)
			}
			if id != 10 {
				t.Fatalf(`Expected WalkAll to return cursor "10", received (%d)`, id)
			}
		})
		t.Run("removal", func(t *testing.T) {
			sf.Remove(0)
			var n uint64
			id := sf.Walk(0, func(i interface{}) bool {
				n++
				return true
			})
			if n != 9 {
				t.Fatalf("Expected WalkAll to execute callback 9 times, received (%d)", n)
			}
			if id != 10 {
				t.Fatalf(`Expected WalkAll to return cursor "10", received (%d)`, id)
			}
		})
	})
}
