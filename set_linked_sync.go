package gord

import (
	"strings"
	"sync"
)

// Constructs a new `SyncLinkedSet` from the provided values, deduplicating them.
func NewSyncLinkedSet(vals ...interface{}) *SyncLinkedSet {
	var set SyncLinkedSet
	for _, val := range vals {
		set.set.Add(val)
	}
	return &set
}

// Concurrency-safe, slightly slower version of `LinkedSet`. Satisfies the
// `OrdSet` interface. A zero value is ready to use, but should never be
// copied. Uses a mutex.
type SyncLinkedSet struct {
	lock sync.Mutex
	set  LinkedSet
}

// Concurrency-safe version of `LinkedSet.Len`.
func (self *SyncLinkedSet) Len() int {
	if self == nil {
		return 0
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.Len()
}

// Concurrency-safe version of `LinkedSet.Has`.
func (self *SyncLinkedSet) Has(val interface{}) bool {
	if self == nil {
		return false
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.Has(val)
}

// Concurrency-safe version of `LinkedSet.Add`.
func (self *SyncLinkedSet) Add(val interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.set.Add(val)
}

// Concurrency-safe version of `LinkedSet.Added`.
func (self *SyncLinkedSet) Added(val interface{}) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.Added(val)
}

// Concurrency-safe version of `LinkedSet.Delete`.
func (self *SyncLinkedSet) Delete(val interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.set.Delete(val)
}

// Concurrency-safe version of `LinkedSet.Deleted`.
func (self *SyncLinkedSet) Deleted(val interface{}) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.Deleted(val)
}

// Concurrency-safe version of `LinkedSet.AddFirst`.
func (self *SyncLinkedSet) AddFirst(val interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.set.AddFirst(val)
}

// Concurrency-safe version of `LinkedSet.AddedFirst`.
func (self *SyncLinkedSet) AddedFirst(val interface{}) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.AddedFirst(val)
}

// Concurrency-safe version of `LinkedSet.AddLast`.
func (self *SyncLinkedSet) AddLast(val interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.set.AddLast(val)
}

// Concurrency-safe version of `LinkedSet.AddedLast`.
func (self *SyncLinkedSet) AddedLast(val interface{}) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.AddedLast(val)
}

// Concurrency-safe version of `LinkedSet.PoppedFirst`.
func (self *SyncLinkedSet) PoppedFirst() (interface{}, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.PoppedFirst()
}

// Concurrency-safe version of `LinkedSet.PoppedLast`.
func (self *SyncLinkedSet) PoppedLast() (interface{}, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.PoppedLast()
}

// Concurrency-safe version of `LinkedSet.Values`.
func (self *SyncLinkedSet) Values() []interface{} {
	if self == nil {
		return nil
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.Values()
}

// Concurrency-safe version of `LinkedSet.String`.
func (self *SyncLinkedSet) String() string {
	if self == nil {
		return `[]`
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.set.String()
}

// Concurrency-safe version of `LinkedSet.GoString`.
func (self *SyncLinkedSet) GoString() string {
	if self == nil {
		return `(*SyncLinkedSet)(nil)`
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	var buf strings.Builder
	buf.WriteString(`NewSyncLinkedSet`)
	self.set.writeGoString(&buf)
	return buf.String()
}
