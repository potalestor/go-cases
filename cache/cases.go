package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type item struct {
	Created int64
	Data    interface{}
}

func (i *item) set(created *time.Time, data interface{}) *item {
	i.Created = created.UnixNano()
	i.Data = data
	return i
}

func (i *item) valid(created *time.Time, life time.Duration) bool {
	return i.Created-created.Add(-life).UnixNano() > 0
}

type items []*item

func newItems(created *time.Time, data interface{}) items {
	i := make([]*item, 1, 5)
	i[0] = aquireItem().set(created, data)
	return i
}

func (i *items) populate(created *time.Time,
	life time.Duration,
	data interface{},
) int {
	*i = append(*i, aquireItem().set(created, data))
	if !(*i)[0].valid(created, life) {
		releaseItem((*i)[0])
		i.removeExpired(created, life)
	}
	return len(*i)
}

func (i *items) removeExpired(created *time.Time, life time.Duration) {
	skipped := 1
	for j := 1; j < len(*i)-1; j++ {
		if !(*i)[j].valid(created, life) {
			releaseItem((*i)[j])
			skipped++
		}
	}
	*i = (*i)[skipped:]
}

var itemPool = sync.Pool{
	New: func() interface{} { return &item{} },
}

func aquireItem() *item {
	return itemPool.Get().(*item)
}
func releaseItem(i *item) {
	if i != nil {
		itemPool.Put(i)
	}
}

type Cases struct {
	m map[string]items
}

func NewCases() *Cases {
	return &Cases{m: make(map[string]items, 1000)}
}

func (c *Cases) Populate(
	key string,
	created *time.Time,
	life time.Duration,
	data interface{},
) int {
	i, ok := c.m[key]
	if !ok || (len(i) == 0) {
		c.m[key] = newItems(created, data)
		return 1
	}
	count := i.populate(created, life, data)
	c.m[key] = i
	return count
}

func (c *Cases) Evict(key string) []interface{} {
	is, ok := c.m[key]
	if ok && (len(is) > 0) {
		data := make([]interface{}, 0, len(is))
		for _, i := range is {
			data = append(data, i.Data)
			releaseItem(i)
		}
		c.m[key] = is[:0]
		return data
	}
	return nil
}

func (c *Cases) ClearCase(key string) {
	is, ok := c.m[key]
	if ok && (len(is) > 0) {
		for _, i := range is {
			releaseItem(i)
		}
		c.m[key] = is[:0]
	}
}

func (c *Cases) String() string {
	var b strings.Builder
	for k, v := range c.m {
		b.WriteString(k)
		b.WriteString(": ")
		for i := range v {
			b.WriteString(fmt.Sprintf("%d ", v[i].Created))
		}
	}
	return b.String()
}
