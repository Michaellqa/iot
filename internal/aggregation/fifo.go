package aggregation

type Fifo interface {
	Add(interface{})
	Get() interface{}
	Len() int
}

type node struct {
	prev *node
	next *node
	val  interface{}
}

// ListFifo queue based on double linked list
type ListFifo struct {
	head *node
	tail *node
	len  int
}

func (f *ListFifo) Add(elem interface{}) {
	n := &node{val: elem}
	if f.len == 0 {
		f.head = n
		f.tail = n
	} else {
		f.head.prev = n
		n.next = f.head
		f.head = n
	}
	f.len++
}

func (f *ListFifo) Get() interface{} {
	if f.len == 0 {
		return nil
	}

	elem := f.tail.val
	if f.len == 1 {
		f.head = nil
		f.tail = nil
	} else {
		f.tail = f.tail.prev
		f.tail.next = nil // feels enough to be GC-ed ...
	}
	f.len--
	return elem
}

func (f *ListFifo) Len() int {
	return f.len
}
