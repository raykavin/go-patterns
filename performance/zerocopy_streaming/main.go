// zero_copy_streaming.go
//
// This example demonstrates "zero-copy" style streaming by avoiding extra allocations.
// In Go, io.Copy uses an internal buffer and can leverage optimized paths.
//
// Key ideas illustrated:
//
//   - Stream data from a Reader to a Writer with io.Copy
//   - Avoid building large intermediate byte slices
//
package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	src := strings.NewReader(strings.Repeat("go!", 1000))

	var dst bytes.Buffer
	n, err := io.Copy(&dst, src)
	fmt.Println("bytes copied:", n, "err:", err)
	fmt.Println("prefix:", dst.String()[:12])
}
