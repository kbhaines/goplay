package main

import (
	"fmt"
	"math/rand"
	"time"
)

type IdGenerator func() int

func NewIdGen() IdGenerator {
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

func main() {

	gen := NewIdGen()

	for i := 0; i < 100000; i++ {
		me := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			fmt.Printf("routine %d got %d\n", me, gen())
		}()
	}
	time.Sleep(time.Second * 2)

}
