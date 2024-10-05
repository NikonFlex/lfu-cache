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

func (l *LinkedList[T]) Head() *Node[T] {
	return l.head
}

func (l *LinkedList[T]) Tail() *Node[T] {
	return l.tail
}

func (l *LinkedList[T]) Size() int {
	return l.size
}

// PushBack pushes element to the end of the list, gets node to push
// panics if node is nil
func (l *LinkedList[T]) PushBack(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

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

	l.size++
}

// PushFront pushes element to the end of the list, gets node to push
// panics if node is nil
func (l *LinkedList[T]) PushFront(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

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

	l.size++
}

// PushAfter pushes element after the existing node (after) of the list, gets node to push
// panics if node or after is nil
func (l *LinkedList[T]) PushAfter(after *Node[T], node *Node[T]) {
	if after == nil || node == nil {
		panic("node is nil")
	}

	node.Next = nil
	node.Prev = nil

	node.Prev = after
	node.Next = after.Next

	// check that after is tail
	if after.Next != nil {
		after.Next.Prev = node
	} else {
		l.tail = node
	}

	after.Next = node
	l.size++
}

// Pop deletes node by node
// panics if node is nil
func (l *LinkedList[T]) Pop(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

	prev := node.Prev
	next := node.Next

	if prev != nil {
		prev.Next = next
	} else {
		l.head = next // if node to delete is head
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.tail = prev // if node to delete is tail
	}

	node.Next = nil
	node.Prev = nil

	l.size--
}
