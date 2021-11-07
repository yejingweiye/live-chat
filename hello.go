package main

import (
	"fmt"
	"time"
)

func main() {

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case t := <-timer.C:
			fmt.Println("定时运行当前时间为", t)
			timer.Reset(3 * time.Second)
		}
	}

}
