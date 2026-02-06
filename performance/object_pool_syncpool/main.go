// object_pool_syncpool.go
//
// This example demonstrates object pooling with sync.Pool.
//
// Key ideas illustrated:
//
//   - Reuse temporary objects to reduce allocations
//   - Keep pooled objects in a valid reset state
//
package main

import (
	"bytes"
	"fmt"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func main() {
	for i := 1; i <= 5; i++ {
		b := bufPool.Get().(*bytes.Buffer)
		b.Reset()
		fmt.Fprintf(b, "hello %d", i)
		fmt.Println(b.String())
		bufPool.Put(b)
	}
}
