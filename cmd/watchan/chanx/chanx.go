package chanx

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

var channelID int = 1

type msg struct {
	value interface{}
	ok    bool
}

type log struct {
	value     interface{}
	event     string
	goroutine string
	ok        bool
}

type c struct {
	mu     sync.Mutex
	cond   *sync.Cond
	origin chan msg
	logger chan log
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
func Make(length int) C {
	logger := make(chan log, length)

	c := &c{
		origin: make(chan msg, length),
		logger: logger,
		cid:    channelID,
	}

	go func() {
		for l := range c.logger {
			// "goroutine 1" SEND 20 via "channel 1"
			fmt.Printf("\"goroutine %v\" %v %v via \"channel %v\"\n", l.goroutine, l.event, l.value, c.cid)
		}
	}()

	channelID++
	c.cond = sync.NewCond(&c.mu)
	return c
}

func (c *c) Send(v interface{}) (ok bool) {
	defer func() { ok = recover() == nil }()
	c.origin <- msg{v, true}
	goroutineID := getGoroutineID()
	c.logger <- log{v, "SEND", goroutineID, true}
	return
}

func (c *c) Recv() (v interface{}, ok bool) {
	select {
	case m := <-c.origin:
		goroutineID := getGoroutineID()
		c.logger <- log{m.value, "RECV", goroutineID, true}
		return m.value, m.ok
	}

}

func (c *c) Close() (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer func() { ok = recover() == nil }()
	close(c.origin)
	close(c.logger)
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

func getGoroutineID() string {
	stackSlice := make([]byte, 2048)
	s := runtime.Stack(stackSlice, false)
	r := strings.Split(strings.Split(string(stackSlice[0:s]), "\n")[0], " ")
	return r[1]
}
