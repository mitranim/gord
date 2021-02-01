## Overview

**Go** **ord**ered sets.

Features:

* `LinkedSet`: ordered set with near-constant-time (O(1)) performance for inserting, deleting, and moving elements. Backed by a map and a doubly-linked list.

* `SyncLinkedSet`: concurrency-safe `LinkedSet`, slightly slower.

* `SliceSet`: slice-backed ordered set. Simpler and faster for small sets, extreme performance degradation for large sets.

* All implementations share a common interface.

* Small with no dependencies.

See the documentation at https://godoc.org/github.com/mitranim/gord.

Example:

```go
import "github.com/mitranim/gord"

set := gord.NewOrdSet()

// Note the order.
set.Add(20)
set.Add(10)
set.Add(30)

// Redundant and doesn't change the order.
set.Add(30)
set.Add(10)
set.Add(20)

set.Has(10)    // true
set.Has(40)    // false
set.Values()   // []interface{}{20, 10, 30}

set.PopFirst() // 20
set.PopLast()  // 30
set.Values()   // []interface{}{10}
```

## Known Limitations

* Has room for performance optimizations.

* `LinkedSet` and `SyncLinkedSet` don't expose a way to iterate over elements without allocating a slice via `.Values()`. Can be rectified on demand.

## License

https://unlicense.org

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
