# semaphore
Very simple goro-safe semaphore

[![GoDoc](https://godoc.org/github.com/cognusion/semaphore?status.svg)](https://godoc.org/github.com/cognusion/semaphore)

Basics
======

```bash
go get github.com/cognusion/semaphore
```

Super simple: NewSemaphore(N) to create a semaphore of size N, Lock() to consume, Unlock() to replace, Free() to see how many locks are available. Additionally you can Add(i) and Sub(i) locks.

```go
package main

import (
	"github.com/cognusion/semaphore"
	"time"
	"fmt"
)

func main() {
	// Make a new semaphore, with the number of
	// simultaneous locks you want to allow
	S := semaphore.NewSemaphore(1)
	
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
```
