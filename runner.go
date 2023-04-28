package main

import (
	"io"
	"net/http"
)

const (
	NOTIFY = true
	EXIT   = false
)

type Runner struct {
	buf        []byte
	readBytes  int
	channel    chan bool
	supervisor *Supervisor
	client     *http.Client
}

const bufSize = 1500

func NewRunner(s *Supervisor) *Runner {
	return &Runner{
		buf:        make([]byte, bufSize),
		channel:    make(chan bool, 16),
		supervisor: s,
		client:     &http.Client{},
	}
}

func (r *Runner) Run() {
	resp := r.makeRequest()
	for {
		select {
		case i := <-r.channel:
			switch i {
			case NOTIFY:
				r.supervisor.speedCnt <- r.readBytes
				r.readBytes = 0
			case EXIT:
				return
			}
		default:
			n, err := resp.Body.Read(r.buf)
			if err == nil {
				r.readBytes += n
			} else {
				_ = resp.Body.Close()
				if err == io.EOF {
					resp = r.makeRequest()
				} else {
					panic("resp.Body.Read() error: " + err.Error())
				}
			}
		}
	}
}

func (r *Runner) makeRequest() *http.Response {
	resp, err := r.client.Get(url)
	if err != nil {
		panic("http.Get() error: " + err.Error())
	}
	return resp
}
