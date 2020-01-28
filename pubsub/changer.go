package pubsub

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TrackChanger struct {
	sync.RWMutex
	subs map[string]chan<- string
	ctx  context.Context
	wg   *sync.WaitGroup
}

func NewTrackChanger(ctx context.Context, wg *sync.WaitGroup) *TrackChanger {
	tc := &TrackChanger{
		subs: map[string]chan<- string{},
		ctx:  ctx,
		wg:   wg,
	}
	go tc.watch()
	return tc
}

func (tc *TrackChanger) Subscribe(key string, eventChan chan<- string) {
	tc.Lock()
	defer tc.Unlock()
	if _, ok := tc.subs[key]; !ok {
		fmt.Println(key, "added")
		tc.subs[key] = eventChan
	}
}

func (tc *TrackChanger) Unsubscribe(key string) {
	tc.Lock()
	defer tc.Unlock()
	if sub, ok := tc.subs[key]; ok {
		close(sub)
		delete(tc.subs, key)
	}
}

func (tc *TrackChanger) watch() {
	defer tc.wg.Done()

	for {
		tc.checkAndNotifyChanges()
		select {
		case <-tc.ctx.Done():
			tc.closeAllSubscribers()
			fmt.Println("watcher closed")
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func (tc *TrackChanger) closeAllSubscribers() {
	tc.Lock()
	defer tc.Unlock()
	for key, sub := range tc.subs {
		close(sub)
		delete(tc.subs, key)
	}

}

func (tc *TrackChanger) checkAndNotifyChanges() {
	tc.RLock()
	defer tc.RUnlock()
	for key, sub := range tc.subs {
		sub <- key
	}
}
