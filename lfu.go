package lfu

import (
	"container/list"
	"fmt"
)

type Cache struct {
	values map[string]*cacheEntry
	freqs *list.List
	len int
}

type cacheEntry struct {
	key string
	value interface{}
	freqNode *list.Element
}

type listEntry struct {
	entries map[*cacheEntry]byte
	freq int
}

func New()*Cache {
	c := new(Cache)
	c.values = make(map[string]*cacheEntry)
	c.freqs = list.New()
	return c
}

func (c *Cache) Get(key string)interface{} {
	if e, ok := c.values[key]; ok {
		c.increment(e)
		return e.value
	}
	return nil
}

func (c *Cache) Set(key string, value interface{}) {
	if e, ok := c.values[key]; ok {
		// value already exists for key.  overwrite
		e.value = value
		c.increment(e)
	} else {
		// value doesn't exist.  insert
		e := new(cacheEntry)
		e.key = key
		e.value = value
		c.values[key] = e
		c.increment(e)
		c.len++
	}
}

func (c *Cache) Len()int {
	return c.len
}

func (c *Cache) Evict(count int)int{
	var evicted int
	for i := 0; i < count; {
		if place := c.freqs.Front(); place != nil {
			for entry, _ := range place.Value.(*listEntry).entries {
				if i < count {
					delete(c.values, entry.key)
					c.remEntry(place, entry)
					evicted++
					c.len--
					i++
				}
			}
		}
	}
	return evicted
}

func (c *Cache) increment(e *cacheEntry) {
	currentPlace := e.freqNode
	var nextFreq int
	var nextPlace *list.Element
	if currentPlace == nil {
		// new entry
		nextFreq = 1
		nextPlace = c.freqs.Front()
	} else {
		// move up
		nextFreq = currentPlace.Value.(*listEntry).freq + 1
		nextPlace = currentPlace.Next()
	}

	if nextPlace != nil {
		fmt.Printf("%v Looking for: %v with freq %v\n", e.key, nextPlace.Value, nextFreq)
	} else {
		fmt.Printf("%v Looking for: %v with freq %v\n", e.key, nextPlace, nextFreq)
	}
	
	if nextPlace == nil || nextPlace.Value.(*listEntry).freq != nextFreq {
		// create a new list entry
		li := new(listEntry)
		li.freq = nextFreq
		li.entries = make(map[*cacheEntry]byte)
		if currentPlace != nil {
			nextPlace = c.freqs.InsertAfter(li, currentPlace)
		} else {
			nextPlace = c.freqs.PushFront(li)
		}
	}
	e.freqNode = nextPlace
	nextPlace.Value.(*listEntry).entries[e] = 1
	fmt.Println(nextPlace.Value)
	if(currentPlace != nil){
		// remove from current position
		c.remEntry(currentPlace, e)
	}
}

func (c *Cache) remEntry(place *list.Element, entry *cacheEntry){
	entries := place.Value.(*listEntry).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqs.Remove(place)
	}
}