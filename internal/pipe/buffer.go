package pipe

import "sync"

type QueueItem struct {
	next     *QueueItem
	previous *QueueItem
	data     []byte
}

type Queue struct {
	cond   *sync.Cond
	head   *QueueItem
	tail   *QueueItem
	size   int
	closed bool
}

func NewQueue() *Queue {
	return &Queue{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (q *Queue) Enqueue(data []byte) bool {
	q.cond.L.Lock()

	if q.closed {
		q.cond.L.Unlock()
		return false
	}

	item := &QueueItem{
		data: data,
		next: q.tail,
	}

	if q.tail != nil {
		q.tail.previous = item
	}

	if q.head == nil {
		q.head = item
	}

	q.tail = item
	q.size += len(data)

	q.cond.Signal()
	q.cond.L.Unlock()
	return true
}

func (q *Queue) Dequeue() ([]byte, bool) {
	q.cond.L.Lock()

	for q.closed == false && q.head == nil {
		q.cond.Wait()
	}

	var data []byte
	if q.head != nil {
		item := q.head

		if q.head.previous != nil {
			q.head.previous.next = nil
			q.head = q.head.previous
		} else {
			q.head = nil
			q.tail = nil
		}

		item.next = nil
		item.previous = nil
		data = item.data
		q.size -= len(data)
	}

	drained := len(data) == 0 && q.closed && q.size == 0
	q.cond.L.Unlock()

	return data, drained
}

func (q *Queue) Close() {
	q.cond.L.Lock()
	q.closed = true
	q.cond.Broadcast()
	q.cond.L.Unlock()
}
