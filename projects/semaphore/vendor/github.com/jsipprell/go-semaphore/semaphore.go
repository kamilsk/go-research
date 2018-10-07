/*
 * Copyright (c) 2014-2015 Jesse Sipprell <jessesipprell@gmail.com>
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 */

// This package implements a low-overhead idiomatic sempahore for go.
//
// The sempahore is entirely channel based and does not use goroutines.
//
// There is a global version and locally allocatable version, neither uses
// mutexes but the global version does use a sync.WaitGroup for syncronization.
package semaphore // "github.com/jsipprell/go-semaphore"

import (
	"errors"
	"sync"
	"time"
)

// The basic interface:
//
//     if err := sem.Acquire(); err != nil {
//         if err == semaphore.ErrUnavailable {
//             // ... resource exhausted
//         }
//     } else {
//         defer sem.Release()
//     }
//
// Methods:
//
// Acquire() error - acquire one resource count in a race-safe fashion
// Release()       - release one previously acquired resource count
// Cap()           - return the total maximum resource capacity
//                   (note: NOT what is available/in-use)
type Semaphore interface {
	Acquire() error
	Release()
	Cap() int
}

// To perform blocking, use the WaitableSemaphore interface
//
//     if err := waitsem.Wait(time.Duration(3300)); err != nil {
//         if err == semaphore.ErrUnavailable {
//             // ... resource exhausted
//         }
//     } else {
//         defer sem.Release()
//     }
//
// Additional Methods:
//
// Wait(duration time.Duration) error
//     Waits up time `duration` for the resource to become available and then
//     returns ErrUnavailable. If duration is 0, waits indefinitely.
type WaitableSemaphore interface {
	Semaphore
	Wait(time.Duration) error
}

var (
	// ErrUnavailable is the only error condition returned by this package
	ErrUnavailable = errors.New("semaphore resource unavailable")

	globalSem = sharedSemaphore{simpleSemaphore: &simpleSemaphore{}}
)

type simpleSemaphore struct {
	c chan struct{}
	// NB: capacity is for informational purposes only
	cap int
}

type sharedSemaphore struct {
	*simpleSemaphore
	init sync.Once
	wg   *sync.WaitGroup
}

// Returns the capacity of a semaphore, which is the max number
// of countable resources, NOT what's in-use or remaining.
func (sem *simpleSemaphore) Cap() int {
	return sem.cap
}

// Acquire() one resource unit from the semaphore, *never* blocks,
// but returns ErrUnavailable if exhausted.
func (sem *simpleSemaphore) Acquire() (err error) {
	select {
	case <-sem.c:
	default:
		err = ErrUnavailable
	}
	return
}

// Release a previously required resource unit. It is an error to
// release a resource not obtained via Acquire()
func (sem *simpleSemaphore) Release() {
	sem.c <- struct{}{}
}

func initSem(sem Semaphore) {
	s := sem.(*simpleSemaphore)
	s.c = make(chan struct{}, s.cap)
	for i := 0; i < s.cap; i++ {
		s.c <- struct{}{}
	}
}

func init() {
	if !GlobalIsInitialized() && globalSem.wg != nil {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		globalSem.wg = wg
	}
}

// Allocate a new private semaphore which can have Acquire()
// called on it `max` times without a matching Release() before it will
// return ErrUnavailable.
func New(max int) WaitableSemaphore {
	sem := &simpleSemaphore{cap: max}
	initSem(sem)
	return sem
}

// Returns the global shared semaphore, optionally initializing its
// maximum resource value (if this has not already been done)
func Global(max int) WaitableSemaphore {
	globalSem.init.Do(func() {
		if globalSem.wg == nil {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			globalSem.wg = wg
		}
		defer globalSem.wg.Done()
		if max > -1 && globalSem.cap == 0 {
			globalSem.cap = max
		}
		initSem(globalSem.simpleSemaphore)
	})

	globalSem.wg.Wait()
	return globalSem
}

// Returns true if the global shared semaphore has already been initialized
func GlobalIsInitialized() bool {
	return globalSem.simpleSemaphore.c != nil && globalSem.wg != nil
}

// Wait for a sempahore to have at least one resource available and then
// Acquire() it. If the time.Duration argument is zero, blocks forever.
//
// If the timer expires before the semaphore has Release() called on it
// by another goroutine, returns ErrUnavailable.
func (sem *simpleSemaphore) Wait(t time.Duration) (err error) {
	var C <-chan time.Time
	if t == time.Duration(0) {
		C = make(<-chan time.Time)
	} else {
		C = time.After(t)
	}
	select {
	case <-sem.c:
	case <-C:
		err = ErrUnavailable
	}
	return
}
