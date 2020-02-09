package concurrent

/*
	lock free linked queue.
	ref:
		1. http://ddrv.cn/a/591069
		2. https://coolshell.cn/articles/8239.html
*/

/* ArrayQueueCursor
 * build cursor size same as cache line can prevent false sharing
 * the size of CPU cache line is 64 bytes
 */
type ArrayQueueCursor [8]int64

/* ArrayQueue
 * ringBuffer: fix length array with an increase index, use index mod cap to locate
 * availableBuffer: to mark data which can be read, when writing process is finished, is set to true
 * capability: equal to 2^n, so that bit operation can be use
 */
type ArrayQueue struct {
	ringBuffer      []interface{}
	availableBuffer []bool
	writeCursor     ArrayQueueCursor
	readCursor      ArrayQueueCursor
	capability      int64
	indexMark       int64
}

func NewArrayQueue() *ArrayQueue {
	cap := int64(16)
	return &ArrayQueue{
		ringBuffer:      make([]interface{}, cap),
		availableBuffer: make([]bool, cap),
		writeCursor:     ArrayQueueCursor{},
		readCursor:      ArrayQueueCursor{},
		capability:      cap,
		indexMark:       cap - int64(1),
	}
}

func (queue *ArrayQueue) Enqueue()  {
	
}

func (queue *ArrayQueue) Dequeue()  {
	
}