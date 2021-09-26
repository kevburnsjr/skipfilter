package skipfilter

import (
	"fmt"
	"sync"

	"github.com/MauriceGit/skiplist"
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/hashicorp/golang-lru/simplelru"
)

type SkipFilter struct {
	i     uint64
	idx   map[interface{}]uint64
	set   *roaring64.Bitmap
	list  skiplist.SkipList
	cache *simplelru.LRU
	test  func(interface{}, interface{}) bool
	mutex sync.RWMutex
}

// New creates a new SkipFilter.
//   test - should return true if the value passes the provided filter.
//   size - controls the size of the LRU cache. Defaults to 100,000 if 0 or less.
//          should be tuned to match or exceed the expected filter cardinality.
func New(test func(value interface{}, filter interface{}) bool, size int) *SkipFilter {
	if size <= 0 {
		size = 1e5
	}
	cache, _ := simplelru.NewLRU(size, nil)
	return &SkipFilter{
		idx:   make(map[interface{}]uint64),
		set:   roaring64.New(),
		list:  skiplist.New(),
		cache: cache,
		test:  test,
	}
}

// Add adds a value to the set
func (sf *SkipFilter) Add(value interface{}) {
	el := &entry{sf.i, value}
	sf.list.Insert(el)
	sf.set.Add(sf.i)
	sf.idx[value] = sf.i
	sf.i++
}

// Remove removes a value from the set
func (sf *SkipFilter) Remove(value interface{}) {
	if id, ok := sf.idx[value]; ok {
		sf.list.Delete(&entry{id: id})
		sf.set.Remove(id)
		delete(sf.idx, value)
	}
}

// Len returns the number of values in the set
func (sf *SkipFilter) Len() int {
	return sf.list.GetNodeCount()
}

// MatchAny returns a slice of values in the set matching any of the provided filters
func (sf *SkipFilter) MatchAny(filters ...interface{}) []interface{} {
	var set = roaring64.New()
	var fs = make([]*filter, len(filters))
	for i, k := range filters {
		fs[i] = sf.getFilter(k)
		set.Or(fs[i].set)
	}
	values, notfound := sf.getValues(set)
	for _, id := range notfound {
		// Clean up references to removed values
		for _, f := range fs {
			f.set.Remove(id)
		}
	}
	return values
}

// Walk executes callback for each value in the set beginning at `start` index.
// Return true in callback to continue iterating, false to stop.
// Returned uint64 is index of `next` element (send as `start` to continue iterating)
func (sf *SkipFilter) Walk(start uint64, callback func(val interface{}) bool) uint64 {
	var i uint64
	var id = start
	var prev uint64
	var first = true
	el, ok := sf.list.FindGreaterOrEqual(&entry{id: start})
	for ok && el != nil {
		if id = el.GetValue().(*entry).id; !first && id <= prev {
			// skiplist loops back to first element so we have to detect loop and break manually
			id = prev + 1
			break
		}
		i++
		if !callback(el.GetValue().(*entry).val) {
			id++
			break
		}
		prev = id
		el = sf.list.Next(el)
		first = false
	}
	return id
}

func (sf *SkipFilter) getFilter(k interface{}) *filter {
	var f *filter
	val, ok := sf.cache.Get(k)
	if ok {
		f = val.(*filter)
	} else {
		f = &filter{0, roaring64.New()}
		sf.cache.Add(k, f)
	}
	var id uint64
	var prev uint64
	var first = true
	for el, ok := sf.list.FindGreaterOrEqual(&entry{id: f.i}); ok && el != nil; el = sf.list.Next(el) {
		if id = el.GetValue().(*entry).id; !first && id <= prev {
			// skiplist loops back to first element so we have to detect loop and break manually
			break
		}
		if sf.test(el.GetValue().(*entry).val, k) {
			f.set.Add(id)
		}
		prev = id
		first = false
	}
	f.i = sf.i
	return f
}

func (sf *SkipFilter) getValues(set *roaring64.Bitmap) ([]interface{}, []uint64) {
	idBuf := make([]uint64, 512)
	iter := set.ManyIterator()
	values := []interface{}{}
	notfound := []uint64{}
	e := &entry{}
	for n := iter.NextMany(idBuf); n > 0; n = iter.NextMany(idBuf) {
		for i := 0; i < n; i++ {
			e.id = idBuf[i]
			el, ok := sf.list.Find(e)
			if ok {
				values = append(values, el.GetValue().(*entry).val)
			} else {
				notfound = append(notfound, idBuf[i])
			}
		}
	}
	return values, notfound
}

type entry struct {
	id  uint64
	val interface{}
}

func (e *entry) ExtractKey() float64 {
	return float64(e.id)
}

func (e *entry) String() string {
	return fmt.Sprintf("%16x", e.id)
}

type filter struct {
	i   uint64
	set *roaring64.Bitmap
}
