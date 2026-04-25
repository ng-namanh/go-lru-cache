package cache

func NewCache(capacity int) *Cache {

	if capacity <= 0 {
		panic("invalid capacity: must be greater than 0")
	}

	head := &CacheNode{}
	tail := &CacheNode{}
	head.next = tail
	tail.prev = head

	return &Cache{
		items: make(map[string]*CacheNode),
		list: DoublyLinkedList{
			head: head,
			tail: tail,
		},
		metadata: CacheMetadata{
			capacity: capacity,
		},
	}
}

func (c *Cache) Get(key string) any {
	return nil
}

func (c *Cache) Put(key string, value any) {

}

func (c *Cache) Len() int {
	return len(c.items)
}

func (c *Cache) remove(node *CacheNode) {

	node.prev.next = node.next
	node.next.prev = node.prev

	node.prev = nil
	node.next = nil
}

func (c *Cache) Clear() {
	c.items = make(map[string]*CacheNode)
	c.list = DoublyLinkedList{
		head: &CacheNode{key: "", value: nil},
		tail: &CacheNode{key: "", value: nil},
	}
}
