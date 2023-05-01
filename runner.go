package main

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
)

type Runner struct {
	supervisor *Supervisor
	buf        []byte
	conn       Conn
}

const bufSize = 1 << 20

func NewRunner(s *Supervisor) *Runner {
	return &Runner{
		buf:        make([]byte, bufSize),
		supervisor: s,
	}
}

func (r *Runner) run() {
	defer r.conn.Close()
	err := r.dial()
	if err != nil {
		panic(err.Error())
	}
	var n, length int
	for {
		n, err = r.conn.Read(r.buf)
		length += n
		if length >= r.supervisor.contentLength {
			return
		}
		if err != nil {
			if err == io.EOF {
				return
			} else {
				panic("Read error: " + err.Error())
			}
		}
	}
}

func (r *Runner) Run() {
	for {
		r.run()
	}
}

type Conn struct {
	netConn net.Conn
	conn    *tls.Conn
}

// Read avoid tls decryption to get better performance
func (c *Conn) Read(b []byte) (n int, err error) {
	return c.netConn.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.conn != nil {
		return c.conn.Write(b)
	}
	return c.netConn.Write(b)
}

func (c *Conn) Close() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			panic("close Conn error: " + err.Error())
		}
	} else {
		if c.netConn != nil {
			err := c.netConn.Close()
			if err != nil {
				panic("close Conn error: " + err.Error())
			}
		}
	}
}

func (r *Runner) dial() error {
	if r.supervisor.isTLS {
		conn, err := tls.Dial("tcp", r.supervisor.addr, &tls.Config{})
		if err != nil {
			return errors.New("dial error: " + err.Error())
		}
		r.conn = Conn{
			conn:    conn,
			netConn: conn.NetConn(),
		}
	} else {
		conn, err := net.Dial("tcp", r.supervisor.addr)
		if err != nil {
			return errors.New("dial error: " + err.Error())
		}
		r.conn = Conn{
			conn:    nil,
			netConn: conn,
		}
	}

	_, err := r.conn.Write(r.supervisor.reqData)
	if err != nil {
		return errors.New("write error: " + err.Error())
	}
	return nil
}
