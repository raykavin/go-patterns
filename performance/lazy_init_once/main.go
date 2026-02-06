// lazy_initialization_sync_once.go
//
// This example demonstrates lazy initialization using sync.Once.
//
// Key ideas illustrated:
//
//   - Initialization logic runs exactly once
//   - Safe for concurrent access
//
package main

import (
	"fmt"
	"sync"
)

type Config struct {
	DSN string
}

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{DSN: "postgres://user:pass@localhost/db"}
	})
	return config
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(GetConfig().DSN)
		}()
	}
	wg.Wait()
}
