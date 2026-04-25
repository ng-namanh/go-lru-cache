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

}

func (c *Cache) Len() int {
	return len(c.items)
}

// Clear the cache
func (c *Cache) Clear() {
	c.items = make(map[string]*CacheNode)
	c.list = DoublyLinkedList{
		head: &CacheNode{key: "", value: nil},
		tail: &CacheNode{key: "", value: nil},
	}
}

func (c *Cache) remove(node *CacheNode) {
	node.prev.next = node.next
	node.next.prev = node.prev

	node.prev = nil
	node.next = nil
}

// Move the node to the front of the list (most recently used)
func (c *Cache) moveToFront(node *CacheNode) {
	if c.list.head.next == node {
		return
	}

	c.remove(node)

	node.prev = c.list.head
	node.next = c.list.head.next
	c.list.head.next.prev = node
	c.list.head.next = node
}
