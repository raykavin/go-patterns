// functional_options.go
//
// This example demonstrates the Functional Options pattern for configurable constructors.
//
// Key ideas illustrated:
//
//   - Option functions mutate config safely
//   - Avoids large constructors with many parameters
//
package main

import (
	"fmt"
	"time"
)

type Client struct {
	baseURL string
	timeout time.Duration
	retries int
}

type Option func(*Client)

func WithBaseURL(v string) Option { return func(c *Client) { c.baseURL = v } }
func WithTimeout(v time.Duration) Option { return func(c *Client) { c.timeout = v } }
func WithRetries(v int) Option { return func(c *Client) { c.retries = v } }

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL: "https://example.local",
		timeout: 2 * time.Second,
		retries: 2,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func main() {
	c := NewClient(
		WithBaseURL("https://api.service"),
		WithTimeout(500*time.Millisecond),
		WithRetries(5),
	)
	fmt.Printf("%+v\n", *c)
}
