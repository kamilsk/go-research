package semaphore_test

import (
	"github.com/bluele/go-semaphore"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	permit := 3
	sem := semaphore.NewSemaphore(permit)
	sem.Aquire(1)
	if sem.Available() != permit-1 {
		t.Errorf("sem.Available() should be %v", permit-1)
	}
	sem.Release()
	if sem.Available() != permit {
		t.Errorf("sem.Available() should be %v", permit)
	}

	if !sem.AquireWithTimeout(time.Millisecond) {
		t.Errorf("sem.AquireWithTimeout(time.Millisecond) should be %v", true)
	}
	sem.Release()
	for i := 0; i < permit; i++ {
		sem.Aquire(1)
	}
	if sem.AquireWithTimeout(time.Millisecond) {
		t.Errorf("sem.AquireWithTimeout(time.Millisecond) should be %v", false)
	}
}

func TestTimeLimitedSemaphore(t *testing.T) {
	permit := 3
	sem := semaphore.NewTimeLimitedSemaphore(permit, time.Second)
	sem.Aquire(1)
	if sem.Available() != permit-1 {
		t.Errorf("sem.Available() should be %v", permit-1)
	}
	sem.Release()
	if sem.Available() != permit-1 {
		t.Errorf("sem.Available() should be %v", permit-1)
	}
	time.Sleep(2 * time.Second)
	if sem.Available() != permit {
		t.Errorf("sem.Available() should be %v", permit)
	}

	if !sem.AquireWithTimeout(time.Millisecond) {
		t.Errorf("sem.AquireWithTimeout(time.Millisecond) should be %v", true)
	}
	sem.Release()
	for i := 0; i < permit; i++ {
		sem.Aquire(1)
	}
	if sem.AquireWithTimeout(time.Millisecond) {
		t.Errorf("sem.AquireWithTimeout(time.Millisecond) should be %v", false)
	}
}

func TestNamedSemaphores(t *testing.T) {
	permit := 3
	name := "test1"
	otherName := "test2"
	sem := semaphore.NewNamedSemaphores(semaphore.NewSemaphore(permit))
	sem.Aquire(name, 1)
	if sem.Available(name) != permit-1 {
		t.Errorf(`sem.Available(%#v) should be %v`, name, permit-1)
	}
	if sem.Available(otherName) != permit {
		t.Errorf(`sem.Available(%#v) should be %v`, otherName, permit)
	}
	sem.Release(name)
	if sem.Available(name) != permit {
		t.Errorf(`sem.Available(%#v) should be %v`, name, permit)
	}

	sem = semaphore.NewNamedSemaphores(semaphore.NewTimeLimitedSemaphore(permit, time.Second))
	sem.Aquire(name, 1)
	if sem.Available(name) != permit-1 {
		t.Errorf(`sem.Available(%#v) should be %v`, name, permit-1)
	}
	if sem.Available(otherName) != permit {
		t.Errorf(`sem.Available(%#v) should be %v`, otherName, permit)
	}
	sem.Release(name)
	if sem.Available(name) != permit-1 {
		t.Errorf(`sem.Available(%#v) should be %v`, name, permit-1)
	}
	time.Sleep(2 * time.Second)
	if sem.Available(name) != permit {
		t.Errorf(`sem.Available(%#v) should be %v`, name, permit)
	}

	if !sem.AquireWithTimeout(name, time.Millisecond) {
		t.Errorf("sem.AquireWithTimeout(name, time.Millisecond) should be %v", true)
	}
	sem.Release(name)
	for i := 0; i < permit; i++ {
		sem.Aquire(name, 1)
	}
	if sem.AquireWithTimeout(name, time.Millisecond) {
		t.Errorf("sem.AquireWithTimeout(name, time.Millisecond) should be %v", false)
	}
}
