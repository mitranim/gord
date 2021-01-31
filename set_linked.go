package gord

import (
	"container/list"
	"fmt"
	"strings"
)

// Constructs a new `LinkedSet` from the provided values, deduplicating them.
func NewLinkedSet(vals ...interface{}) *LinkedSet {
	var set LinkedSet
	for _, val := range vals {
		set.Add(val)
	}
	return &set
}

// Ordered set. Satisfies the `OrdSet` interface. A zero value is ready to use,
// but should not be copied after the first mutation. Has near-constant-time
// (O(1)) performance for inserting, deleting, and moving elements.
//
// Concurrency-unsafe; use `SyncLinkedSet` for concurrent access.
type LinkedSet struct {
	set map[interface{}]*list.Element
	ord list.List
}

// Satisfy `Set`.
func (self *LinkedSet) Len() int {
	if self == nil {
		return 0
	}
	return len(self.set)
}

// Satisfy `Set`.
func (self *LinkedSet) Has(val interface{}) bool {
	if self == nil {
		return false
	}
	_, ok := self.set[val]
	return ok
}

// Satisfy `Set`.
func (self *LinkedSet) Add(val interface{}) {
	_ = self.Added(val)
}

// Satisfy `Set`.
func (self *LinkedSet) Added(val interface{}) bool {
	if self.Has(val) {
		return false
	}

	self.init()
	self.set[val] = self.ord.PushBack(val)
	return true
}

// Satisfy `Set`.
func (self *LinkedSet) Delete(val interface{}) {
	_ = self.Deleted(val)
}

// Satisfy `Set`.
func (self *LinkedSet) Deleted(val interface{}) bool {
	elem := self.set[val]
	if elem == nil {
		return false
	}
	self.removeElem(elem)
	return true
}

// Satisfy `OrdSet`.
func (self *LinkedSet) AddFirst(val interface{}) {
	_ = self.AddedFirst(val)
}

// Satisfy `OrdSet`.
func (self *LinkedSet) AddedFirst(val interface{}) bool {
	self.init()

	elem := self.set[val]
	if elem != nil {
		self.ord.MoveToFront(elem)
		return false
	}

	self.set[val] = self.ord.PushFront(val)
	return true
}

// Satisfy `OrdSet`.
func (self *LinkedSet) AddLast(val interface{}) {
	_ = self.AddedLast(val)
}

// Satisfy `OrdSet`.
func (self *LinkedSet) AddedLast(val interface{}) bool {
	self.init()

	elem := self.set[val]
	if elem != nil {
		self.ord.MoveToBack(elem)
		return false
	}

	self.set[val] = self.ord.PushBack(val)
	return true
}

// Satisfy `OrdSet`.
func (self *LinkedSet) PoppedFirst() (interface{}, bool) {
	return self.poppedElem(self.ord.Front())
}

// Satisfy `OrdSet`.
func (self *LinkedSet) PoppedLast() (interface{}, bool) {
	return self.poppedElem(self.ord.Back())
}

// Satisfy `OrdSet`. The slice is allocated every time and is OK to mutate.
func (self *LinkedSet) Values() []interface{} {
	if self == nil {
		return nil
	}

	out := make([]interface{}, 0, self.ord.Len())
	self.each(func(_ int, val interface{}) {
		out = append(out, val)
	})
	return out
}

// Satisfy `StringerOrdSet`.
func (self *LinkedSet) String() string {
	if self == nil {
		return `[]`
	}
	// Inefficient but simple.
	return fmt.Sprint(self.Values())
}

// Satisfy `StringerOrdSet`.
func (self *LinkedSet) GoString() string {
	if self == nil {
		return `(*LinkedSet)(nil)`
	}

	var buf strings.Builder
	buf.WriteString(`NewLinkedSet`)
	self.writeGoString(&buf)
	return buf.String()
}

func (self *LinkedSet) writeGoString(buf *strings.Builder) {
	buf.WriteString(`(`)
	self.each(func(i int, val interface{}) {
		if i > 0 {
			buf.WriteString(`, `)
		}
		fmt.Fprintf(buf, "%#v", val)
	})
	buf.WriteString(`)`)
}

func (self *LinkedSet) init() {
	if self.set == nil {
		self.set = map[interface{}]*list.Element{}
	}
}

func (self *LinkedSet) removeElem(elem *list.Element) {
	self.ord.Remove(elem)
	delete(self.set, elem.Value)
}

func (self *LinkedSet) poppedElem(elem *list.Element) (interface{}, bool) {
	if elem == nil {
		return nil, false
	}
	self.removeElem(elem)
	return elem.Value, true
}

func (self LinkedSet) each(fun func(i int, val interface{})) {
	i := 0
	for elem := self.ord.Front(); elem != nil; elem = elem.Next() {
		fun(i, elem.Value)
		i++
	}
}
