package main

import (
	"fmt"
	"time"

	"github.com/m4oyu/visualizeChannel/chanx"

	viz "github.com/m4oyu/goroutine-viz"
)

func main() {
	// useChanx()
	callRuntimeStack()
}

func useChanx() {
	// 定義 chan, chan_id
	// ch1 := make(chan int)
	// ch2 := make(chan string)
	ch1 := chanx.Make(1, "ch1")
	ch2 := chanx.Make(1, "ch2")

	// 送信
	// ch1 <- 100
	// ch2 <- "hi"
	ch1.Send(100)

	go func() {
		ch2.Send("hi")
	}()

	// 受信
	// v1 := <-ch1
	// v2 := <-ch2
	v1, _ := ch1.Recv()
	v2, _ := ch2.Recv()

	fmt.Println(v1)
	fmt.Println(v2)

}

func wantToDo() {
	// 定義 chan, chan_id
	ch1 := make(chan int, 1)
	ch2 := make(chan string, 1)

	// 送信
	ch1 <- 100
	ch2 <- "hi"

	// 受信
	v1 := <-ch1
	v2 := <-ch2

	fmt.Println(v1)
	fmt.Println(v2)
}

func callRuntimeStack() {
	viz.WatchGoroutine("BREAKPOINT1")

	go func() {
		viz.WatchGoroutine("BREAKPOINT2")
		<-time.After(time.Second * 1)

	}()
}
