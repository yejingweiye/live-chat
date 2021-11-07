package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	go func() {
		close(ch1)
	}()

	time.Sleep(1 * time.Second)
	for {
		select {
		case v, ok := <-ch1:
			fmt.Println("close", v, ok)
		}
	}
}
