package cache

import (
	"sync"
)

type CacheNode struct {
	key   string     // Primary identifier (e.g., "asset_12345")
	value any        // Cached data (serialized metadata, content refs)
	prev  *CacheNode // Pointer to previous node in the list
	next  *CacheNode // Pointer to next node in the list
}

// Essential for implementing the LRU cache eviction policy
type DoublyLinkedList struct {
	head *CacheNode // Sentinel node (most recently used end)
	tail *CacheNode // Sentinel node (least recently used end)
}

type CacheMetadata struct {
	capacity    int   // Configuration parameter
	totalHits   int64 // Successful get() operations
	totalMisses int64 // Failed get() operations
}

type Cache struct {
	mu       sync.Mutex
	metadata CacheMetadata
	items    map[string]*CacheNode
	list     DoublyLinkedList
}
