// Simple demo of interleaving reads and writes from channels. Shows that the
// variable 'v' is not bound to the sending channel 'ch' until the channel is
// ready for communication (i.e. the 'reader' reads). At that time the value
// for 'v' is used in the communication. This is demonstrated by the result
// that the reader only sees every 5th value.
//
// The 'blip' channel is sent to every 0.25s, the reader channel 'ch' is read
// from every 1s, the net result is 'v' is incremented by 5 every 1s.

package main

import (
	"fmt"
	"time"
)

func server(ch chan<- int, blip <-chan struct{}) {
	v := 0
	for {
		fmt.Printf("waiting at %d ----- ", v)
		select {
		case ch <- v:
			fmt.Printf("sent %d\n", v)
			v++
		case <-blip:
			fmt.Println("blip")
			v++
		}
	}
}

func reader(ch <-chan int) {
	for {
		fmt.Printf("reader got %d\n", <-ch)
		time.Sleep(time.Second)
	}
}

func blipper(ch chan<- struct{}) {
	for {
		fmt.Println("blipping...")
		ch <- struct{}{}
		time.Sleep(time.Second / 4)
	}
}

func main() {
	ch := make(chan int)
	blip := make(chan struct{})

	go server(ch, blip)
	go reader(ch)
	go blipper(blip)
	time.Sleep(time.Second * 10)
}
