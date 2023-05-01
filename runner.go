package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
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

func request() {
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "teicn.oss-cn-hongkong.aliyuncs.com:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			panic("close conn error: " + err.Error())
		}
	}(conn)

	req, _ := http.NewRequest("GET", "https://teicn.oss-cn-hongkong.aliyuncs.com/teicarmx64.7z", nil)
	req.Header.Set("User-Agent", UA)
	reqData, _ := httputil.DumpRequestOut(req, false)

	n, err := conn.Write(reqData)
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 1024)
	n, err = conn.NetConn().Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))
}

func getReqDataAndAddr(url string) ([]byte, string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UA)
	reqData, _ := httputil.DumpRequestOut(req, false)
	var port int
	if req.URL.Scheme == "https" {
		port = 443
	} else if req.URL.Scheme == "http" {
		port = 80
	} else {
		panic("invalid scheme")
	}
	addr := req.URL.Hostname() + ":" + strconv.Itoa(port)
	return reqData, addr
}
