package main

import (
	"fmt"
	"time"
)

type Supervisor struct {
	time     int
	speed    int
	accum    int
	avg      int
	speedCnt chan int
	runners  []*Runner
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		speedCnt: make(chan int, 128),
		runners:  make([]*Runner, 0, 128),
	}
}

func (s *Supervisor) AddRunner() {
	r := NewRunner(s)
	s.runners = append(s.runners, r)
	go r.Run()
}

func (s *Supervisor) Run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case i := <-s.speedCnt:
			s.speed += i
		case <-ticker.C:
			s.time++
			s.accum += s.speed
			s.avg = s.accum / s.time
			s.show()
			s.speed = 0
			s.notify()
		}
	}
}

func (s *Supervisor) notify() {
	for _, r := range s.runners {
		r.channel <- NOTIFY
	}
}

func (s *Supervisor) show() {
	if s.speed == 0 {
		return
	}
	fmt.Printf(
		"\rCur: %s/s  Avg: %s/s  Acu: %s   ",
		toSize(s.speed), toSize(s.avg), toSize(s.accum),
	)
}

const (
	K = 1 << 10
	M = 1 << 20
	G = 1 << 30
	T = 1 << 40
)

func toSize(size int) string {
	if size < K {
		return fmt.Sprintf("%dB", size)
	} else if size < M {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(K))
	} else if size < G {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(M))
	} else if size < T {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(G))
	} else {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(T))
	}
}
