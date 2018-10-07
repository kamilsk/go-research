// tests for context-aware sempahores
package semaphore

import (
	"golang.org/x/net/context"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

type MockLogger interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

type LogAdapter struct {
	*log.Logger
}

type State struct {
	sync.Mutex
	MockLogger
	Context context.Context
	track   map[string]bool
}

func (L *LogAdapter) Log(v ...interface{}) {
	L.Println(v...)
}

func (L *LogAdapter) Logf(format string, v ...interface{}) {
	L.Printf(format, v...)
}

func (s *State) Track(ident string) {
	s.Lock()
	defer s.Unlock()
	if s.track == nil {
		s.track = make(map[string]bool)
	}
	s.track[ident] = false
}

func (s *State) Done(ident string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.track[ident]; ok {
		s.track[ident] = true
	} else {
		s.Fatalf("%v is not tracked", ident)
	}
}

func (s *State) Check() {
	<-time.After(time.Duration(50) * time.Millisecond)
	s.Lock()
	defer s.Unlock()

	if s.Context != nil {
		<-s.Context.Done()
	}
	for ident, ok := range s.track {
		if ok {
			s.Logf("%v: OKAY (completed)", ident)
		} else {
			s.Fatalf("%v: FAIL (did not complete)", ident)
		}
	}
}

func initLogger(prefix string) (L *LogAdapter) {
	if prefix != "" {
		prefix += "/"
	}
	return &LogAdapter{log.New(os.Stdout, prefix, 0)}
}

func acquireAndWaitWith(S Semaphore, state *State, secs int, timeout int, label string, t MockLogger) {
	_, sem := WithContext(S, state.Context)
	acquireAndWait(sem, state, secs, timeout, label, t)
}

func acquireAndWait(sem ContextSemaphore, state *State, secs int, timeout int, label string, t MockLogger) {
	var cf context.CancelFunc
	var TC <-chan time.Time
	var timer *time.Timer
	var wait bool
	defer func() {
		if cf != nil {
			cf()
			t.Logf("%s: ====== DONE ======", label)
		}
	}()
	state.Track(label)
	defer state.Done(label)
	if timeout > 0 {
		timer = time.NewTimer(time.Duration(timeout) * time.Second)
		TC = timer.C
	} else {
		TC = make(chan time.Time)
	}
	defer func() {
		if timer != nil {
			timer.Stop()
		}
		if cf != nil {
			cf()
		}
	}()
	ctx := state.Context
	if ctx == nil {
		ctx = context.TODO()
	} else {
		wait = true
	}
	select {
	case <-ctx.Done():
		t.Logf("%s: main context done early, aborting", label)
		return
	case cf = <-sem.Acquire():
		if cf == nil {
			if timeout != secs {
				t.Fatalf("%s: acquire returned unexpected nil cancelfunc, aborting", label)
			} else {
				t.Logf("%s: acquire reteurned expected nil cancelfunc", label)
			}
			return
		}
	case <-TC:
		t.Logf("%s: timer expired (%v seconds)", label, timeout)
		return
	}

	t.Logf("%s: ====== BEGIN ======", label)
	if timer == nil {
		timer = time.NewTimer(time.Duration(secs) * time.Second)
	}
	select {
	case <-timer.C:
	case <-ctx.Done():
		t.Logf("%s: Parent context cancelled", label)
	}
	//t.Logf("%s: releasing semaphore",label)
	if cf != nil {
		cf()
		cf = nil
		t.Logf("%s: ====== DONE ======", label)
	}
	if wait {
		t.Logf("%s: waiting for context completion", label)
		<-ctx.Done()
	}
}

func TestContextCancel(t *testing.T) {
	L := initLogger("cancel")
	L.Logf("starting parent context cancellation test, will take 5 seconds")
	parent, cf := context.WithCancel(context.TODO())
	state := &State{MockLogger: t, Context: parent}
	ctx, sem := WithContext(New(1), state.Context)
	if sem == nil {
		t.Fatal("nil sem")
	} else if sem.Acquire() == nil {
		t.Fatal("nil acquire")
	}
	go acquireAndWait(sem, state, 7, 7, "wait-7", t)
	go acquireAndWait(sem, state, 3, 20, "wait-20", t)
	acquireAndWait(sem, state, 5, 5, "wait-5", t)
	go cf()

	<-ctx.Done()
	state.Check()
}

func TestContextSync(t *testing.T) {
	L := initLogger("sync")
	L.Logf("starting context test, will take 20 seconds")
	ctx, cf := context.WithCancel(context.TODO())
	sem := New(1)
	state := &State{MockLogger: L, Context: ctx}
	go acquireAndWaitWith(sem, state, 3, 0, "a-wait-3", t)
	go acquireAndWaitWith(sem, state, 5, 0, "a-wait-5", t)
	go acquireAndWaitWith(sem, state, 3, 0, "b-wait-3", t)
	go acquireAndWaitWith(sem, state, 2, 0, "a-wait-2", t)
	go acquireAndWaitWith(sem, state, 7, 20, "a-wait-7", t)
	select {
	case <-time.After(time.Duration(20) * time.Second):
	case <-ctx.Done():
	}
	state.Logf("cancelling")
	cf()
	<-ctx.Done()
	state.MockLogger = t
	state.Logf("Finalized")
	state.Check()
}
