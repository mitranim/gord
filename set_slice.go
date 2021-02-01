package gord

import (
	"fmt"
	"strings"
)

// Constructs a new `SliceSet` from the provided values. Very similar to
// `SliceSet{}` or `&SliceSet{}`, but discards duplicates.
func NewSliceSet(vals ...interface{}) *SliceSet {
	var set SliceSet
	for _, val := range vals {
		set.Add(val)
	}
	return &set
}

// Ordered set implemented as a slice. Compared to `LinkedSet`, this is simpler
// and more efficient (memory and CPU wise) for small sets, but some operations
// have extreme performance degradation for large sets. Giving the exact values
// of "small" and "large" requires more testing; might depend on hardware.
//
// There's no "concurrent" version of this type, mainly because it would have to
// sacrifice elegance. It may be added on demand.
type SliceSet []interface{}

func (self *SliceSet) Len() int {
	if self == nil {
		return 0
	}
	return len(*self)
}

// Satisfy `Set`.
func (self *SliceSet) Has(val interface{}) bool {
	if self == nil {
		return false
	}
	for _, value := range *self {
		if value == val {
			return true
		}
	}
	return false
}

// Satisfy `Set`.
func (self *SliceSet) Add(val interface{}) {
	_ = self.Added(val)
}

// Satisfy `Set`.
func (self *SliceSet) Added(val interface{}) bool {
	if !self.Has(val) {
		*self = append(*self, val)
		return true
	}
	return false
}

// Satisfy `Set`.
func (self *SliceSet) Delete(val interface{}) {
	_ = self.Deleted(val)
}

// Satisfy `Set`.
func (self *SliceSet) Deleted(val interface{}) bool {
	slice := *self

	for i, value := range slice {
		if value == val {
			slice.shiftLeft(i)
			*self = slice.init()
			return true
		}
	}

	return false
}

// Satisfy `OrdSet`.
func (self *SliceSet) AddFirst(val interface{}) {
	_ = self.AddedFirst(val)
}

// Satisfy `OrdSet`.
func (self *SliceSet) AddedFirst(val interface{}) bool {
	slice := *self

	for i, value := range slice {
		if value == val {
			slice.shiftRight(i)
			slice[0] = val
			return false
		}
	}

	var longer SliceSet
	if cap(slice) > len(slice) {
		longer = slice[:len(slice)+1]
	} else {
		longer = make([]interface{}, len(slice)+1, moreCap(len(slice)))
	}

	copy(longer[1:], slice)
	longer[0] = val
	*self = longer
	return true
}

// Satisfy `OrdSet`.
func (self *SliceSet) AddLast(val interface{}) {
	_ = self.AddedLast(val)
}

// Satisfy `OrdSet`.
func (self *SliceSet) AddedLast(val interface{}) bool {
	slice := *self

	for i, value := range slice {
		if value == val {
			slice.shiftLeft(i)
			slice[len(slice)-1] = val
			return false
		}
	}

	*self = append(*self, val)
	return true
}

// Satisfy `OrdSet`.
func (self *SliceSet) PoppedFirst() (interface{}, bool) {
	slice := *self

	if len(slice) > 0 {
		val := slice[0]
		slice.shiftLeft(0)
		*self = slice.init()
		return val, true
	}

	return nil, false
}

// Satisfy `OrdSet`.
func (self *SliceSet) PoppedLast() (interface{}, bool) {
	slice := *self

	if len(slice) > 0 {
		val := slice[len(slice)-1]
		*self = slice.init()
		return val, true
	}

	return nil, false
}

// Satisfy `OrdSet`. Returns self. This is a free cast with no reallocation, and
// any mutations of the resulting slice are reflected in the set.
func (self *SliceSet) Values() []interface{} {
	if self == nil {
		return nil
	}
	return []interface{}(*self)
}

// Satisfy `StringerOrdSet`.
func (self *SliceSet) String() string {
	return fmt.Sprint(self.Values())
}

// Satisfy `StringerOrdSet`.
func (self *SliceSet) GoString() string {
	if self == nil {
		return `(*SliceSet)(nil)`
	}

	if *self == nil {
		return `SliceSet(nil)`
	}

	var buf strings.Builder
	buf.WriteString(`SliceSet{`)

	for i, val := range *self {
		if i > 0 {
			buf.WriteString(`, `)
		}
		fmt.Fprintf(&buf, "%#v", val)
	}

	buf.WriteString(`}`)
	return buf.String()
}

func (self SliceSet) shiftLeft(index int) {
	copy(self[index:], self[index+1:])
}

func (self SliceSet) shiftRight(index int) {
	copy(self[1:index+1], self[:index])
}

func (self SliceSet) init() SliceSet {
	return self[:len(self)-1]
}

// At the time of writing, this is a reasonable approximation of how the Go
// runtime increases slice capacity when appending. The exact amount may vary
// by type and Go version.
func moreCap(size int) int {
	if size > 0 {
		return size * 2
	}
	return 1
}
