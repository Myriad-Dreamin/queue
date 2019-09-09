package queue

import (
	"fmt"
	"runtime"
)

// UnsafeBytesQueue is the special queue which value is bytes
type UnsafeBytesQueue struct {
	capacity uint32
	capMod   uint32
	putPos   uint32
	getPos   uint32
	cache    [][]byte
}

// NewUnsafeBytesQueue return a pointer of UnsafeBytesQueue
// note: capacity <= 2^{31}
func NewUnsafeBytesQueue(capacity uint32) *UnsafeBytesQueue {
	q := new(UnsafeBytesQueue)
	q.capacity = minQuantity(capacity)
	q.capMod = q.capacity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([][]byte, q.capacity)
	return q
}

// String return the information of the queue
func (q *UnsafeBytesQueue) String() string {
	return fmt.Sprintf("Queue{capacity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capacity, q.capMod, q.getPos, q.getPos)
}

// Capacity return the capacity of the queue
func (q *UnsafeBytesQueue) Capacity() uint32 {
	return q.capacity
}

// Quantity return the quantity of the queue
func (q *UnsafeBytesQueue) Quantity() uint32 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Size return the size of the queue
func (q *UnsafeBytesQueue) Size() uint32 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Put an bytes to queue
func (q *UnsafeBytesQueue) Put(val []byte) bool {

	if q.putPos+2 == q.getPos {
		runtime.Gosched()
		return false
	}
	q.cache[q.putPos&q.capMod] = val
	q.putPos++
	return true
}

// Get an bytes from queue
func (q *UnsafeBytesQueue) Get() (val []byte, ok bool) {
	if q.putPos == q.getPos {
		runtime.Gosched()
		return nil, false
	}
	val = q.cache[q.getPos&q.capMod]
	q.getPos++
	return val, true
}
