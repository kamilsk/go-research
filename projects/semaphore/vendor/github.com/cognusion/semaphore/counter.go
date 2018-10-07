package semaphore

// Counter is a Semaphore-locked integer
type Counter struct {
	lock    Semaphore
	counter int
}

// Inc[rements] the counter
func (c *Counter) Inc() {
	c.lock.Lock()
	c.counter++
	c.lock.Unlock()
}

// Dec[rements] the counter
func (c *Counter) Dec() {
	c.lock.Lock()
	c.counter--
	c.lock.Unlock()
}

// Value returns the current counter value
func (c *Counter) Value() int {
	return c.counter
}

// NewCounter returns a new Counter, with a 0 value
func NewCounter() Counter {
	return NewSetCounter(0)
}

// NewSetCounter returns a new Counter with the specified value
func NewSetCounter(value int) Counter {
	c := Counter{
		lock:    NewSemaphore(1),
		counter: value,
	}
	return c
}
