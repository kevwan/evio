package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"evio"
)

func main() {
	var lock sync.Mutex
	conns := make(map[string]evio.Conn)

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				var cc []evio.Conn
				lock.Lock()
				for _, conn := range conns {
					cc = append(cc, conn)
				}
				lock.Unlock()
				fmt.Println("conns:", len(conns))
				for _, conn := range cc {
					conn.Write([]byte("hello"))
				}
			}
		}
	}()

	var count int32
	var events evio.Events
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		out = in
		return
	}
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		lock.Lock()
		conns[c.RemoteAddr().String()] = c
		lock.Unlock()
		v := atomic.AddInt32(&count, 1)
		if v%100 == 0 {
			fmt.Println(v)
		}
		return
	}
	events.Closed = func(c evio.Conn, err error) (action evio.Action) {
		lock.Lock()
		delete(conns, c.RemoteAddr().String())
		lock.Unlock()
		v := atomic.AddInt32(&count, -1)
		if v%100 == 0 {
			fmt.Println(v)
		}
		return
	}
	if err := evio.Serve("tcp://0.0.0.0:2269", events); err != nil {
		panic(err.Error())
	}
}
