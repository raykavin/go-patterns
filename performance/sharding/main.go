// sharding.go
//
// This example demonstrates a basic sharding strategy: distribute keys across shards.
//
// Key ideas illustrated:
//
//   - Choose shard by hash(key) % numShards
//   - Each shard has its own lock and map
//
package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type shard struct {
	mu sync.RWMutex
	m  map[string]string
}

type ShardedKV struct {
	shards []shard
}

func NewShardedKV(n int) *ShardedKV {
	s := &ShardedKV{shards: make([]shard, n)}
	for i := range s.shards {
		s.shards[i].m = map[string]string{}
	}
	return s
}

func (s *ShardedKV) idx(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32()) % len(s.shards)
}

func (s *ShardedKV) Set(key, val string) {
	sh := &s.shards[s.idx(key)]
	sh.mu.Lock()
	sh.m[key] = val
	sh.mu.Unlock()
}

func (s *ShardedKV) Get(key string) (string, bool) {
	sh := &s.shards[s.idx(key)]
	sh.mu.RLock()
	v, ok := sh.m[key]
	sh.mu.RUnlock()
	return v, ok
}

func main() {
	kv := NewShardedKV(4)
	kv.Set("user:1", "Ada")
	kv.Set("user:2", "Grace")

	v, ok := kv.Get("user:2")
	fmt.Println("ok?", ok, "val:", v)
}
