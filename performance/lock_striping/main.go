// lock_striping.go
//
// This example demonstrates lock striping: split a big lock into many smaller locks
// based on key hashing to increase concurrency.
//
// Key ideas illustrated:
//
//   - Multiple mutexes protect disjoint key ranges
//   - Reduce contention in hot maps
//
package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type StripedMap struct {
	stripes []struct {
		mu sync.Mutex
		m  map[string]int
	}
}

func NewStripedMap(n int) *StripedMap {
	s := &StripedMap{stripes: make([]struct {
		mu sync.Mutex
		m  map[string]int
	}, n)}
	for i := range s.stripes {
		s.stripes[i].m = map[string]int{}
	}
	return s
}

func (s *StripedMap) idx(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32()) % len(s.stripes)
}

func (s *StripedMap) Inc(key string) {
	i := s.idx(key)
	st := &s.stripes[i]
	st.mu.Lock()
	st.m[key]++
	st.mu.Unlock()
}

func (s *StripedMap) Get(key string) int {
	i := s.idx(key)
	st := &s.stripes[i]
	st.mu.Lock()
	v := st.m[key]
	st.mu.Unlock()
	return v
}

func main() {
	sm := NewStripedMap(8)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("k-%d", i%5)
			sm.Inc(k)
		}(i)
	}
	wg.Wait()

	for i := 0; i < 5; i++ {
		k := fmt.Sprintf("k-%d", i)
		fmt.Println(k, sm.Get(k))
	}
}
