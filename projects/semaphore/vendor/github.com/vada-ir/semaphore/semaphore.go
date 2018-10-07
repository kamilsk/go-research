// Package semaphore create a simple semaphore with channels.
package semaphore

import "time"

// Semaphore is used to handle permits
type Semaphore interface {
	// PermitCount return the number of permit in this semaphore
	PermitCount() int
	//Acquire try to acquire n permit, block until all permits are available
	Acquire(int)
	// TryAcquire try to acquire n permit in a time limit. return the number of
	// permit it get in duration
	TryAcquire(int, time.Duration) int
	// Release release the permits. may called from anothe routine, so maybe it
	// blocked if the number of released permit is not available for release
	Release(int)
	// TryRelease try to release permits and if no permit already taken, return
	// the number of actually relaesed permits
	TryRelease(int, time.Duration) int
}

type empty struct{}

type semaphore struct {
	c chan empty
}

func (s *semaphore) PermitCount() int {
	return cap(s.c)
}

func (s *semaphore) Acquire(n int) {
	for i := 0; i < n; i++ {
		s.c <- empty{}
	}
}

func (s *semaphore) TryAcquire(n int, d time.Duration) int {
	var total int
	dc := time.After(d)
	for i := 0; i < n; i++ {
		select {
		case s.c <- empty{}:
			total++
		case <-dc:
			return total
		}
	}

	return total
}

func (s *semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<-s.c
	}
}

func (s *semaphore) TryRelease(n int, d time.Duration) int {
	var total int
	dc := time.After(d)
	for i := 0; i < n; i++ {
		select {
		case <-s.c:
			total++
		case <-dc:
			return total
		}
	}

	return total
}

// NewSemaphore return a new semaphore with requested size
func NewSemaphore(size int) Semaphore {
	return &semaphore{make(chan empty, size)}
}
