package lfu

import (
	"errors"
	"iter"
	"lfucache/internal/linkedlist"
)

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

// Cache
// O(capacity) memory
type Cache[K comparable, V any] interface {
	// Get returns the value of the key if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	Get(key K) (V, error)

	// Put updates the value of the key if present, or inserts the key if not already present.
	//
	// When the cache reaches its capacity, it should invalidate and remove the least frequently used key
	// before inserting a new item. For this problem, when there is a tie
	// (i.e., two or more keys with the same frequency), the least recently used key would be invalidated.
	//
	// O(1)
	Put(key K, value V)

	// All returns the iterator in descending order of frequency.
	// If two or more keys have the same frequency, the most recently used key will be listed first.
	//
	// O(capacity)
	All() iter.Seq2[K, V]

	// Size returns the cache size.
	//
	// O(1)
	Size() int

	// Capacity returns the cache capacity.
	//
	// O(1)
	Capacity() int

	// GetKeyFrequency returns the element's frequency if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	GetKeyFrequency(key K) (int, error)
}

// struct for containing input element for frequencies
type cacheElement[K comparable, V any] struct {
	key   K
	value V
	freq  int
}

// CacheImpl realization of Cache interface
type CacheImpl[K comparable, V any] struct {
	capacity    int
	size        int
	items       map[K]*linkedlist.Node[*cacheElement[K, V]]
	frequencies map[int]*linkedlist.Node[*linkedlist.LinkedList[*cacheElement[K, V]]]
	cache       *linkedlist.LinkedList[*linkedlist.LinkedList[*cacheElement[K, V]]]
}

func New[K comparable, V any](capacity ...int) *CacheImpl[K, V] {
	lfuCapacity := DefaultCapacity
	// determination of capacity
	if len(capacity) > 0 {
		if capacity[0] < 0 {
			panic("lfu cache capacity is negative")
		}
		lfuCapacity = capacity[0]
	}

	// constructing new element
	return &CacheImpl[K, V]{
		capacity:    lfuCapacity,
		items:       make(map[K]*linkedlist.Node[*cacheElement[K, V]], lfuCapacity),
		frequencies: make(map[int]*linkedlist.Node[*linkedlist.LinkedList[*cacheElement[K, V]]]),
		cache:       linkedlist.New[*linkedlist.LinkedList[*cacheElement[K, V]]](),
	}
}

func (l *CacheImpl[K, V]) Get(key K) (V, error) {
	var zero V
	elementNode, ok := l.items[key]

	// new element
	if !ok {
		return zero, ErrKeyNotFound
	}

	l.incrementFrequency(elementNode)
	return elementNode.Value.value, nil
}

func (l *CacheImpl[K, V]) Put(key K, value V) {
	// unable to put
	if l.capacity == 0 {
		return
	}

	// put logic
	if elementNode, ok := l.items[key]; ok {
		// element exists
		elementNode.Value.value = value
		l.incrementFrequency(elementNode)
	} else {
		// element doesn't exist
		// delete lfu element
		if l.size == l.capacity {
			l.evictLFU()
		}

		// add new one
		l.addNewEntry(key, value)
	}
}

// adds new entry to a list
// creates new frequency list if needed
func (l *CacheImpl[K, V]) addNewEntry(key K, value V) {
	newElement := &cacheElement[K, V]{key: key, value: value, freq: 1}
	newElementNode := &linkedlist.Node[*cacheElement[K, V]]{Value: newElement}

	// determination of list to push
	var frequencyList *linkedlist.LinkedList[*cacheElement[K, V]]
	if frequencyNode, ok := l.frequencies[1]; ok {
		frequencyList = frequencyNode.Value
	} else {
		// create new list if frequency 1 doesn't exist
		frequencyList = linkedlist.New[*cacheElement[K, V]]()
		frequencyNode = &linkedlist.Node[*linkedlist.LinkedList[*cacheElement[K, V]]]{Value: frequencyList}
		l.cache.PushFront(frequencyNode)
		l.frequencies[1] = frequencyNode
	}

	// push new element
	frequencyList.PushBack(newElementNode)
	l.items[key] = newElementNode
	l.size++
}

// IncrementFrequency increments frequency of a node
// replaces it to another frequency list
func (l *CacheImpl[K, V]) incrementFrequency(elementNode *linkedlist.Node[*cacheElement[K, V]]) {
	frequency := elementNode.Value.freq
	frequencyNode := l.frequencies[frequency]
	frequencyList := frequencyNode.Value

	// delete node from current frequency
	frequencyList.Pop(elementNode)

	elementNode.Value.freq++
	newFrequency := elementNode.Value.freq

	// determination of list to push new
	var newFrequencyList *linkedlist.LinkedList[*cacheElement[K, V]]
	if newFrequencyNode, ok := l.frequencies[newFrequency]; ok {
		// frequency + 1 exists
		newFrequencyList = newFrequencyNode.Value
	} else {
		// frequency + 1 doesn't exist,
		// create new list
		newFrequencyList = linkedlist.New[*cacheElement[K, V]]()
		newFrequencyNode = &linkedlist.Node[*linkedlist.LinkedList[*cacheElement[K, V]]]{Value: newFrequencyList}
		l.cache.PushAfter(frequencyNode, newFrequencyNode)
		l.frequencies[newFrequency] = newFrequencyNode
	}

	// delete empty frequency list
	if frequencyList.Size() == 0 {
		l.cache.Pop(frequencyNode)
		delete(l.frequencies, frequency)
	}

	newFrequencyList.PushBack(elementNode)
}

// delete lfu element
func (l *CacheImpl[K, V]) evictLFU() {
	frequencyNode := l.cache.Head()

	// no frequencies (error)
	if frequencyNode == nil {
		return
	}
	frequencyList := frequencyNode.Value

	// the lowest value of lru element is head
	elementNode := frequencyList.Head()

	// no elements in this list (error)
	if elementNode == nil {
		return
	}

	delete(l.items, elementNode.Value.key)

	frequencyList.Pop(elementNode)

	// delete empty frequency node
	if frequencyList.Size() == 0 {
		l.cache.Pop(frequencyNode)
		delete(l.frequencies, elementNode.Value.freq)
	}

	l.size--
}

func (l *CacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// starts from tail because it is the biggest frequency
		for frequencyNode := l.cache.Tail(); frequencyNode != nil; frequencyNode = frequencyNode.Prev {
			frequencyList := frequencyNode.Value
			// starts from tail because it is the lru element at current frequency
			for elementNode := frequencyList.Tail(); elementNode != nil; elementNode = elementNode.Prev {
				if !yield(elementNode.Value.key, elementNode.Value.value) {
					return
				}
			}
		}
	}
}

func (l *CacheImpl[K, V]) Size() int {
	return l.size
}

func (l *CacheImpl[K, V]) Capacity() int {
	return l.capacity
}

func (l *CacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	elementNode, ok := l.items[key]
	// if key doesn't exist return error
	if !ok {
		return 0, ErrKeyNotFound
	}

	// if exists return frequency
	return elementNode.Value.freq, nil
}
