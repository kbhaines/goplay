// package counter provides two examples of unique id generators using channels
// to control access to the counter.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// IdGenFunc is a func that returns a unique Id
type IdGenFunc func() int

// NewIdGen returns an IdGenFunc that, when called, gives an incrementing Id
// guaranteed to be unique within the program. It leaks one go-routine.
func NewIdGen() IdGenFunc {
	ch := make(chan int)
	server := func() {
		var counter int
		for {
			ch <- counter
			counter++
		}
	}
	client := func() int {
		return <-ch
	}

	go server()
	return client
}

// NewIdGenerator returns a new IdGenerator object that may be used to obtain
// program-wide unique and incrementing Ids. An IdGenerator attempts to not
// leak any Go routines.
func NewIdGenerator() *IdGenerator {
	idgen := &IdGenerator{ch: make(chan int), done: make(chan struct{})}

	go func() {
		for {
			select {
			case idgen.ch <- idgen.counter:
				idgen.counter++
			case <-idgen.done:
				close(idgen.ch)
				return
			}
		}
	}()

	return idgen
}

type IdGenerator struct {
	counter int
	ch      chan int
	done    chan struct{}
}

// Get returns a unique integer, next in the sequence. If Done() has already
// been called on this generator then Get returns 0.
func (i *IdGenerator) Get() int {
	return <-i.ch
}

// Done tells the generator to close down.
func (i *IdGenerator) Done() {
	close(i.done)
}

func main() {

	gen := NewIdGenerator()
	limit := 100000

	var wg sync.WaitGroup
	wg.Add(limit)
	for i := 0; i < limit; i++ {
		me := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			fmt.Printf("routine %d got %d\n", me, gen.Get())
			wg.Done()
		}()
	}
	wg.Wait()
	gen.Done()
}
