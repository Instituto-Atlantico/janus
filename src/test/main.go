package main

import (
	"fmt"
	"time"
)

func main() {
	channel := make(chan string)

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			channel <- time.Now().GoString()
		}
	}()

	go func() {
		for {
			a := <-channel
			fmt.Println(a)
		}
	}()

	go func() {
		for {
			a := <-ticker.C
			fmt.Println(a)
		}
	}()

	for {
	}

	// ticker := time.NewTicker(1 * time.Second)

	// go func() {
	// 	for range ticker.C {
	// 		a := <-ticker.C
	// 		fmt.Println("done ", a)
	// 	}
	// }()

	// for {
	// }
}
