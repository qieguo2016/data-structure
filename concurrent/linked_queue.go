package concurrent

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

/*
	lock free linked queue.
	ref:
		1. http://ddrv.cn/a/591069
		2. https://coolshell.cn/articles/8239.html
*/

type ConcurrentLinkedQueueNode struct {
	Value interface{}
	Next  *ConcurrentLinkedQueueNode
}

func (node *ConcurrentLinkedQueueNode) casNext(oldV, newV *ConcurrentLinkedQueueNode) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&node.Next)),
		unsafe.Pointer(oldV),
		unsafe.Pointer(newV),
	)
}

func (node *ConcurrentLinkedQueueNode) loadNext() *ConcurrentLinkedQueueNode {
	return (*ConcurrentLinkedQueueNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&node.Next)),
	))
}

type ConcurrentLinkedQueue struct {
	head *ConcurrentLinkedQueueNode
	tail *ConcurrentLinkedQueueNode
	size int64
	m    sync.Mutex
}

func NewConcurrentLinkedQueue() *ConcurrentLinkedQueue {
	dummy := &ConcurrentLinkedQueueNode{}
	dummy.Value = nil
	dummy.Next = nil
	return &ConcurrentLinkedQueue{ // like container/list, use same node
		head: dummy,
		tail: dummy,
	}
}

func (queue *ConcurrentLinkedQueue) casTail(oldV, newV *ConcurrentLinkedQueueNode) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&queue.tail)),
		unsafe.Pointer(oldV),
		unsafe.Pointer(newV),
	)
}

func (queue *ConcurrentLinkedQueue) casHead(oldV, newV *ConcurrentLinkedQueueNode) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&queue.head)),
		unsafe.Pointer(oldV),
		unsafe.Pointer(newV),
	)
}

func (queue *ConcurrentLinkedQueue) loadHead() *ConcurrentLinkedQueueNode {
	return (*ConcurrentLinkedQueueNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&queue.head)),
	))
}

func (queue *ConcurrentLinkedQueue) loadTail() *ConcurrentLinkedQueueNode {
	return (*ConcurrentLinkedQueueNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&queue.tail)),
	))
}

func (queue *ConcurrentLinkedQueue) Enqueue(v interface{}) bool {
	newNode := &ConcurrentLinkedQueueNode{Value: v, Next: nil}
	var tail, next *ConcurrentLinkedQueueNode
	for {
		// use atomic load and cas
		tail = queue.loadTail()
		next = tail.loadNext()
		if tail == queue.loadTail() { // double check
			if next == nil { // queue tail
				if tail.casNext(next, newNode) { // link to queue
					break
				}
			} else {
				queue.casTail(tail, next) // move tail pointer to real tail
			}
		}
	}

	queue.casTail(tail, newNode) // failure is ok, another thread has update
	atomic.AddInt64(&queue.size, 1)
	return true
}

func (queue *ConcurrentLinkedQueue) Dequeue() interface{} {
	var head, tail, first *ConcurrentLinkedQueueNode
	for {
		// use atomic load and cas
		head = queue.loadHead()  // dummy
		tail = queue.loadTail()  // dummy
		first = head.loadNext()  // nil
		if head == queue.loadHead() { // double check
			if first == nil { // empty list
				return nil
			}
			if head == tail { // empty list
				queue.casTail(tail, first) // move tail to real pointer
				continue
			}
			if queue.casHead(head, first) {
				break
			}
		}
	}

	atomic.AddInt64(&queue.size, -1)
	return first.Value
}

func (queue *ConcurrentLinkedQueue) Size() int64 {
	return atomic.LoadInt64(&queue.size)
}

func (queue *ConcurrentLinkedQueue) EnqueueLock(v interface{}) bool {
	newNode := &ConcurrentLinkedQueueNode{Value: v, Next: nil}
	queue.m.Lock()
	defer queue.m.Unlock()
	tail := queue.tail
	tail.Next = newNode
	queue.tail = newNode
	queue.size += 1
	return true
}

func (queue *ConcurrentLinkedQueue) DequeueLock() interface{} {
	var head, tail, first *ConcurrentLinkedQueueNode
	queue.m.Lock()
	defer queue.m.Unlock()
	head = queue.head
	tail = queue.tail
	first = head.Next
	if head == tail {
		return nil
	}
	queue.head = first
	head.Next = nil
	queue.size -= 1
	return first.Value
}
