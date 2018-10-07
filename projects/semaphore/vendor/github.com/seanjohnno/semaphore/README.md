## Counting semaphore for golang

### Description

A locking construct. The semaphore maintains an internal count which starts at 0 and is:

* Incremented by *Signal()*
* Decremented by *Wait()*

If a goroutine calls *Wait()* and the count falls below zero then it'll block until theres a matching call to *Signal()*.

If *Signal()* is called multiple times before a *Wait()*, lets say an arbitrary 3x, then *Wait()* can be called 3x without blocking and return immediately 

### When would I want to use this?

Perhaps you have a channel or buffered channel and sending to this is blocking because it's not being read fast enough and the buffer is becoming full. Using a semaphore the writes are non-blocking (apart from the synchronisation within the semaphore methods) and more calls to *Signal()* mean the matching calls to *Wait()* are non-blocking

### Quick Example

```
  func main() {
    sem := semaphore.New()
    go waitAndSignal(sem)
    sem.Wait()
    fmt.Println("Exiting...")
  }
  
  func waitAndSignal(sem *semaphore.CountingSemaphore) {
    time.Sleep(time.Duration(time.Second))
    sem.Signal()
  }
```

### Full Example

You can see an example of a publisher/consumer queue using the semaphore [here](https://github.com/seanjohnno/goexamples/blob/master/semaphore_example.go). Continue reading if you require instructions on how to grab the sourcecode and/or example from within the command-line...

### Setup

Create your Go folder structure on the filesystem (if you have't already):

```
GoProjects
  |- src
  |- pkg
  |- bin
```
In your command-line set your **GOPATH** environment variable:

* Linux: `export GOPATH=<Replace_me_with_path_to>\GoProjects`
* Windows: `set GOPATH="<Replace_me_with_path_to>\GoProjects"`

Browse to your *GoProjects* folder in the command-line and enter:

  `go get github.com/seanjohnno/semaphore`

You should see the folders */github.com/seanjohnno/semaphore* under *src* and the code inside *semaphore*

If you want to run the example then make sure you're in your *GoProjects* folder and run:

  `go get github.com/seanjohnno/goexamples`

Navigate to the *goexamples directory* and run the following:

```
  go build semaphore_example.go
```

...and then depending on your OS:

* Linux: `./semaphore_example`
* Windows: `semaphore_example.exe`
  



