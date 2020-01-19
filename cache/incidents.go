package cache

import (
	"time"
)

type Incidents struct {
	m map[string]*item
}

func NewIncidents() *Incidents {
	return &Incidents{m: make(map[string]*item)}
}

func (i *Incidents) TryGetValue(key string, created *time.Time, life time.Duration) (interface{}, bool) {
	v, ok := i.m[key]
	if ok && v.valid(created, life) {
		return v.Data, true
	}
	return nil, false
}

func (i *Incidents) ClearIncident(key string) {
	v, ok := i.m[key]
	if ok {
		v.Created = 0
	}
}

func (i *Incidents) set(key string, created *time.Time, data interface{}) {
	v, ok := i.m[key]
	if !ok {
		i.m[key] = aquireItem().set(created, data)
		return
	}
	v.set(created, data)
}
