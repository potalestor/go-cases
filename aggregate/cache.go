package aggregate

import (
	"time"
)

const (
	MaxCaseSize = 100
)

var (
	IntervalLength int64 = int64(10 * time.Minute)
)

type DataHandler func(data []interface{}) error

type Case struct {
	Created int64
	Data    []interface{}
}

type Aggregate struct {
	DataHandler DataHandler
	Map         map[string]*Case
}

func NewAggregate() *Aggregate {
	return &Aggregate{Map: map[string]*Case{}}
}

func (c *Aggregate) Populate(key string, value interface{}) int {
	now := time.Now().UnixNano()
	cs, ok := c.Map[key]
	if ok {
		if cs.Created-now <= 0 {
			c.DataHandler(cs.Data)
			cs.Data = cs.Data[:0]
			cs.Created = now + IntervalLength
		}
		cs.Data = append(cs.Data, value)
		c.Map[key] = cs
		return len(cs.Data)
	}
	cs = &Case{Data: make([]interface{}, 1, 100)}
	cs.Created = now + IntervalLength
	cs.Data[0] = value
	return 1
}

func (c *Aggregate) Evict() {
	now := time.Now().UnixNano()
	for key, cs := range c.Map {
		if cs.Created-now <= 0 {
			c.DataHandler(cs.Data)
			cs.Data = cs.Data[:0]
			cs.Created = 0
			c.Map[key] = cs
		}
	}
}
