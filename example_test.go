package skipfilter_test

import (
	"fmt"

	"github.com/kevburnsjr/skipfilter"
)

func ExampleNew() {
	sf := skipfilter.New(func(value interface{}, filter interface{}) bool {
		return value.(int)%filter.(int) == 0
	}, 10)
	fmt.Printf("%d", sf.Len())
	// Output:
	// 0
}

func ExampleSkipFilter_Add() {
	sf := skipfilter.New(func(value interface{}, filter interface{}) bool {
		return value.(int)%filter.(int) == 0
	}, 10)
	for i := 0; i < 10; i++ {
		sf.Add(i)
	}
	sf.Remove(0)
	fmt.Printf("%d", sf.Len())
	// Output:
	// 9
}

func ExampleSkipFilter_MatchAny() {
	sf := skipfilter.New(func(value interface{}, filter interface{}) bool {
		return value.(int)%filter.(int) == 0
	}, 10)
	for i := 0; i < 10; i++ {
		sf.Add(i)
	}
	fmt.Printf("Multiples of 2: %+v\n", sf.MatchAny(2))
	fmt.Printf("Multiples of 3: %+v\n", sf.MatchAny(3))
	// Output:
	// Multiples of 2: [0 2 4 6 8]
	// Multiples of 3: [0 3 6 9]
}

func ExampleSkipFilter_Walk_all() {
	sf := skipfilter.New(nil, 10)
	for i := 0; i < 10; i++ {
		sf.Add(i)
	}
	var n []int
	sf.Walk(0, func(v interface{}) bool {
		n = append(n, v.(int))
		return true
	})
	fmt.Printf("%d", len(n))
	// Output:
	// 10
}

func ExampleSkipFilter_Walk_limit() {
	sf := skipfilter.New(nil, 10)
	for i := 0; i < 10; i++ {
		sf.Add(i)
	}
	var n []int
	sf.Walk(0, func(v interface{}) bool {
		n = append(n, v.(int))
		return len(n) < 5
	})
	fmt.Printf("%d", len(n))
	// Output:
	// 5
}

func ExampleSkipFilter_Walk_start() {
	sf := skipfilter.New(nil, 10)
	for i := 0; i < 10; i++ {
		sf.Add(i)
	}
	var n []int
	sf.Walk(5, func(v interface{}) bool {
		n = append(n, v.(int))
		return true
	})
	fmt.Printf("%d", len(n))
	// Output:
	// 5
}
