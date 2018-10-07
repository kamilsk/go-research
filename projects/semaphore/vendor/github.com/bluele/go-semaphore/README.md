# go-semaphore

Implements several semaphore patterns on go language.

# Examples

## BasicSemaphore

`BasicSemaphore` is a semaphore manages a fixed number of concurrent task.

```go
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
}
```

## TimeLimitedSemaphore

`TimeLimitedSemaphore` is a semaphore manages a fixed number of concurrent task in time.

```go
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
}
```

## NamedSemaphores

NamedSemaphores is a semaphore mangeger can simply manage multiple named semaphores.

```go
func TestNamedSemaphores(t *testing.T) {
  permit := 3
  name := "test1"
  otherName := "test2"
  // You can use `TimeLimitedSemaphore` instead of `Semaphore`.
  // sem := semaphore.NewNamedSemaphores(semaphore.NewTimeLimitedSemaphore(permit, time.Second))
  sem := semaphore.NewNamedSemaphores(semaphore.NewSemaphore(permit))
  sem.Aquire(name, 1)
  if sem.Available(name) != permit-1 {
    t.Errorf(`sem.Available("%v") should be %v`, name, permit-1)
  }
  if sem.Available(otherName) != permit {
    t.Errorf(`sem.Available("%v") should be %v`, otherName, permit)
  }
  sem.Release(name)
  if sem.Available(name) != permit {
    t.Errorf(`sem.Available("%v") should be %v`, name, permit)
  }
}
```

# Test

```
$ go test
```

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>