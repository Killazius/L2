package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	res := make(chan interface{})
	go func() {
		defer close(res)
		switch len(channels) {
		case 0:
			return
		case 1:
			<-channels[0]
		default:
			select {
			case <-channels[0]:
			case <-or(channels[1:]...):
			}
		}
	}()
	return res
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
			fmt.Println("stopped after", after)
		}()
		return c
	}

	start := time.Now()

	res := or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	<-res
	fmt.Printf("done after %v", time.Since(start))
}
