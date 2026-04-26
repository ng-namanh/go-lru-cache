package cache

import (
	"sync"
	"testing"
)

func TestNewCache_InvalidCapacity(t *testing.T) {
	t.Parallel()
	c, err := NewCache(0)
	if err == nil {
		t.Fatal("expected error for capacity 0")
	}
	if c != nil {
		t.Fatal("expected nil cache on error")
	}
	_, err = NewCache(-1)
	if err == nil {
		t.Fatal("expected error for negative capacity")
	}
}

func TestGet_Miss(t *testing.T) {
	t.Parallel()
	c, err := NewCache(2)
	if err != nil {
		t.Fatal(err)
	}
	v, ok := c.Get("missing")
	if ok || v != nil {
		t.Fatalf("Get miss: ok=%v value=%v", ok, v)
	}
}

func TestPutGet_Hit(t *testing.T) {
	t.Parallel()
	c, err := NewCache(2)
	if err != nil {
		t.Fatal(err)
	}
	c.Put("a", 1)
	v, ok := c.Get("a")
	if !ok {
		t.Fatal("expected hit")
	}
	if v != 1 {
		t.Fatalf("value: got %v want 1", v)
	}
	if c.Len() != 1 {
		t.Fatalf("Len: got %d want 1", c.Len())
	}
}

func TestEviction_Classic(t *testing.T) {
	t.Parallel()
	// Put(1) Put(2) Get(1) Put(3) -> 2 evicted, 1 and 3 remain
	c, err := NewCache(2)
	if err != nil {
		t.Fatal(err)
	}
	c.Put("1", 1)
	c.Put("2", 2)
	if c.Len() != 2 {
		t.Fatalf("Len after two puts: got %d", c.Len())
	}
	if _, ok := c.Get("1"); !ok {
		t.Fatal("Get 1: expected hit before eviction sequence")
	}
	c.Put("3", 3)

	_, ok2 := c.Get("2")
	if ok2 {
		t.Fatal("key 2 should be evicted")
	}
	if v, ok := c.Get("1"); !ok || v != 1 {
		t.Fatalf("1: ok=%v v=%v", ok, v)
	}
	if v, ok := c.Get("3"); !ok || v != 3 {
		t.Fatalf("3: ok=%v v=%v", ok, v)
	}
	if c.Len() != 2 {
		t.Fatalf("Len: got %d want 2", c.Len())
	}
}

func TestPut_Existing_DoesNotEvictOthers(t *testing.T) {
	t.Parallel()
	c, err := NewCache(2)
	if err != nil {
		t.Fatal(err)
	}
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("a", 10) // update only
	if c.Len() != 2 {
		t.Fatalf("Len: got %d want 2", c.Len())
	}
	if v, _ := c.Get("b"); v != 2 {
		t.Fatalf("b should still be present, got %v", v)
	}
	if v, _ := c.Get("a"); v != 10 {
		t.Fatalf("a: got %v want 10", v)
	}
}

func TestCapacity1(t *testing.T) {
	t.Parallel()
	c, err := NewCache(1)
	if err != nil {
		t.Fatal(err)
	}
	c.Put("a", 1)
	c.Put("b", 2)
	if _, ok := c.Get("a"); ok {
		t.Fatal("a should be evicted")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("b: ok=%v v=%v", ok, v)
	}
}

func TestClear(t *testing.T) {
	t.Parallel()
	c, err := NewCache(2)
	if err != nil {
		t.Fatal(err)
	}
	c.Put("a", 1)
	c.Clear()
	if c.Len() != 0 {
		t.Fatalf("Len after Clear: got %d", c.Len())
	}
	if _, ok := c.Get("a"); ok {
		t.Fatal("Get after Clear: expected miss")
	}
	// Re-use after clear
	c.Put("x", 9)
	if v, ok := c.Get("x"); !ok || v != 9 {
		t.Fatalf("after re-Put: ok=%v v=%v", ok, v)
	}
}

func TestGet_PromotesToMRU(t *testing.T) {
	t.Parallel()
	c, err := NewCache(2)
	if err != nil {
		t.Fatal(err)
	}
	c.Put("a", 1)
	c.Put("b", 2)
	// a is LRU; touch a with Get, then add c; b (current LRU) should be evicted
	if _, ok := c.Get("a"); !ok {
		t.Fatal("Get a")
	}
	c.Put("c", 3)
	if _, ok := c.Get("b"); ok {
		t.Fatal("b should be evicted after a was promoted")
	}
	if v, _ := c.Get("a"); v != 1 {
		t.Fatalf("a: %v", v)
	}
	if v, _ := c.Get("c"); v != 3 {
		t.Fatalf("c: %v", v)
	}
}

func TestConcurrentGetPut(t *testing.T) {
	c, err := NewCache(64)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				k := "k"
				c.Put(k, i+base)
				_, _ = c.Get(k)
			}
		}(g)
	}
	wg.Wait()
}
