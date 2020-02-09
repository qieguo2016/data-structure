/*
	lock free linked queue.
	ref:
		1. https://zhuanlan.zhihu.com/p/23863915
*/

package concurrent

import (
	"runtime"
	"sync/atomic"
	//"unsafe"
)

/* ArrayQueue
 * ringBuffer: fix length array with an increase index, use index mod cap to locate
 * availableBuffer: to mark data which can be read, when writing process is finished, is set to true
 * capability: equal to 2^n, so that bit operation can be use
 */
type ArrayQueue struct {
	ringBuffer      []*ArrayQueueNode
	availableBuffer []int32
	//writeCursor     *ArrayQueueCursor
	//readCursor      *ArrayQueueCursor
	writeCursor int64
	readCursor  int64
	capability  int64
	indexMark   int64
}

// ArrayQueueNode
type ArrayQueueNode struct {
	Value interface{}
}

/* ArrayQueueCursor
 * build cursor size same as cache line can prevent false sharing
 * the size of CPU cache line is 64 bytes
 */
//type ArrayQueueCursor [8]int64
//
//func (cursor *ArrayQueueCursor) Get() int64 {
//	return cursor[0]
//}
//
//func (cursor *ArrayQueueCursor) Set(val int64) {
//	cursor[0] = val
//}

func NewArrayQueue() *ArrayQueue {
	cap := int64(16) // 2 ^ n
	return &ArrayQueue{
		ringBuffer:      make([]*ArrayQueueNode, cap),
		availableBuffer: make([]int32, cap),
		//writeCursor:     &ArrayQueueCursor{},
		//readCursor:      &ArrayQueueCursor{},
		writeCursor: int64(0),
		readCursor:  int64(0),
		capability:  cap,
		indexMark:   cap - int64(1),
	}
}

func (queue *ArrayQueue) loadWriteCursor() int64 {
	//return (*ArrayQueueCursor)(atomic.LoadPointer(
	//	(*unsafe.Pointer)(unsafe.Pointer(&queue.writeCursor)),
	//))
	return atomic.LoadInt64(&queue.writeCursor)
}

func (queue *ArrayQueue) loadReadCursor() int64 {
	//return (*ArrayQueueCursor)(atomic.LoadPointer(
	//	(*unsafe.Pointer)(unsafe.Pointer(&queue.readCursor)),
	//))
	return atomic.LoadInt64(&queue.readCursor)
}

func (queue *ArrayQueue) casWriteCursor(oldV, newV int64) bool {
	//return (*ArrayQueueCursor)(atomic.LoadPointer(
	//	(*unsafe.Pointer)(unsafe.Pointer(&queue.writeCursor)),
	//))
	return atomic.CompareAndSwapInt64(&queue.writeCursor, oldV, newV)
}

func (queue *ArrayQueue) casReadCursor(oldV, newV int64) bool {
	//return (*ArrayQueueCursor)(atomic.LoadPointer(
	//	(*unsafe.Pointer)(unsafe.Pointer(&queue.readCursor)),
	//))
	return atomic.CompareAndSwapInt64(&queue.readCursor, oldV, newV)
}

// Enqueue 补充写写读读4并发时序图
func (queue *ArrayQueue) Enqueue(val interface{}) bool {
	// 申请空间, if writeCursor < readCursor+capability-2 还可插入，否则block（condition）
	//     PS:  writeCursor == readCursor+capability-2，此时writeCursor对应的位置还可能在读，不能更新
	// cas 更新writeCursor，解决写写冲突
	// 写ring buffer，写独享不需要cas
	// cas 更新available buffer，解决读写冲突

	var index, rc, wc int64
	for {
		rc = queue.loadReadCursor()
		wc = queue.loadWriteCursor()
		if wc-rc >= queue.capability-2 {
			runtime.Gosched()
			continue
		}
		if queue.casWriteCursor(wc, wc+1) {
			break
		}
	}
	index = wc & queue.indexMark
	if buf := queue.ringBuffer[index]; buf != nil {
		buf.Value = val
	} else {
		queue.ringBuffer[index] = &ArrayQueueNode{Value: val}
	}
	for {
		if atomic.CompareAndSwapInt32(&queue.availableBuffer[index], 0, 1) {
			break
		}
	}
	return true
}

func (queue *ArrayQueue) Dequeue() interface{} {
	// 判断可读，if readCursor < writeCursor 可读，否则已空block
	// cas 更新readCursor，解决读读冲突
	// 再判断available buffer是否true，true可读，否则自旋，解决读写冲突
	// 读ringbuffer，读独享不需要cas
	// 更新available buffer为false，不需要cas？
	var index, rc, wc int64
	for {
		rc = queue.loadReadCursor()
		wc = queue.loadWriteCursor()
		if rc >= wc {
			runtime.Gosched()
			continue
		}
		if queue.casReadCursor(rc, rc+1) {
			break
		}
	}
	index = rc & queue.indexMark
	for {
		if atomic.LoadInt32(&queue.availableBuffer[index]) == 1 {
			break
		}
	}

	buf := queue.ringBuffer[index]
	val := buf.Value
	for {
		if atomic.CompareAndSwapInt32(&queue.availableBuffer[index], 1, 0) {
			break
		}
	}
	return val
}
