// pipeline_staged_processing.go
//
// This example demonstrates a staged pipeline where each stage runs in its own
// goroutine and communicates via channels.
//
// Key ideas illustrated:
//
//   - Separating processing into stages
//   - Backpressure naturally occurs through bounded channels
//   - Closing output channels when input is drained
//
// Common uses:
//   - ETL-style transformations
//   - Streaming processing
//   - IO -> parse -> validate -> persist
//
package main

import (
	"fmt"
	"strings"
)

func gen(lines ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, s := range lines {
			out <- s
		}
	}()
	return out
}

func upper(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			out <- strings.ToUpper(s)
		}
	}()
	return out
}

func filterPrefix(in <-chan string, prefix string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			if strings.HasPrefix(s, prefix) {
				out <- s
			}
		}
	}()
	return out
}

func main() {
	src := gen("go", "gopher", "java", "golang", "rust")
	stage1 := upper(src)
	stage2 := filterPrefix(stage1, "GO")

	for v := range stage2 {
		fmt.Println(v)
	}
}
