package semaphore

type empty struct{}

type Semaphore chan empty

func New(capacity int) Semaphore {
	return make(Semaphore, capacity)
}

func (sem Semaphore) Wait() {
	sem <- empty{}
}

func (sem Semaphore) Signal() {
	<-sem
}

func (sem Semaphore) Finish() {
	for i := 0; i < cap(sem); i++ {
		sem.Wait()
	}
}
