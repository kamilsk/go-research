Simple Go Semaphores
====================


Usage
=====

```go
var sem semaphore.Semaphore

// 5 resource units available max
sem = semaphore.New(5)
if err := sem.Acquire(); err != nil {
    if err == semaphore.ErrUnavailable {
        // ... resource exhausted
    } else {
        panic(err.Error());
    }
} else {
    defer sem.Release()
}
```


Global
======

A global shared singleton semaphore is available but not initialized until
first use:


```go
sem := semaphore.Global(10)
if err := sem.Acquire(); err != nil {
    if err == semaphore.ErrUnavailable {
        // ... resource exhausted
    } else {
        panic(err.Error());
    }
} else {
    defer sem.Release()
}
```
