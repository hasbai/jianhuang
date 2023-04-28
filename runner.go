package main

import (
	"net/http"
)

const (
	NOTIFY = true
	EXIT   = false
)

type Runner struct {
	buf        []byte
	written    int
	channel    chan bool
	supervisor *Supervisor
}

const bufSize = 1500

func NewRunner(s *Supervisor) *Runner {
	return &Runner{
		buf:        make([]byte, bufSize),
		channel:    make(chan bool, 16),
		supervisor: s,
	}
}

func (r *Runner) Run() {
	resp := makeRequest()
	for {
		select {
		case i := <-r.channel:
			switch i {
			case NOTIFY:
				r.supervisor.speedCnt <- r.written
				r.written = 0
			case EXIT:
				return
			}
		default:
			n, err := resp.Body.Read(r.buf)
			if err == nil {
				r.written += n
			} else {
				_ = resp.Body.Close()
				if err.Error() == "EOF" {
					resp = makeRequest()
				} else {
					panic("resp.Body.Read() error: " + err.Error())
				}
			}
		}
	}
}

func makeRequest() *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		panic("http.Get() error: " + err.Error())
	}
	return resp
}
