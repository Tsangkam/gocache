package lru

import "container/list"

// core cache obj, use map and Double-linked list, map for index, list for sort
type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List // new->old
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

// create new cache unit
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted, // remove callback
	}
}

// get value by key
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) // move the element to front after it is searched
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele) // remove from the double-linked list
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // remove from the map
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) // update the size of cache
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value) // callback
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // element is already exist
		c.ll.MoveToFront(ele) // move the element to front
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value}) // move the new element to front
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest() // remove the oldest element while oversize
	}
}
