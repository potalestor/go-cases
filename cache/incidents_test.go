package cache

import (
	"testing"
	"time"
)

func TestClearIncident(t *testing.T) {
	c := NewIncidents()
	t0 := time.Now().Add(-time.Hour)
	t1 := time.Now()

	c.set("1", &t0, nil)
	if _, ok := c.TryGetValue("1", &t0, time.Minute); !ok {
		t.Fatal("expected exist key")
	}
	c.ClearIncident("1")
	if _, ok := c.TryGetValue("1", &t0, time.Minute); ok {
		t.Fatal("expected not exist key")
	}

	c.set("1", &t1, nil)
	if _, ok := c.TryGetValue("1", &t1, time.Minute); !ok {
		t.Fatal("expected exist key")
	}
}
