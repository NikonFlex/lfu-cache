package linkedlist

// Node to use in LinkedList
type Node[T any] struct {
	Value      T
	Next, Prev *Node[T]
}

// LinkedList realization
type LinkedList[T any] struct {
	head, tail *Node[T]
	size       int
}

func New[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Head returns the node that is head of the list
func (l *LinkedList[T]) Head() *Node[T] {
	return l.head
}

// Tail returns the node that is tail of the list
func (l *LinkedList[T]) Tail() *Node[T] {
	return l.tail
}

// Size returns the size of the list
func (l *LinkedList[T]) Size() int {
	return l.size
}

// PushBack pushes element to the end of the list, gets node to push
// panics if node is nil
func (l *LinkedList[T]) PushBack(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

	// to be sure that this node doesn't link to other nodes in other lists
	node.Next = nil
	node.Prev = nil

	// if tail != nil some elements are contained in the list
	if l.tail != nil {
		l.tail.Next = node
		node.Prev = l.tail
		l.tail = node
	} else {
		// list is empty
		// init list
		l.head = node
		l.tail = node
	}

	// increment size
	l.size++
}

// PushFront pushes element to the end of the list, gets node to push
// panics if node is nil
func (l *LinkedList[T]) PushFront(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

	// to be sure that this node doesn't link to other nodes in other lists
	node.Next = nil
	node.Prev = nil

	// if head != nil some elements are contained in the list
	if l.head != nil {
		l.head.Prev = node
		node.Next = l.head
		l.head = node
	} else {
		// list is empty
		// init list
		l.head = node
		l.tail = node
	}

	// increment size
	l.size++
}

// PushAfter pushes element after the existing node (after) of the list, gets node to push
// panics if node or after is nil
func (l *LinkedList[T]) PushAfter(after *Node[T], node *Node[T]) {
	if after == nil || node == nil {
		panic("node is nil")
	}

	// to be sure that this node doesn't link to other nodes in other lists
	node.Next = nil
	node.Prev = nil

	// set links
	node.Prev = after
	node.Next = after.Next

	// check that after is tail
	if after.Next != nil {
		after.Next.Prev = node
	} else {
		l.tail = node
	}

	after.Next = node
	// increment size
	l.size++
}

// Pop deletes node by node
// panics if node is nil
func (l *LinkedList[T]) Pop(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

	// for convenient use
	prev := node.Prev
	next := node.Next

	// set links for next
	if prev != nil {
		prev.Next = next
	} else {
		l.head = next // if node to delete is head
	}

	// set link for prev
	if next != nil {
		next.Prev = prev
	} else {
		l.tail = prev // if node to delete is tail
	}

	// ensure that the node doesn't link to other nodes
	node.Next = nil
	node.Prev = nil

	// decrement size
	l.size--
}
