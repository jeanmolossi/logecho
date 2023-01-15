package logecho

import (
	"sync/atomic"
)

// safeCounter is a simple counter. It counts how many
// call are running concurrently and how many has running
type safeCounter struct {
	// counter is counting of running goroutines
	counter uint64
	// maxCounter is the counting max of concurrently
	// goroutines was run
	maxCounter uint64
}

// initialize singleton safeCounter
var c = &safeCounter{0, 0}

// compareAndSwap will check if current maxCounter still
// with old value and updates it to new value
func (s *safeCounter) compareAndSwap(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&s.maxCounter, old, new)
}

// values will load counter and maxCounter respective then return it
func (s *safeCounter) values() (uint64, uint64) {
	return atomic.LoadUint64(&s.counter), atomic.LoadUint64(&s.maxCounter)
}

// increment will adds 1 to counter value
func (s *safeCounter) increment() {
	atomic.AddUint64(&s.counter, 1)
}

// decrement will sub 1 to counter value
func (s *safeCounter) decrement() {
	atomic.AddUint64(&s.counter, ^uint64(0))
}

// incrementRequestCounter will increment counter and
// update maxCounter every counter pass maxCounter counting
func incrementRequestCounter() {
	c.increment()

	counter, maxCounter := c.values()
	if counter > maxCounter {
		c.compareAndSwap(maxCounter, counter)
	}
}

// CurrentCount will load counter value
func CurrentCount() uint64 {
	counter, _ := c.values()
	return counter
}

// TransactionCounter will load maxCounter value
func TransactionCounter() uint64 {
	_, maxCounter := c.values()
	return maxCounter
}

// decrementRequestCounter will decrement counter and
// update maxCounter when counter downs to zero
func decrementRequestCounter() {
	c.decrement()

	counter, maxCounter := c.values()
	if counter == 0 {
		c.compareAndSwap(maxCounter, 0)
	}
}
