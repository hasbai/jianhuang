package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Supervisor struct {
	url           string
	isTLS         bool
	reqData       []byte
	addr          string
	contentLength int
}

func NewSupervisor(url string) *Supervisor {
	s := &Supervisor{
		url: url,
	}
	s.init()
	return s
}

func (s *Supervisor) init() {
	fmt.Println("init...")

	// make request data
	req, _ := http.NewRequest("GET", s.url, nil)
	req.Header.Set("User-Agent", UA)
	req.Header.Set("Connection", "keep-alive")
	s.reqData, _ = httputil.DumpRequestOut(req, false)

	// set isTLS
	if req.URL.Scheme == "https" {
		s.isTLS = true
	} else if req.URL.Scheme == "http" {
		s.isTLS = false
	} else {
		panic("invalid scheme")
	}

	// set addr
	port := req.URL.Port()
	if port == "" {
		if s.isTLS {
			port = "443"
		} else {
			port = "80"
		}
	}
	s.addr = req.URL.Hostname() + ":" + port

	// set contentLength
	s.contentLength = s.getContentLength()
}

func (s *Supervisor) Run(threads int) {
	for i := 0; i < threads; i++ {
		runner := NewRunner(s)
		go runner.Run()
	}
}

// getContentLength returns the length of the whole response
func (s *Supervisor) getContentLength() int {
	resp, err := http.Get(s.url)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode >= 400 {
		panic("request failed")
	}

	buf := bytes.NewBuffer(nil)
	_ = resp.Write(buf)
	return buf.Len()
}
