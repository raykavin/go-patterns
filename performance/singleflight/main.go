// singleflight_duplicate_suppression.go
//
// This example demonstrates suppressing duplicate work with a singleflight-like
// helper (implemented locally to keep the example dependency-free).
//
// Key ideas illustrated:
//
//   - Concurrent callers for same key share one in-flight execution
//   - Only one expensive computation runs
//
package main

import (
	"fmt"
	"sync"
	"time"
)

type call struct {
	wg     sync.WaitGroup
	val    any
	err    error
	shared bool
}

type Singleflight struct {
	mu sync.Mutex
	m  map[string]*call
}

func NewSingleflight() *Singleflight {
	return &Singleflight{m: map[string]*call{}}
}

func (g *Singleflight) Do(key string, fn func() (any, error)) (any, error, bool) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		c.shared = true
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}
	c := &call{}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err, c.shared
}

func main() {
	sf := NewSingleflight()
	var wg sync.WaitGroup

	expensive := func() (any, error) {
		time.Sleep(200 * time.Millisecond)
		return "computed-value", nil
	}

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v, err, shared := sf.Do("key", expensive)
			fmt.Println("caller", id, "value=", v, "err=", err, "shared=", shared)
		}(i)
	}

	wg.Wait()
}
