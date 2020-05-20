package main

import (
	"fmt"
	"sync"
	"time"
)

var done = false

func a(in <-chan int, out chan<- int, waitGroup *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Second)

	looping := true
	count := 0
	for looping {
		select {
		case t := <-ticker.C:
			fmt.Printf("Tick at %v\n", t)
			out <- count
			count++
		case i, b := <-in:
			fmt.Printf("a(%v, %v)\n", i, b)
			looping = b
		}
	}
	close(out)
	waitGroup.Done()
}

func b(in <-chan int, out chan<- int, waitGroup *sync.WaitGroup) {
	looping := true
	for looping {
		select {
		case i, b := <-in:
			fmt.Printf("b(%v, %v)\n", i, b)
			if b {
				out <- i
			} else {
				looping = false
			}
		}
	}
	close(out)
	waitGroup.Done()
}

func c(in <-chan int, out chan<- int, waitGroup *sync.WaitGroup) {
	looping := true
	for looping {
		select {
		case i, b := <-in:
			fmt.Printf("c(%v, %v)\n", i, b)
			if b == false { // error case that we never expect
				fmt.Printf("ERROR: in c and b is false!")
				looping = false
				break
			}

			if i >= 10 {
				fmt.Printf("c() is shutting everything down.\n")
				looping = false
			} else {
				out <- i
			}
		}
	}
	close(out)
	waitGroup.Done()
}

func main() {
	fmt.Printf("Starting up.\n")

	var waitGroup sync.WaitGroup
	channelAtoB := make(chan int)
	channelBtoC := make(chan int)
	channelCtoA := make(chan int)

	waitGroup.Add(3)
	go a(channelCtoA, channelAtoB, &waitGroup)
	go b(channelAtoB, channelBtoC, &waitGroup)
	go c(channelBtoC, channelCtoA, &waitGroup)
	waitGroup.Wait()

	fmt.Printf("fin\n")
}
