/*
Simple ordered sets. Several implementations are available:

• `LinkedSet`: ordered set with near-constant-time (O(1)) performance for
inserting, deleting, and moving elements. Backed by a map and a doubly-linked
list.

• `SyncLinkedSet`: concurrency-safe `LinkedSet`, slightly slower.

• `SliceSet`: slice-backed ordered set. Simpler and faster for small sets,
extreme performance degradation for large sets.

Installation

Simply import:

	import "github.com/mitranim/gord"

	set := NewOrdSet()

Examples

See the example below for `OrdSet` / `NewOrdSet`.
*/
package gord

import "fmt"

// Interface that describes an arbitrary set, but not necessarily an ordered
// set. Satisfied by every type in this package. See `OrdSet` for the full
// interface.
type Set interface {
	// Current set size, replacement for `len(set)`.
	Len() int

	// Answers whether the value is in the set.
	Has(val interface{}) bool

	// Void version of `.Added`.
	Add(val interface{})

	// If `set.Has(val)`, has no effect and returns `false`.
	// If `!set.Has(val)`, appends `val` to the end and returns `true`.
	// In ordered sets (every type in this package), must not change the order.
	Added(val interface{}) bool

	// Void version of `.Deleted`.
	Delete(val interface{})

	// If `set.Has(val)`, deletes the value and returns `true`.
	// If `!set.Has(val)`, does nothing and returns `false`.
	Deleted(val interface{}) bool
}

// Interface that describes an ordered set. Strict superset of `Set`. Satisfied
// by every type in this package.
type OrdSet interface {
	Set

	// Void version of `.AddedFirst`.
	AddFirst(val interface{})

	// If `set.Has(val)`, moves the value to the first position and returns `false`.
	// If `!set.Has(val)`, prepends the value at the start and returns `true`.
	AddedFirst(val interface{}) bool

	// Void version of `.AddedLast`.
	AddLast(val interface{})

	// If `set.Has(val)`, moves the value to the last position and returns `false`.
	// If `!set.Has(val)`, appends the value at the end and returns `true`.
	AddedLast(val interface{}) bool

	// If the set is empty, returns `(nil, false)`.
	// Otherwise removes the first value and returns `(val, true)`.
	//
	// There's no corresponding `.PopFirst` (without the boolean) to avoid the
	// potential gotcha where the first value was `nil`, but the code would
	// erroneously think that the set was empty.
	PoppedFirst() (interface{}, bool)

	// If the set is empty, returns `(nil, false)`.
	// Otherwise removes the last value and returns `(val, true)`.
	//
	// There's no corresponding `.PopLast` (without the boolean) to avoid the
	// potential gotcha where the last value was `nil`, but the code would
	// erroneously think that the set was empty.
	PoppedLast() (interface{}, bool)

	// Returns the set's values as a slice, in the same order. Allowed to return
	// either `nil` or `[]interface{}{}`. Callers are expected to not care about
	// the distinction, or the slice's capacity. Whether the callers are allowed
	// to mutate the slice depends on the implementation of the backing type.
	Values() []interface{}
}

// Describes a set with extra printing methods. Satisfied by every type in this
// package.
type StringerSet interface {
	Set
	fmt.Stringer
	fmt.GoStringer
}

// Describes an ordered set with extra printing methods. Satisfied by every type
// in this package.
type StringerOrdSet interface {
	OrdSet
	fmt.Stringer
	fmt.GoStringer
}

// Default way to create an ordered set, using `LinkedSet`.
func NewOrdSet(vals ...interface{}) OrdSet {
	return NewLinkedSet(vals...)
}
