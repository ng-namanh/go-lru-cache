package cache

import (
	"errors"
	"fmt"
)

func NewCache(capacity int) (*Cache, error) {

	if capacity <= 0 {
		return nil, errors.New("invalid capacity: must be greater than 0")
	}

	head := &CacheNode{}
	tail := &CacheNode{}
	head.next = tail
	tail.prev = head

	cache := &Cache{
		items: make(map[string]*CacheNode),
		list: DoublyLinkedList{
			head: head,
			tail: tail,
		},
		metadata: CacheMetadata{
			capacity: capacity,
		},
	}

	return cache, nil
}

func (c *Cache) Get(key string) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.items[key]
	if !ok {
		c.metadata.totalMisses++
		return nil, fmt.Errorf("key not found: %s", key)
	}

	c.metadata.totalHits++
	c.moveToFront(node)

	return node.value, nil
}

func (c *Cache) Put(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the key already exists
	if node, ok := c.items[key]; ok {
		node.value = value
		c.moveToFront(node)
		return
	}

	// Check capacity (only for new keys)
	if len(c.items) >= c.metadata.capacity {
		c.evictLRU()
	}

	// Insert new node as MRU
	node := &CacheNode{key: key, value: value}
	c.addToFront(node)
	c.items[key] = node
}

func (c *Cache) Len() int {
	return len(c.items)
}

// Clear the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list.head = &CacheNode{}
	c.list.tail = &CacheNode{}
	c.list.head.next = c.list.tail
	c.list.tail.prev = c.list.head

	c.items = make(map[string]*CacheNode)

}

func (c *Cache) remove(node *CacheNode) {
	// Defensive guards: do nothing for nil, sentinels, or unlinked nodes.
	if node == nil || node == c.list.head || node == c.list.tail || node.prev == nil || node.next == nil {
		return
	}
	node.prev.next = node.next
	node.next.prev = node.prev

	node.prev = nil
	node.next = nil
}

// Purpose: Move the node to the front of the list (most recently used)
//
// Used in:
//   - Get() for a hit: existing nodes are already in the list, so we just move them to the MRU.
func (c *Cache) moveToFront(node *CacheNode) {
	if c.list.head.next == node {
		return
	}

	c.remove(node)
	c.addToFront(node)
}

// Purpose: inserts a node as the MRU entry (right after the `head` sentinel)
//
// Used in:
//   - Put() for a new key: new nodes aren't in the list yet, so we must insert them at the MRU
//   - (Optionally) moveToFront() can be implemented as: remove(node) then addToFront(node) for existing nodes.
//
// What it does: pointer slicing only (no `remove`)
func (c *Cache) addToFront(node *CacheNode) {
	mru := c.list.head.next

	node.prev = c.list.head
	node.next = mru
	mru.prev = node
	c.list.head.next = node
}

// Purpose: evicts the least recently used node from the cache (LRU)
//
// Used in:
//   - Put() when the cache is full: we need to remove the LRU to make room for the new node.
//
// What it does: pointer slicing only (no `remove`) and deleting the node from the map
func (c *Cache) evictLRU() {
	lru := c.list.tail.prev

	// If the LRU is the head, return
	if lru == c.list.head {
		return
	}

	c.remove(lru)
	delete(c.items, lru.key)

}

func (c *Cache) checkKeyExists(key string) bool {
	_, ok := c.items[key]
	return ok
}
