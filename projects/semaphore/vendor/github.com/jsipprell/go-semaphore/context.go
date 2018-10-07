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

package semaphore

import (
	"golang.org/x/net/context"
)

// An optional context-aware semaphore that is compatible with standard
// semaphores.  Usage:
//
//      s := semaphore.New(max)
//      ctx, sem := semaphore.WithContext(s, context)
//      // the above returns a new child context
//      release := <-sem.Acquire():
//      if release != nil {
//  	      // will be nil if context completed before acquiring
//            release()
//      }
//
// Context aware semaphores atomically acquire a sempahore as long as their
// parent context has not been cancelled. They return a cancellation function
// that will both release the sempahore and cancel the sub-context created for
// the sempahore. This cancellation function *must* be called in order to
// properly release the semaphore.
//
// If the parent context is cancelled before the semaphore is acquired,
// Acquire() returns nil.
//
// It is acceptable to call the cancellation function even if the parent or
// child contexts have already completed.
type ContextSemaphore interface {
	Acquire() <-chan context.CancelFunc
}

type acquireFunc func() <-chan context.CancelFunc

func (f acquireFunc) Acquire() <-chan context.CancelFunc {
	return f()
}

// Acquire a context-aware semaphore resource available on a channel.  Returns
// nil if the parent context was cancelled or completed before semaphore
// acquisition, othewise returns a context cancellation function which *MUST*
// be used to release the sempahore and cancel the subcontext.
//
// Note: returns a new child context which can be monitored just as with
// context.WithCacenl().
func WithContext(sem Semaphore, parent context.Context) (context.Context, ContextSemaphore) {
	var S *simpleSemaphore
	if s, ok := sem.(*simpleSemaphore); ok {
		S = s
	} else {
		panic("not a sempahore")
	}
	if parent == nil {
		parent = context.TODO()
	}
	ctx, cf := context.WithCancel(parent)
	c := make(chan context.CancelFunc, 1)
	go func(c chan<- context.CancelFunc) {
		select {
		case <-ctx.Done():
			close(c)
		case <-S.c:
			c <- context.CancelFunc(func() {
				defer cf()
				S.Release()
			})
		}
	}(c)
	return ctx, acquireFunc(func() <-chan context.CancelFunc {
		return c
	})
}
