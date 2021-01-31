package gord

import (
	"container/list"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

type T = testing.T
type B = testing.B

var rnd = rand.New(rand.NewSource(0))

func ExampleOrdSet() {
	set := NewOrdSet()

	// Note the order.
	fmt.Println(set.Added(20)) // true
	fmt.Println(set.Added(10)) // true
	fmt.Println(set.Added(30)) // true

	// Redundant and doesn't change the order.
	fmt.Println(set.Added(30)) // false
	fmt.Println(set.Added(10)) // false
	fmt.Println(set.Added(20)) // false

	fmt.Println(set.Has(10))  // true
	fmt.Println(set.Has(40))  // false
	fmt.Println(set.Values()) // []interface{}{20, 10, 30}

	fmt.Println(set.PoppedFirst()) // 20, true
	fmt.Println(set.PoppedLast())  // 30, true
	fmt.Println(set.Values())      // []interface{}{10}

	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
	// true
	// false
	// [20 10 30]
	// 20 true
	// 30 true
	// [10]
}

// Other tests rely on this.
//
// Known issue: doesn't guarantee that `.Values()` will not accidentally
// deduplicate the values, if we have somehow stored duplicates.
func TestLinkedSetValues(t *T) {
	var ord list.List

	vals := map[interface{}]*list.Element{}
	vals[20] = ord.PushBack(20)
	vals[10] = ord.PushBack(10)
	vals[30] = ord.PushBack(30)

	set := LinkedSet{set: vals, ord: ord}

	requireEqual([]interface{}{20, 10, 30}, set.Values())
}

// Other tests rely on this.
func TestSliceSetValues(t *T) {
	requireEqual([]interface{}{20, 10, 30}, (&SliceSet{20, 10, 30}).Values())
}

func TestLinkedSet(t *T) {
	testSet(t, func() OrdSet { return new(LinkedSet) })

	t.Run("New", func(t *T) {
		testSetNew(func(args ...interface{}) OrdSet { return NewLinkedSet(args...) })
	})

	t.Run("String", func(t *T) {
		requireEqual(`[]`, (*LinkedSet)(nil).String())
		testSetString(new(LinkedSet))
	})

	t.Run("GoString", func(t *T) {
		requireEqual(`(*LinkedSet)(nil)`, (*LinkedSet)(nil).GoString())
		testLinkedSetGoString(new(LinkedSet), `NewLinkedSet`)
	})
}

func TestSyncLinkedSet(t *T) {
	testSet(t, func() OrdSet { return new(SyncLinkedSet) })

	t.Run("New", func(t *T) {
		testSetNew(func(args ...interface{}) OrdSet { return NewSyncLinkedSet(args...) })
	})

	t.Run("String", func(t *T) {
		requireEqual(`[]`, (*SyncLinkedSet)(nil).String())
		testSetString(new(SyncLinkedSet))
	})

	t.Run("GoString", func(t *T) {
		requireEqual(`(*SyncLinkedSet)(nil)`, (*SyncLinkedSet)(nil).GoString())
		testLinkedSetGoString(new(SyncLinkedSet), `NewSyncLinkedSet`)
	})
}

func TestSliceSet(t *T) {
	testSet(t, func() OrdSet { return new(SliceSet) })

	t.Run("New", func(t *T) {
		testSetNew(func(args ...interface{}) OrdSet { return NewSliceSet(args...) })
	})

	t.Run("String", func(t *T) {
		requireEqual(`[]`, (*SliceSet)(nil).String())
		var set SliceSet
		requireEqual(`[]`, set.String())
		testSetString(new(SliceSet))
	})

	t.Run("GoString", func(t *T) {
		requireEqual(`(*SliceSet)(nil)`, (*SliceSet)(nil).GoString())

		var set SliceSet
		requireEqual(`SliceSet(nil)`, set.GoString())

		set.Add(20)
		requireEqual(`SliceSet{20}`, set.GoString())

		set.Add(10)
		set.Add(30)
		requireEqual(`SliceSet{20, 10, 30}`, set.GoString())
	})
}

func testSet(t *T, newSet func() OrdSet) {
	t.Run("Len", func(t *T) { testSetLen(newSet()) })
	t.Run("Has", func(t *T) { testSetHas(newSet()) })
	t.Run("Add", func(t *T) { testSetAdd(newSet()) })
	t.Run("Added", func(t *T) { testSetAdded(newSet()) })
	t.Run("Delete", func(t *T) { testSetDelete(newSet()) })
	t.Run("Deleted", func(t *T) { testSetDeleted(newSet()) })
	t.Run("AddFirst", func(t *T) { testSetAddFirst(newSet()) })
	t.Run("AddedFirst", func(t *T) { testSetAddedFirst(newSet()) })
	t.Run("AddLast", func(t *T) { testSetAddLast(newSet()) })
	t.Run("AddedLast", func(t *T) { testSetAddedLast(newSet()) })
	t.Run("PoppedFirst", func(t *T) { testSetPoppedFirst(newSet()) })
	t.Run("PoppedLast", func(t *T) { testSetPoppedLast(newSet()) })
	t.Run("Values", func(t *T) { testSetValues(newSet()) })
}

// Relies on correctness of `Add`, which is tested later.
func testSetLen(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	requireEqual(1, set.Len())
	requireEqual([]interface{}{20}, set.Values())

	set.Add(20)
	requireEqual(1, set.Len())
	requireEqual([]interface{}{20}, set.Values())

	set.Add(10)
	requireEqual(2, set.Len())
	requireEqual([]interface{}{20, 10}, set.Values())

	set.Add(10)
	requireEqual(2, set.Len())
	requireEqual([]interface{}{20, 10}, set.Values())

	set.Add(30)
	requireEqual(3, set.Len())
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	set.Add(30)
	requireEqual(3, set.Len())
	requireEqual([]interface{}{20, 10, 30}, set.Values())
}

func testSetHas(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	requireEqual(true, set.Has(20))
	requireEqual(false, set.Has(10))
	requireEqual(false, set.Has(30))

	set.Add(20)
	requireEqual(true, set.Has(20))
	requireEqual(false, set.Has(10))
	requireEqual(false, set.Has(30))

	set.Add(10)
	requireEqual(true, set.Has(20))
	requireEqual(true, set.Has(10))
	requireEqual(false, set.Has(30))

	set.Add(10)
	requireEqual(true, set.Has(20))
	requireEqual(true, set.Has(10))
	requireEqual(false, set.Has(30))

	set.Add(30)
	requireEqual(true, set.Has(20))
	requireEqual(true, set.Has(10))
	requireEqual(true, set.Has(30))

	set.Add(30)
	requireEqual(true, set.Has(20))
	requireEqual(true, set.Has(10))
	requireEqual(true, set.Has(30))
}

func testSetAdd(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	requireEqual([]interface{}{20}, set.Values())

	set.Add(20)
	requireEqual([]interface{}{20}, set.Values())

	set.Add(10)
	requireEqual([]interface{}{20, 10}, set.Values())

	set.Add(10)
	requireEqual([]interface{}{20, 10}, set.Values())

	set.Add(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	set.Add(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())
}

func testSetAdded(set OrdSet) {
	requireEqual(0, set.Len())

	requireEqual(true, set.Added(20))
	requireEqual([]interface{}{20}, set.Values())

	requireEqual(false, set.Added(20))
	requireEqual([]interface{}{20}, set.Values())

	requireEqual(true, set.Added(10))
	requireEqual([]interface{}{20, 10}, set.Values())

	requireEqual(false, set.Added(10))
	requireEqual([]interface{}{20, 10}, set.Values())

	requireEqual(true, set.Added(30))
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	requireEqual(false, set.Added(30))
	requireEqual([]interface{}{20, 10, 30}, set.Values())
}

func testSetDelete(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	set.Add(10)
	set.Add(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	set.Delete(10)
	requireEqual([]interface{}{20, 30}, set.Values())

	set.Delete(10)
	requireEqual([]interface{}{20, 30}, set.Values())

	set.Delete(20)
	requireEqual([]interface{}{30}, set.Values())

	set.Delete(20)
	requireEqual([]interface{}{30}, set.Values())

	set.Delete(30)
	requireEqual([]interface{}{}, set.Values())

	set.Delete(30)
	requireEqual([]interface{}{}, set.Values())
}

func testSetDeleted(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	set.Add(10)
	set.Add(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	requireEqual(true, set.Deleted(10))
	requireEqual([]interface{}{20, 30}, set.Values())

	requireEqual(false, set.Deleted(10))
	requireEqual([]interface{}{20, 30}, set.Values())

	requireEqual(true, set.Deleted(20))
	requireEqual([]interface{}{30}, set.Values())

	requireEqual(false, set.Deleted(20))
	requireEqual([]interface{}{30}, set.Values())

	requireEqual(true, set.Deleted(30))
	requireEqual([]interface{}{}, set.Values())

	requireEqual(false, set.Deleted(30))
	requireEqual([]interface{}{}, set.Values())
}

func testSetAddFirst(set OrdSet) {
	requireEqual(0, set.Len())

	set.AddFirst(20)
	requireEqual([]interface{}{20}, set.Values())

	set.AddFirst(20)
	requireEqual([]interface{}{20}, set.Values())

	set.AddFirst(10)
	requireEqual([]interface{}{10, 20}, set.Values())

	set.AddFirst(10)
	requireEqual([]interface{}{10, 20}, set.Values())

	set.AddFirst(30)
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	set.AddFirst(30)
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	set.AddFirst(10)
	requireEqual([]interface{}{10, 30, 20}, set.Values())

	set.AddFirst(10)
	requireEqual([]interface{}{10, 30, 20}, set.Values())

	set.AddFirst(30)
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	set.AddFirst(30)
	requireEqual([]interface{}{30, 10, 20}, set.Values())
}

func testSetAddedFirst(set OrdSet) {
	requireEqual(0, set.Len())

	requireEqual(true, set.AddedFirst(20))
	requireEqual([]interface{}{20}, set.Values())

	requireEqual(false, set.AddedFirst(20))
	requireEqual([]interface{}{20}, set.Values())

	requireEqual(true, set.AddedFirst(10))
	requireEqual([]interface{}{10, 20}, set.Values())

	requireEqual(false, set.AddedFirst(10))
	requireEqual([]interface{}{10, 20}, set.Values())

	requireEqual(true, set.AddedFirst(30))
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	requireEqual(false, set.AddedFirst(30))
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	requireEqual(false, set.AddedFirst(10))
	requireEqual([]interface{}{10, 30, 20}, set.Values())

	requireEqual(false, set.AddedFirst(10))
	requireEqual([]interface{}{10, 30, 20}, set.Values())

	requireEqual(false, set.AddedFirst(30))
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	requireEqual(false, set.AddedFirst(30))
	requireEqual([]interface{}{30, 10, 20}, set.Values())
}

func testSetAddLast(set OrdSet) {
	requireEqual(0, set.Len())

	set.AddLast(20)
	requireEqual([]interface{}{20}, set.Values())

	set.AddLast(20)
	requireEqual([]interface{}{20}, set.Values())

	set.AddLast(10)
	requireEqual([]interface{}{20, 10}, set.Values())

	set.AddLast(10)
	requireEqual([]interface{}{20, 10}, set.Values())

	set.AddLast(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	set.AddLast(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	set.AddLast(10)
	requireEqual([]interface{}{20, 30, 10}, set.Values())

	set.AddLast(10)
	requireEqual([]interface{}{20, 30, 10}, set.Values())

	set.AddLast(20)
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	set.AddLast(20)
	requireEqual([]interface{}{30, 10, 20}, set.Values())
}

func testSetAddedLast(set OrdSet) {
	requireEqual(0, set.Len())

	requireEqual(true, set.AddedLast(20))
	requireEqual([]interface{}{20}, set.Values())

	requireEqual(false, set.AddedLast(20))
	requireEqual([]interface{}{20}, set.Values())

	requireEqual(true, set.AddedLast(10))
	requireEqual([]interface{}{20, 10}, set.Values())

	requireEqual(false, set.AddedLast(10))
	requireEqual([]interface{}{20, 10}, set.Values())

	requireEqual(true, set.AddedLast(30))
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	requireEqual(false, set.AddedLast(30))
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	requireEqual(false, set.AddedLast(10))
	requireEqual([]interface{}{20, 30, 10}, set.Values())

	requireEqual(false, set.AddedLast(10))
	requireEqual([]interface{}{20, 30, 10}, set.Values())

	requireEqual(false, set.AddedLast(20))
	requireEqual([]interface{}{30, 10, 20}, set.Values())

	requireEqual(false, set.AddedLast(20))
	requireEqual([]interface{}{30, 10, 20}, set.Values())
}

func testSetPoppedFirst(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	set.Add(10)
	set.Add(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	requireEqual(pair{20, true}, toPair(set.PoppedFirst()))
	requireEqual([]interface{}{10, 30}, set.Values())

	requireEqual(pair{10, true}, toPair(set.PoppedFirst()))
	requireEqual([]interface{}{30}, set.Values())

	requireEqual(pair{30, true}, toPair(set.PoppedFirst()))
	requireEqual([]interface{}{}, set.Values())

	requireEqual(pair{nil, false}, toPair(set.PoppedFirst()))
	requireEqual([]interface{}{}, set.Values())
}

func testSetPoppedLast(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	set.Add(10)
	set.Add(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())

	requireEqual(pair{30, true}, toPair(set.PoppedLast()))
	requireEqual([]interface{}{20, 10}, set.Values())

	requireEqual(pair{10, true}, toPair(set.PoppedLast()))
	requireEqual([]interface{}{20}, set.Values())

	requireEqual(pair{20, true}, toPair(set.PoppedLast()))
	requireEqual([]interface{}{}, set.Values())

	requireEqual(pair{nil, false}, toPair(set.PoppedLast()))
	requireEqual([]interface{}{}, set.Values())
}

func testSetValues(set OrdSet) {
	requireEqual(0, set.Len())

	set.Add(20)
	set.Add(10)
	set.Add(30)
	requireEqual([]interface{}{20, 10, 30}, set.Values())
}

func testSetNew(new func(...interface{}) OrdSet) {
	requireEqual(0, new().Len())
	requireEqual([]interface{}{20}, new(20).Values())
	requireEqual([]interface{}{20, 10}, new(20, 10).Values())
	requireEqual([]interface{}{20, 10, 30}, new(20, 10, 30).Values())
	requireEqual([]interface{}{20, 10, 30}, new(20, 10, 30, 20, 30, 10).Values())
}

func testSetString(set StringerOrdSet) {
	requireEqual(`[]`, set.String())

	set.Add(20)
	requireEqual(`[20]`, set.String())

	set.Add(10)
	set.Add(30)
	requireEqual(`[20 10 30]`, set.String())
}

// Applies to both `LinkedSet` and `SyncLinkedSet`.
func testLinkedSetGoString(set StringerOrdSet, constr string) {
	requireEqual(constr+`()`, set.GoString())

	set.Add(20)
	requireEqual(constr+`(20)`, set.GoString())

	set.Add(10)
	set.Add(30)
	requireEqual(constr+`(20, 10, 30)`, set.GoString())
}

func BenchmarkLinkedSet(b *B)     { bench(b, func() OrdSet { return new(LinkedSet) }) }
func BenchmarkSyncLinkedSet(b *B) { bench(b, func() OrdSet { return new(SyncLinkedSet) }) }
func BenchmarkSliceSet(b *B)      { bench(b, func() OrdSet { return new(SliceSet) }) }

func bench(b *B, newSet func() OrdSet) {
	b.Run("small", func(b *B) { benchSized(b, newSet, 1<<3) })
	b.Run("bigger", func(b *B) { benchSized(b, newSet, 1<<8) })
	b.Run("more bigger", func(b *B) { benchSized(b, newSet, 1<<12) })
	b.Run("MORE bigger", func(b *B) { benchSized(b, newSet, 1<<16) })
}

// Caution: this benchmark may be misleading. It combines MULTIPLE operations
// into one, testing the average, looking for worse-case outliers. Ideally, we
// would split this into individual benchmarks for different operations.
//
// This benchmark's main purpose is to demonstrate the degradation of `SliceSet`
// for large data sets, and identify the size after which it gets really bad.
func benchSized(b *B, newSet func() OrdSet, size int) {
	set, first, mid, last, next := benchInit(newSet, size)
	b.ResetTimer()

	for range counter(b.N) {
		set.Add(first)
		set.Add(mid)
		set.Add(last)
		set.Add(next)

		_ = set.Has(first)
		_ = set.Has(mid)
		_ = set.Has(last)
		_ = set.Has(next)

		set.AddFirst(first)
		set.AddFirst(mid)
		set.AddFirst(last)
		set.AddFirst(next)

		set.AddLast(first)
		set.AddLast(mid)
		set.AddLast(last)
		set.AddLast(next)

		set.Delete(first)
		set.Delete(mid)
		set.Delete(last)
		set.Delete(next)

		_, _ = set.PoppedFirst()
		_, _ = set.PoppedLast()
	}
}

func benchInit(
	newSet func() OrdSet, size int,
) (
	set OrdSet, first interface{}, mid interface{}, last interface{}, next interface{},
) {
	set = newSet()

	vals := rand.Perm(size)
	for _, val := range vals {
		set.Add(val)
	}

	first = vals[0]
	mid = vals[len(vals)/2]
	last = vals[len(vals)-1]
	next = rand.Int()
	return
}

func counter(n int) []struct{} { return make([]struct{}, n) }

// Note: the panic makes a stack trace, convenient for finding the line.
func requireEqual(expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		panic(fmt.Errorf(`
failed equality check
expected: %#v
actual:   %#v
`, expected, actual))
	}
}

func shuffled(list []int) []int {
	out := make([]int, len(list))
	copy(out, list)
	shuffle(out)
	return out
}

func shuffle(list []int) {
	rand.Shuffle(len(list), func(a, b int) {
		list[a], list[b] = list[b], list[a]
	})
}

func toPair(val interface{}, did bool) pair {
	return pair{val, did}
}

type pair struct {
	val interface{}
	did bool
}
