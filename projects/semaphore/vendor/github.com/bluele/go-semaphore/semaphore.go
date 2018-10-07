package semaphore

import (
	"sync"
	"time"
)

const (
	resource byte = 0
)

type Semaphore interface {
	Aquire(int)
	AquireWithTimeout(time.Duration) bool
	Release()
	Available() int
	Wait()
	deepCopy() Semaphore
}

type BasicSemaphore struct {
	channel chan byte
	permit  int
}

type TimeLimitedSemaphore struct {
	channel chan byte
	buffer  chan byte
	per     time.Duration
	permit  int
}

func NewSemaphore(permit int) *BasicSemaphore {
	sm := &BasicSemaphore{
		permit:  permit,
		channel: make(chan byte, permit),
	}
	for i := 0; i < permit; i++ {
		sm.channel <- resource
	}
	return sm
}

func NewTimeLimitedSemaphore(permit int, per time.Duration) *TimeLimitedSemaphore {
	sm := &TimeLimitedSemaphore{
		permit:  permit,
		per:     per,
		channel: make(chan byte, permit),
		buffer:  make(chan byte, permit),
	}
	for i := 0; i < permit; i++ {
		sm.channel <- resource
	}
	go sm.gc()
	return sm
}

func (sm *BasicSemaphore) Aquire(n int) {
	for i := 0; i < n; i++ {
		<-sm.channel
	}
}

func (sm *BasicSemaphore) AquireWithTimeout(timeout time.Duration) bool {
	ch := time.After(timeout)
	select {
	case <-sm.channel:
		return true
	case <-ch:
		return false
	}
}

func (sm *BasicSemaphore) Release() {
	select {
	case sm.channel <- resource:
	default:
		panic("too many release")
	}
}

func (sm *BasicSemaphore) Available() int {
	return len(sm.channel)
}

func (sm *BasicSemaphore) Wait() {
	sm.Aquire(sm.permit)
}

func (sm *BasicSemaphore) deepCopy() Semaphore {
	return NewSemaphore(sm.permit)
}

func (sm *TimeLimitedSemaphore) Aquire(n int) {
	for i := 0; i < n; i++ {
		<-sm.channel
	}
}

func (sm *TimeLimitedSemaphore) AquireWithTimeout(timeout time.Duration) bool {
	ch := time.After(timeout)
	select {
	case <-sm.channel:
		return true
	case <-ch:
		return false
	}
}

func (sm *TimeLimitedSemaphore) Release() {
	select {
	case sm.buffer <- resource:
	default:
		panic("too many release")
	}
}

func (sm *TimeLimitedSemaphore) Available() int {
	return len(sm.channel)
}

func (sm *TimeLimitedSemaphore) Wait() {
	sm.Aquire(sm.permit)
}

func (sm *TimeLimitedSemaphore) gc() {
	tick := time.NewTicker(sm.per)
	for {
	Main:
		select {
		case <-tick.C:
			for {
				select {
				case b := <-sm.buffer:
					sm.channel <- b
				default:
					break Main
				}
			}

		}
	}
}

func (sm *TimeLimitedSemaphore) deepCopy() Semaphore {
	return NewTimeLimitedSemaphore(sm.permit, sm.per)
}

type NamedSemaphores struct {
	model Semaphore
	sems  map[string]Semaphore
	mu    sync.RWMutex
}

func NewNamedSemaphores(model Semaphore) *NamedSemaphores {
	nss := &NamedSemaphores{}
	nss.model = model
	nss.sems = make(map[string]Semaphore)
	return nss
}

func (ns *NamedSemaphores) Aquire(name string, n int) {
	ns.mu.RLock()
	sm, ok := ns.sems[name]
	if ok {
		defer ns.mu.RUnlock()
		sm.Aquire(n)
		return
	}
	ns.mu.RUnlock()
	ns.createNamedSemaphore(name, func(sm Semaphore) interface{} {
		sm.Aquire(n)
		return nil
	})
}

func (ns *NamedSemaphores) AquireWithTimeout(name string, timeout time.Duration) bool {
	ns.mu.RLock()
	sm, ok := ns.sems[name]
	if ok {
		defer ns.mu.RUnlock()
		return sm.AquireWithTimeout(timeout)
	}
	ns.mu.RUnlock()
	return ns.createNamedSemaphore(name, func(sm Semaphore) interface{} {
		return sm.AquireWithTimeout(timeout)
	}).(bool)
}

func (ns *NamedSemaphores) Release(name string) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	sm, ok := ns.sems[name]
	if !ok {
		panic("no such name :'" + name + "'")
	}
	sm.Release()
}

func (ns *NamedSemaphores) Available(name string) int {
	ns.mu.RLock()
	sm, ok := ns.sems[name]
	if ok {
		defer ns.mu.RUnlock()
		return sm.Available()
	}
	ns.mu.RUnlock()
	return ns.createNamedSemaphore(name, func(sm Semaphore) interface{} {
		return sm.Available()
	}).(int)
}

func (ns *NamedSemaphores) Wait(name string) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	sm, ok := ns.sems[name]
	if !ok {
		panic("no such name :'" + name + "'")
	}
	sm.Wait()
}

func (ns *NamedSemaphores) createNamedSemaphore(name string, cb func(Semaphore) interface{}) interface{} {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	sm, ok := ns.sems[name]
	if !ok {
		sm = ns.model.deepCopy()
		ns.sems[name] = sm
	}
	return cb(sm)
}
