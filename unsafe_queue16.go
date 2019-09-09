package queue

import (
	"fmt"
	"runtime"
)

// UnsafeQueue16 is the special queue which value is uint16
type UnsafeQueue16 struct {
	capacity uint32
	capMod   uint16
	putPos   uint16
	getPos   uint16
	cache    []uint16
}

// NewUnsafeQueue16 return a pointer of UnsafeQueue16
// note: capacity <= 2^{15}
func NewUnsafeQueue16(capacity uint16) *UnsafeQueue16 {
	q := new(UnsafeQueue16)
	q.capacity = minQuantity(uint32(capacity))
	q.capMod = uint16(q.capacity - 1)
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]uint16, q.capacity)
	return q
}

// String return the information of the queue
func (q *UnsafeQueue16) String() string {
	return fmt.Sprintf("Queue{capacity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capacity, q.capMod, q.getPos, q.getPos)
}

// Capacity return the capacity of the queue
func (q *UnsafeQueue16) Capacity() uint32 {
	return q.capacity
}

// Quantity return the quantity of the queue
func (q *UnsafeQueue16) Quantity() uint16 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Size return the size of the queue
func (q *UnsafeQueue16) Size() uint16 {
	if q.putPos >= q.getPos {
		return q.putPos - q.getPos
	}
	return q.capMod + (q.putPos - q.getPos)
}

// Put an uint16 to queue
func (q *UnsafeQueue16) Put(val uint16) bool {

	if q.putPos+2 == q.getPos {
		runtime.Gosched()
		return false
	}
	q.cache[q.putPos&q.capMod] = val
	q.putPos++
	return true
}

// Get an uint16 from queue
func (q *UnsafeQueue16) Get() (val uint16, ok bool) {
	if q.putPos == q.getPos {
		runtime.Gosched()
		return 0, false
	}
	val = q.cache[q.getPos&q.capMod]
	q.getPos++
	return val, true
}
