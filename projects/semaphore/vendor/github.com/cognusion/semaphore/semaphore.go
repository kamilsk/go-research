/*
Package semaphore is a super simple goro-safe semaphore struct for Go.
* NewSemaphore(N) to create a semaphore of size N
* Lock() to consume
* Unlock() to replace
* Add(i) to add i to the lock count
* Sub(i) to subtract i to the lock count
* Free() to see how many locks are available

	import (
		"github.com/cognusion/semaphore"
		"time"
		"fmt"
	)

	func main() {
		// Make a new semaphore, with the number of
		// simultaneous locks you want to allow
		S := NewSemaphore(1)

		go func() {
			// Call lock, which will block if there aren't free locks
			// and defer the unlock until the function ends
			S.Lock()
			defer S.Unlock()

			// Do some stuff
			fmt.Println("Doing some stuff")
			time.Sleep(1 * time.Second)
		}()

		go func() {
			// Call lock, which will block if there aren't free locks
			// and defer the unlock until the function ends
			S.Lock()
			defer S.Unlock()

			// Do some other stuff
			fmt.Println("Doing some other stuff")
			time.Sleep(50 * time.Millisecond)
		}()

		time.Sleep(1 * time.Millisecond)
		fmt.Printf("Free locks? %d\n",S.Free())
		time.Sleep(3 * time.Second)
		fmt.Printf("Free locks now? %d\n",S.Free())
	}
*/
package semaphore

import "fmt"

// Semaphore is a goro-safe simple semaphore
type Semaphore struct {
	lock chan bool
}

// NewSemaphore returns a Semaphore allowing up to 'size' locks before blocking
func NewSemaphore(size int) Semaphore {
	var S Semaphore
	S.lock = make(chan bool, size)
	return S
}

// Lock consumes a lock in the semaphore, blocking if none is available
func (s *Semaphore) Lock() {
	s.lock <- true
}

// Unlock replaces a lock in the semaphore, blocking if no locks are consumed
func (s *Semaphore) Unlock() {
	<-s.lock
}

// Add consumes numLocks locks in the semaphore, blocking if none is available
func (s *Semaphore) Add(numLocks int) {
	for i := 0; i < numLocks; i++ {
		s.lock <- true
	}
}

// Sub replaces numLocks locks in the semaphore, blocking if no locks are consumed
func (s *Semaphore) Sub(numLocks int) {
	for i := 0; i < numLocks; i++ {
		<-s.lock
	}
}

// Free returns the number of available locks in the semaphore
func (s *Semaphore) Free() int {
	return cap(s.lock) - len(s.lock)
}

// String returns the string representation of the semaphore
func (s *Semaphore) String() string {
	return fmt.Sprintf("%d of %d free", s.Free(), cap(s.lock))
}
