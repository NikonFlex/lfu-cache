package linkedlist

import "iter"

// Node to use in linkedList
// Value - value containing in list
// next - next node in list
// prev - prev node in list
type Node[T any] struct {
	Value      T
	next, prev *Node[T]
}

type LinkedList[T any] interface {
	// Head returns the node that is head of the list (real first element, not dummy)
	Head() *Node[T]

	// Tail returns the node that is tail of the list (real last element, not dummy)
	Tail() *Node[T]

	// Size returns the size of the list
	Size() int

	// PushBack pushes element to the end of the list, gets node to push
	// panics if node is nil
	PushBack(node *Node[T])

	// PushFront pushes element to the start of the list, gets node to push
	// panics if node is nil
	PushFront(node *Node[T])

	// PushAfter pushes element after the existing node (after) of the list, gets node to push
	// panics if node or after is nil
	PushAfter(after *Node[T], node *Node[T])

	// Pop deletes node by node
	// panics if node is nil
	Pop(node *Node[T])

	// FromTail iterates list from the tail
	FromTail() iter.Seq[T]
}

// linkedList realization
// head - dummy node in the start of the list
// tail - dummy node in the end of the list
// size - size of list
type linkedList[T any] struct {
	head, tail *Node[T]
	size       int
}

func New[T any]() *linkedList[T] {
	// init dummy nodes
	dummyHead := &Node[T]{}
	dummyTail := &Node[T]{}
	dummyHead.next = dummyTail
	dummyTail.prev = dummyHead
	return &linkedList[T]{head: dummyHead, tail: dummyTail}
}

func (l *linkedList[T]) Head() *Node[T] {
	// list is empty
	if l.size == 0 {
		return nil
	}
	return l.head.next
}

func (l *linkedList[T]) Tail() *Node[T] {
	// list is empty
	if l.size == 0 {
		return nil
	}
	return l.tail.prev
}

func (l *linkedList[T]) Size() int {
	return l.size
}

func (l *linkedList[T]) PushBack(node *Node[T]) {
	// incorrect input data
	if node == nil {
		panic("node is nil")
	}

	// Insert the node before the dummy tail
	l.insertAfter(l.tail.prev, node)
}

func (l *linkedList[T]) PushFront(node *Node[T]) {
	// incorrect input data
	if node == nil {
		panic("node is nil")
	}

	// Insert the node after the dummy head
	l.insertAfter(l.head, node)
}

func (l *linkedList[T]) PushAfter(after *Node[T], node *Node[T]) {
	// incorrect input data
	if after == nil || node == nil {
		panic("node is nil")
	}

	// Insert the node after the specified node
	l.insertAfter(after, node)
}

func (l *linkedList[T]) Pop(node *Node[T]) {
	// incorrect input data
	if node == nil {
		panic("node is nil")
	}

	// Remove the node from the list
	l.remove(node)
}

func (l *linkedList[T]) FromTail() iter.Seq[T] {
	return func(yield func(T) bool) {
		// start from tail
		for node := l.Tail(); node != nil && node != l.head; node = node.prev {
			if !yield(node.Value) {
				return
			}
		}
	}
}

// Helper method to insert a node after a specified node
func (l *linkedList[T]) insertAfter(after *Node[T], node *Node[T]) {
	// Prepare the node for insertion
	node.next = nil
	node.prev = nil

	// Insert node between 'after' and 'after.next'
	next := after.next
	after.next = node
	node.prev = after
	node.next = next
	next.prev = node

	// Increment size
	l.size++
}

// Helper method to remove a node from the list
func (l *linkedList[T]) remove(node *Node[T]) {
	prev := node.prev
	next := node.next

	// Reconnect prev and next nodes
	prev.next = next
	next.prev = prev

	// Clear node's links
	node.next = nil
	node.prev = nil

	// Decrement size
	l.size--
}
