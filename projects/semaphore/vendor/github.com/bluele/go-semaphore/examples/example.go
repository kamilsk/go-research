package main

import (
	"fmt"
	"github.com/bluele/go-semaphore"
	"time"
)

func task(i int) {
	time.Sleep(time.Millisecond * time.Duration(time.Now().UnixNano()%1000))
	fmt.Printf("done: task-%v\n", i)
}

func getMainCh() chan byte {
	ch := make(chan byte)
	go func() {
		taskNumber := 20
		permit := 6
		sem := semaphore.NewTimeLimitedSemaphore(permit, time.Second)

		for i := 0; i < taskNumber; i++ {
			sem.Aquire(1)
			go func(i int) {
				task(i)
				sem.Release()
			}(i)
		}
		sem.Wait()
		ch <- 0
	}()
	return ch
}

func main() {
	tick := time.NewTicker(time.Second)
	ch := getMainCh()
	counter := 0
Main:
	for {
		select {
		case <-ch:
			break Main
		case <-tick.C:
			counter += 1
			fmt.Println(counter)
		}
	}
}
