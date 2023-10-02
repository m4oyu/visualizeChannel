package main

import (
	"fmt"
	"time"
)

func main() {
	ch := chanx.Make(int, 1)

	go func() {
		for i := 0; i < 5; i++ {
			chanx.Send(i)
			time.Sleep(time.Second)
		}
		close(ch)
	}()

	for value := range ch {
		fmt.Println("Received:", value)
	}
}
