package queue

import (
	"fmt"
	"runtime"
)

// UnsafeIQueue is the special queue which value is interface
type UnsafeIQueue struct {
	capacity uint32
	capMod   uint32
	putPos   uint32
	getPos   uint32
	cache    []interface{}
}

// NewUnsafeInterfaceQueue return a pointer of UnsafeIQueue
// note: capacity <= 2^{31}
func NewUnsafeInterfaceQueue(capacity uint32) *UnsafeIQueue {
	q := new(UnsafeIQueue)
	q.capacity = minQuantity(capacity)
	q.capMod = q.capacity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]interface{}, q.capacity)
	return q
}

// String return the information of the queue
func (q *UnsafeIQueue) String() string {
	return fmt.Sprintf("Queue{capacity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capacity, q.capMod, q.getPos, q.getPos)
}

// Capacity return the capacity of the queue
func (q *UnsafeIQueue) Capacity() uint32 {
	return q.capacity
}

// Quantity return the quantity of the queue
func (q *UnsafeIQueue) Quantity() uint32 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Size return the size of the queue
func (q *UnsafeIQueue) Size() uint32 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Put an interface to queue
func (q *UnsafeIQueue) Put(val interface{}) bool {

	if q.putPos+2 == q.getPos {
		runtime.Gosched()
		return false
	}
	q.cache[q.putPos&q.capMod] = val
	q.putPos++
	return true
}

// Get an interface from queue
func (q *UnsafeIQueue) Get() (val interface{}, ok bool) {
	if q.putPos == q.getPos {
		runtime.Gosched()
		return nil, false
	}
	val = q.cache[q.getPos&q.capMod]
	q.getPos++
	return val, true
}
