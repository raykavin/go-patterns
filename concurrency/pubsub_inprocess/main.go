// publish_subscribe_inprocess.go
//
// This example demonstrates an in-process publish/subscribe pattern using channels.
//
// Key ideas illustrated:
//
//   - A broker goroutine manages subscribers
//   - Subscribers receive published messages
//   - Unsubscribe on context cancellation
//
package main

import (
	"context"
	"fmt"
	"time"
)

type Broker struct {
	subscribe   chan chan string
	unsubscribe chan chan string
	publish     chan string
}

func NewBroker() *Broker {
	b := &Broker{
		subscribe:   make(chan chan string),
		unsubscribe: make(chan chan string),
		publish:     make(chan string),
	}
	go b.run()
	return b
}

func (b *Broker) run() {
	subs := map[chan string]struct{}{}
	for {
		select {
		case ch := <-b.subscribe:
			subs[ch] = struct{}{}
		case ch := <-b.unsubscribe:
			delete(subs, ch)
			close(ch)
		case msg := <-b.publish:
			for ch := range subs {
				// best-effort send: drop if subscriber is slow
				select {
				case ch <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Subscribe(ctx context.Context, buf int) <-chan string {
	ch := make(chan string, buf)
	b.subscribe <- ch

	go func() {
		<-ctx.Done()
		b.unsubscribe <- ch
	}()

	return ch
}

func (b *Broker) Publish(msg string) {
	b.publish <- msg
}

func main() {
	b := NewBroker()

	ctx, cancel := context.WithCancel(context.Background())
	sub := b.Subscribe(ctx, 4)

	go func() {
		for msg := range sub {
			fmt.Println("sub got:", msg)
		}
	}()

	for i := 1; i <= 5; i++ {
		b.Publish(fmt.Sprintf("event-%d", i))
		time.Sleep(60 * time.Millisecond)
	}
	cancel()
	time.Sleep(100 * time.Millisecond)
}
