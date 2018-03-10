package main

import (
	"fmt"
	"runtime"
	"time"
)

func query() int {
	n := 1
	time.Sleep(time.Duration(n) * time.Millisecond)
	return n
}

func queryAll() int {
	ch := make(chan int)
	go func() { ch <- query() }()
	go func() { ch <- query() }()
	go func() { ch <- query() }()
	return <-ch
}

func testLeakingGoRoutines() {
	fmt.Println("MainGoRoutine", runtime.NumGoroutine())
	for i := 0; i <= 4; i++ {
		queryAll()
		fmt.Println("GoRoutines", runtime.NumGoroutine())
	}
}

func deadLock() {
	var ch chan struct{}
	ch <- struct{}{}
}

//TestGoRoutine
type Server struct{ quit chan bool }

func NewServer() *Server {
	s := &Server{make(chan bool)}
	go s.run()
	return s
}

func (s *Server) run() {
	for {
		select {
		case <-s.quit:
			fmt.Println("finishing task")
			time.Sleep(time.Second)
			fmt.Println("task done")
			s.quit <- true
			return
		case <-time.After(time.Second):
			fmt.Println("running task")
		}
	}
}

func (s *Server) Stop() {
	fmt.Println("server stopping")
	s.quit <- true
	<-s.quit
	fmt.Println("server stopped")
}

func main() {
	//testLeakingGoRoutines()
	//deadLock()
	s := NewServer()
	time.Sleep(2 * time.Second)
	s.Stop()
}
