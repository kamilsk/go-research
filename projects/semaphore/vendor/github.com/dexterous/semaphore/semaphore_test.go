package semaphore

import (
  "testing"
)


type A struct {
  *testing.T
}

func assert(t *testing.T) *A {
  return &A{t}
}

func (a *A) equal(expected uint, actual uint, message string) {
  if actual != expected {
    a.Errorf(message, expected, actual)
  }
}


func TestSemaphore_Capacity_OfNew(t *testing.T) {
  assert(t).equal(1, NewSemaphore().Capacity(), "Expected capacity %d got %d")
}

func TestSemaphore_Capacity_OfNewWithArg(t *testing.T) {
	assert(t).equal(5, NewSemaphoreWith(5).Capacity(), "Expected capacity of %d, got %d")
}

func TestSemaphore_QueueLength_OfNew(t *testing.T) {
  assert(t).equal(0, NewSemaphore().QueueLength(), "Expected queue of %d got %d")
}

func TestSemaphore_QueueLength_WithAcquiredPermit(t *testing.T) {
  var s = NewSemaphore()

  s.Acquire()

  assert(t).equal(1, s.QueueLength(),"Expected queue of %d got %d") 
}

func TestSemaphore_QueueLength_WithAcquiredPermitReleased(t *testing.T) {
  var s = NewSemaphoreWith(3)

  s.Acquire()
  s.Acquire()
  s.Release()

  assert(t).equal(1, s.QueueLength(), "Expected queue of %d got %d")
}

func TestSemaphore_Acquire_WithTimeout_AcquirePermit(t *testing.T) {
  var s = NewSemaphore()

  assert(t).equal(0, s.QueueLength(),"Expected queue of %d got %d") 

  if !s.TryAcquire() {
    t.Error("Could not acquire permit from Semaphore with spare")
  }

  assert(t).equal(1, s.QueueLength(),"Expected queue of %d got %d") 
}

func TestSemaphore_Acquire_WithTimeout_AcquireTimedout(t *testing.T) {
  var s = NewSemaphore()

  assert(t).equal(0, s.QueueLength(),"Expected queue of %d got %d") 

  s.Acquire()

  assert(t).equal(1, s.QueueLength(),"Expected queue of %d got %d") 

  if s.TryAcquire() {
    t.Error("Acquired permit from empty Semaphore")
  }

  assert(t).equal(1, s.QueueLength(),"Expected queue of %d got %d") 
}
