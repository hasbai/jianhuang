package main

import (
	"io"
	"net/http"
)

type Runner struct {
	buf       []byte
	readBytes int
	client    *http.Client
}

const bufSize = 1500

func NewRunner() *Runner {
	return &Runner{
		buf:    make([]byte, bufSize),
		client: &http.Client{},
	}
}

func (r *Runner) Run() {
	resp := r.makeRequest()
	for {
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

func (r *Runner) makeRequest() *http.Response {
	resp, err := r.client.Get(url)
	if err != nil {
		panic("http.Get() error: " + err.Error())
	}
	return resp
}
