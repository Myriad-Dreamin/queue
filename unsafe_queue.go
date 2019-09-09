package queue

import (
	"fmt"
	"runtime"
)

// UnsafeQueue is the special queue which value is uint32
type UnsafeQueue struct {
	capacity uint32
	capMod   uint32
	putPos   uint32
	getPos   uint32
	cache    []uint32
}

// NewUnsafeQueue return a pointer of UnsafeQueue
// note: capacity <= 2^{31}
func NewUnsafeQueue(capacity uint32) *UnsafeQueue {
	q := new(UnsafeQueue)
	q.capacity = minQuantity(capacity)
	q.capMod = q.capacity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]uint32, q.capacity)
	return q
}

// String return the information of the queue
func (q *UnsafeQueue) String() string {
	return fmt.Sprintf("Queue{capacity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capacity, q.capMod, q.getPos, q.getPos)
}

// Capacity return the capacity of the queue
func (q *UnsafeQueue) Capacity() uint32 {
	return q.capacity
}

// Quantity return the quantity of the queue
func (q *UnsafeQueue) Quantity() uint32 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Size return the size of the queue
func (q *UnsafeQueue) Size() uint32 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Put an uint32 to queue
func (q *UnsafeQueue) Put(val uint32) bool {

	if q.putPos+2 == q.getPos {
		runtime.Gosched()
		return false
	}
	q.cache[q.putPos&q.capMod] = val
	q.putPos++
	return true
}

// Get an uint32 from queue
func (q *UnsafeQueue) Get() (val uint32, ok bool) {
	if q.putPos == q.getPos {
		runtime.Gosched()
		return 0, false
	}
	val = q.cache[q.getPos&q.capMod]
	q.getPos++
	return val, true
}
