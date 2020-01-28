package pubsub

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestChanger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	tc := NewTrackChanger(ctx, &wg)
	var subs []Subscriber
	for i := 0; i < 5; i++ {
		sub := NewSub(strconv.Itoa(i))
		tc.Subscribe(sub.Key(), sub.EventChan())
		subs = append(subs, sub)
	}
	time.Sleep(15 * time.Second)
	cancel()
	wg.Wait()
	time.Sleep(time.Second)

}
