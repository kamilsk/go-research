package semaphore

import (
	"sync"
)

// ------------------------------------------------------------------------------------------------------------------------
// Struct: CountingSemaphore
// ------------------------------------------------------------------------------------------------------------------------

// CountingSemaphore is a locking mechanism which counts Wait/Signal calls and only locks if Count goes below zero
//
// For example:
// 	CountingSemaphore.Signal() // Count = 1
// 	CountingSemaphore.Signal() // Count = 2
//  CountingSemaphore.Wait() // Count = 1, the method returns immediately
// 	CountingSemaphore.Wait() // Count = 0, the method returns immediately
//  CountingSemaphore.Wait() // Count = -1, this call will block until Signal is called and the count returns to zero
type CountingSemaphore struct {
	
	// Count is the current count of Wait/Signal. Signal increments, Wait decrements
	Count int

	// SyncLock is used to synchronise the methods of this objects
	SyncLock sync.Mutex

	// WaitingChan is the underlying channel that a Wait() call will block on if Count < 0
	WaitingChan chan bool
}

// Wait decrements the Count by 1 and will pause/block on the underlying lock if the count goes below zero
func (this *CountingSemaphore) Wait() {

	// Sync method as we're altering state
	this.SyncLock.Lock()

	this.Count--
	if this.Count < 0 {

		// If count has gone below zero then we need to release sync lock and wait on channel (a call to Signal will wake us up)
		this.SyncLock.Unlock()
		<- this.WaitingChan

	} else {
		this.SyncLock.Unlock()
	}
}

// TryAcquire is non blocking. It'll return true and decrement the count if its currently positive, false if Count <= 0
func (this *CountingSemaphore) TryAcquire() bool {

	// Sync method as we're altering state
	this.SyncLock.Lock()
	defer this.SyncLock.Unlock()

	// If we're < 1 then Wait would block, return false
	if this.Count < 1 {
		return false

	// We've had enough signals so decrement count & return true
	} else {
		this.Count--
		return true
	}
}

// Signal increments the Count by 1 and will unblock a single goroutine blocked on the Wait call
func (this *CountingSemaphore) Signal() {

	// Sync method as we're altering state
	this.SyncLock.Lock()
	defer this.SyncLock.Unlock()

	// If count is below zero then we have goroutines waiting on this channel (it won't block) so wake them up
	if this.Count < 0 {
		this.WaitingChan <- true
	}

	// Increment count
	this.Count++
}

// ------------------------------------------------------------------------------------------------------------------------
// Construction
// ------------------------------------------------------------------------------------------------------------------------

// New creates a new CountingSemaphore
func New() *CountingSemaphore {
	cs := &CountingSemaphore{ WaitingChan: make(chan bool)}
	return cs
}