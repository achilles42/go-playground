package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
		case <-time.After(time.Second):
			fmt.Println("running task")
		case <-s.quit:
			fmt.Println("finishing task")
			time.Sleep(time.Second)
			fmt.Println("task done")
			s.quit <- true
			return
		}
	}
}

func (s *Server) Stop() {
	fmt.Println("server stopping")
	s.quit <- true
	<-s.quit
	fmt.Println("server stopped")
}

func TestTicker() {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tict at ", t)
		}
	}()
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func OsIntrupt() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-sigs
			fmt.Println()
			fmt.Println(sig)
			done <- true
			return
		}
	}()
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

func main() {
	//  testLeakingGoRoutines()
	//  deadLock()
	//	s := NewServer()
	//	time.Sleep(12 * time.Second)
	//	s.Stop()
	//	TestTicker()
	OsIntrupt()
}
