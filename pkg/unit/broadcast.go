package unit

import (
	"context"
)

type Value interface{}

type Subscription struct {
	c chan Value
}

func (s Subscription) Input() Input {
	return s.c
}

func NewSubscription() Subscription {
	return Subscription{
		c: make(chan Value, 1),
	}
}

type Subscriptions map[string]Subscription

func (s Subscriptions) Inputs() map[string]Input {
	result := make(map[string]Input)
	for k, v := range s {
		result[k] = v.Input()
	}
	return result
}

type Broadcasts map[string]Broadcast

func (s Broadcasts) Outputs() map[string]Output {
	result := make(map[string]Output)
	for k, v := range s {
		result[k] = v.Output()
	}
	return result
}

type Broadcast interface {
	Output() Output
	Start()
	Stop()
	Subscribe(Subscription)
	Unsubscribe(Subscription)
}

type BroadcastAgent struct {
	values      chan Value
	subscribers map[chan<- Value]struct{}
	cancel      context.CancelFunc

	caching bool
	cache   *Value
}

func NewBroadcast(ctx context.Context, schema OutputSchema) *BroadcastAgent {
	b := &BroadcastAgent{
		values:      make(chan Value, 1),
		subscribers: map[chan<- Value]struct{}{},
		caching:     !schema.NoCaching,
	}

	cCtx, cancel := context.WithCancel(ctx)
	b.cancel = cancel
	go func() {
		for {
			select {
			case <-cCtx.Done():
				break
			case v := <-b.values:
				for s := range b.subscribers {
					s <- v
				}
				val := (Value)(v)
				b.cache = &val
			}
		}
	}()

	return b
}

func (f *BroadcastAgent) Output() Output {
	return f.values
}
func (f *BroadcastAgent) Start() {}
func (f *BroadcastAgent) Stop() {
	f.cancel()
}
func (f *BroadcastAgent) Subscribe(sub Subscription) {
	f.subscribers[sub.c] = struct{}{}

	if f.caching && f.cache != nil {
		sub.c <- *f.cache
	}
}
func (f *BroadcastAgent) Unsubscribe(sub Subscription) {
	delete(f.subscribers, sub.c)
}
