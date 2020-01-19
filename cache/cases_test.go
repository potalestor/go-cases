package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCases()
	if c.m == nil {
		t.Fatal("c.m expected not nil")
	}
}

func TestPopulate01(t *testing.T) {
	t0 := time.Now().Add(-time.Hour)
	t1 := t0.Add(time.Minute)
	t2 := t1.Add(time.Minute)
	c := NewCases()
	count := c.Populate("1", &t0, time.Hour, nil)
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
	count = c.Populate("1", &t1, time.Hour, nil)
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
	count = c.Populate("1", &t2, time.Hour, nil)
	if count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
	fmt.Println(c.String())
}

func TestPopulate02(t *testing.T) {
	t0 := time.Now().Add(-time.Hour)
	fmt.Println(t0.UnixNano())
	t1 := t0.Add(10 * time.Minute)
	fmt.Println(t1.UnixNano())
	t2 := t1.Add(time.Minute)
	fmt.Println(t2.UnixNano())
	c := NewCases()
	count := c.Populate("1", &t0, 2*time.Minute, nil)
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
	count = c.Populate("1", &t1, 2*time.Minute, nil)
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
	count = c.Populate("1", &t2, 2*time.Minute, nil)
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
	fmt.Println(c.String())
}

func TestPopulate03(t *testing.T) {
	t0 := time.Now().Add(-time.Hour)
	fmt.Println(t0.UnixNano())
	t1 := t0.Add(time.Minute)
	fmt.Println(t1.UnixNano())
	t2 := t1.Add(10 * time.Minute)
	fmt.Println(t2.UnixNano())
	c := NewCases()
	count := c.Populate("1", &t0, 2*time.Minute, nil)
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
	count = c.Populate("1", &t1, 2*time.Minute, nil)
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
	count = c.Populate("1", &t2, 2*time.Minute, nil)
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
	fmt.Println(c.String())
}

func BenchmarkPopulate(b *testing.B) {
	c := NewCases()
	for i := 0; i < b.N; i++ {
		t0 := time.Now().Add(-time.Hour)
		key := fmt.Sprintf("%d", i)
		count := c.Populate(key, &t0, 2*time.Minute, nil)
		if count >= 3 {
			c.Evict(key)
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	c := map[string]*item{}
	for i := 0; i < b.N; i++ {
		t0 := time.Now().Add(-time.Hour)
		c[fmt.Sprintf("%d", i)] = aquireItem().set(&t0, nil)
	}
}
