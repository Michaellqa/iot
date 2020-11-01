package aggregation

type node struct {
	prev *node
	next *node
	val  interface{}
}

// fifo queue based on double linked list
type fifo struct {
	head *node
	tail *node
	len  int
}

func (f *fifo) add(elem interface{}) {
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

func (f *fifo) get() interface{} {
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
