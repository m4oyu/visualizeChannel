package chanx

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

var IDCounter int = 1

type msg struct {
	v  interface{}
	ok bool
}

type c struct {
	mu     sync.Mutex
	cond   *sync.Cond
	c      chan msg
	closed bool
	cid    int
}

type C interface {
	// Send a messge to the channel. Returns false if the channel is closed.
	Send(v interface{}) (ok bool)
	// Recv a messge from the channel. Returns false if the channel is closed.
	Recv() (v interface{}, ok bool)
	// Close the channel. Returns false if the channel is already closed.
	Close() (ok bool)
	// Wait for the channel to close. Returns immediately if the channel is
	// already closed
	Wait()
}

// Make new channel. Provide a length to make a buffered channel.
func Make(length int, name string) C {
	c := &c{c: make(chan msg, length), cid: IDCounter}
	IDCounter++
	c.cond = sync.NewCond(&c.mu)
	return c
}

func (c *c) Send(v interface{}) (ok bool) {
	defer func() { ok = recover() == nil }()
	c.c <- msg{v, true}

	stackSlice := make([]byte, 2048)
	s := runtime.Stack(stackSlice, false)
	r := strings.Split(strings.Split(string(stackSlice[0:s]), "\n")[0], " ")
	fmt.Printf("%v %v SEND \"%v\" via channel %v\n", r[0], r[1], v, c.cid)

	return
}

func (c *c) Recv() (v interface{}, ok bool) {
	select {
	case msg := <-c.c:

		stackSlice := make([]byte, 2048)
		s := runtime.Stack(stackSlice, false)
		r := strings.Split(strings.Split(string(stackSlice[0:s]), "\n")[0], " ")
		fmt.Printf("%v %v RECV \"%v\" via channel %v\n", r[0], r[1], msg.v, c.cid)

		return msg.v, msg.ok
	}

}

func (c *c) Close() (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer func() { ok = recover() == nil }()
	close(c.c)
	c.closed = true
	c.cond.Broadcast()
	return
}

func (c *c) Wait() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for {
		if c.closed {
			return
		}
		c.cond.Wait()
	}
}
