package linkedList

type Node[T any] struct {
	Value      T
	Next, Prev *Node[T]
}

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

// PushBack добавляет существующий узел в конец списка
func (l *LinkedList[T]) PushBack(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

	// Отключаем узел от предыдущего списка
	node.Next = nil
	node.Prev = nil

	if l.tail != nil {
		l.tail.Next = node
		node.Prev = l.tail
		l.tail = node
	} else {
		l.head = node
		l.tail = node
	}

	l.size++
}

// PushFront добавляет существующий узел в начало списка
func (l *LinkedList[T]) PushFront(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

	// Отключаем узел от предыдущего списка
	node.Next = nil
	node.Prev = nil

	if l.head != nil {
		l.head.Prev = node
		node.Next = l.head
		l.head = node
	} else {
		l.head = node
		l.tail = node
	}

	l.size++
}

// PushAfter вставляет существующий узел после указанного узла
func (l *LinkedList[T]) PushAfter(after *Node[T], node *Node[T]) {
	if after == nil || node == nil {
		panic("node is nil")
	}

	// Отключаем узел от предыдущего списка
	node.Next = nil
	node.Prev = nil

	node.Prev = after
	node.Next = after.Next

	if after.Next != nil {
		after.Next.Prev = node
	} else {
		l.tail = node
	}

	after.Next = node
	l.size++
}

func (l *LinkedList[T]) Pop(node *Node[T]) {
	if node == nil {
		panic("node is nil")
	}

	prev := node.Prev
	next := node.Next

	if prev != nil {
		prev.Next = next
	} else {
		l.head = next // если удаляемый узел - голова
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.tail = prev // обновляем l.tail при удалении хвоста
	}

	// Обнуляем ссылки в удаленном узле
	node.Next = nil
	node.Prev = nil

	l.size--
}

func (l *LinkedList[T]) PopBack() {
	if l.tail == nil {
		panic("list is empty")
	}

	oldTail := l.tail

	if l.tail.Prev != nil {
		l.tail = l.tail.Prev
		l.tail.Next = nil
	} else {
		l.head = nil
		l.tail = nil
	}

	// Обнуляем ссылки в удаленном узле
	oldTail.Next = nil
	oldTail.Prev = nil

	l.size--
}
