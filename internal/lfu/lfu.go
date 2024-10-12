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
type сacheImpl[K comparable, V any] struct {
	capacity    int
	size        int
	items       map[K]*linkedlist.Node[cacheElement[K, V]]
	frequencies map[int]*linkedlist.Node[linkedlist.LinkedList[cacheElement[K, V]]]
	cache       linkedlist.LinkedList[linkedlist.LinkedList[cacheElement[K, V]]]
}

func New[K comparable, V any](capacity ...int) *сacheImpl[K, V] {
	lfuCapacity := DefaultCapacity
	// determination of capacity
	if len(capacity) > 0 {
		if capacity[0] < 0 {
			panic("lfu cache capacity is negative")
		}
		lfuCapacity = capacity[0]
	}

	// constructing new element
	return &сacheImpl[K, V]{
		capacity:    lfuCapacity,
		items:       make(map[K]*linkedlist.Node[cacheElement[K, V]], lfuCapacity),
		frequencies: make(map[int]*linkedlist.Node[linkedlist.LinkedList[cacheElement[K, V]]]),
		cache:       linkedlist.New[linkedlist.LinkedList[cacheElement[K, V]]](),
	}
}

func (l *сacheImpl[K, V]) Get(key K) (V, error) {
	var zero V
	elementNode, ok := l.items[key]

	// new element
	if !ok {
		return zero, ErrKeyNotFound
	}

	l.incrementFrequency(elementNode)
	return elementNode.Value.value, nil
}

func (l *сacheImpl[K, V]) Put(key K, value V) {
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
func (l *сacheImpl[K, V]) addNewEntry(key K, value V) {
	newElement := cacheElement[K, V]{key: key, value: value, freq: 1}
	newElementNode := &linkedlist.Node[cacheElement[K, V]]{Value: newElement}

	// push to init frequency
	l.pushToFrequency(1, newElementNode)

	l.items[key] = newElementNode
	l.size++
}

// IncrementFrequency increments frequency of a node
// replaces it to another frequency list
func (l *сacheImpl[K, V]) incrementFrequency(elementNode *linkedlist.Node[cacheElement[K, V]]) {
	frequency := elementNode.Value.freq
	frequencyNode := l.frequencies[frequency]
	frequencyList := frequencyNode.Value

	// delete node from current frequency
	frequencyList.Pop(elementNode)

	// increment frequency
	elementNode.Value.freq++

	// push element
	l.pushToFrequency(elementNode.Value.freq, elementNode)

	// delete empty frequency list
	if frequencyList.Size() == 0 {
		l.cache.Pop(frequencyNode)
		delete(l.frequencies, frequency)
	}
}

func (l *сacheImpl[K, V]) pushToFrequency(frequency int, elementNode *linkedlist.Node[cacheElement[K, V]]) {
	var frequencyList linkedlist.LinkedList[cacheElement[K, V]]
	if frequencyNode, ok := l.frequencies[frequency]; ok {
		frequencyList = frequencyNode.Value
	} else {
		// create new list if frequency doesn't exist
		frequencyList = linkedlist.New[cacheElement[K, V]]()
		frequencyNode = &linkedlist.Node[linkedlist.LinkedList[cacheElement[K, V]]]{Value: frequencyList}

		// init frequency
		if frequency == 1 {
			l.cache.PushFront(frequencyNode)
		} else { // next frequency
			l.cache.PushAfter(l.frequencies[frequency-1], frequencyNode)
		}
		l.frequencies[frequency] = frequencyNode
	}

	frequencyList.PushBack(elementNode)
}

// delete lfu element
func (l *сacheImpl[K, V]) evictLFU() {
	frequencyNode := l.cache.Head()

	// no frequencies (error)
	if frequencyNode == nil {
		panic("Cache is empty, evict to early")
	}
	frequencyList := frequencyNode.Value

	// the lowest value of lru element is head
	elementNode := frequencyList.Head()

	// no elements in this list (error)
	if elementNode == nil {
		panic("Frequency list is empty")
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

func (l *сacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// starts from tail because it is the biggest frequency
		for frequencyList := range l.cache.FromTail() {
			// starts from tail because it is the lru element at current frequency
			for elementNode := range frequencyList.FromTail() {
				if !yield(elementNode.key, elementNode.value) {
					return
				}
			}
		}
	}
}

func (l *сacheImpl[K, V]) Size() int {
	return l.size
}

func (l *сacheImpl[K, V]) Capacity() int {
	return l.capacity
}

func (l *сacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	elementNode, ok := l.items[key]
	// if key doesn't exist return error
	if !ok {
		return 0, ErrKeyNotFound
	}

	// if exists return frequency
	return elementNode.Value.freq, nil
}
