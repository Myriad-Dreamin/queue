package queue

import (
	"fmt"
	"runtime"
)

// UnsafeQueue64 is the special queue which value is uint64
type UnsafeQueue64 struct {
	capacity uint64
	capMod   uint64
	putPos   uint64
	getPos   uint64
	cache    []uint64
}

// NewUnsafeQueue64 return a pointer of UnsafeQueue64
// note: capacity <= 2^{63}
func NewUnsafeQueue64(capacity uint64) *UnsafeQueue64 {
	q := new(UnsafeQueue64)
	q.capacity = minQuantity64(capacity)
	q.capMod = q.capacity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]uint64, q.capacity)
	return q
}

// String return the information of the queue
func (q *UnsafeQueue64) String() string {
	return fmt.Sprintf("Queue{capacity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capacity, q.capMod, q.getPos, q.getPos)
}

// Capacity return the capacity of the queue
func (q *UnsafeQueue64) Capacity() uint64 {
	return q.capacity
}

// Quantity return the quantity of the queue
func (q *UnsafeQueue64) Quantity() uint64 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Size return the size of the queue
func (q *UnsafeQueue64) Size() uint64 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Put an uint64 to queue
func (q *UnsafeQueue64) Put(val uint64) bool {

	if q.putPos+2 == q.getPos {
		runtime.Gosched()
		return false
	}
	q.cache[q.putPos&q.capMod] = val
	q.putPos++
	return true
}

// Get an uint64 from queue
func (q *UnsafeQueue64) Get() (val uint64, ok bool) {
	if q.putPos == q.getPos {
		runtime.Gosched()
		return 0, false
	}
	val = q.cache[q.getPos&q.capMod]
	q.getPos++
	return val, true
}
