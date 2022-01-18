package utils

import (
	"net"
	"time"
)

// TimeoutListener wraps a net.Listener, and gives a place to store the timeout
// parameters. On Accept, it will wrap the net.Conn with our own TimeoutConn for us.
type TimeoutListener struct {
	net.Listener
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewTimeoutListener read write timeout
func NewTimeoutListener(addr string, readTimeout, writeTimeout time.Duration) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	tl := &TimeoutListener{
		Listener:     l,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return tl, nil
}

// Accept wrap raw listener function
func (tl *TimeoutListener) Accept() (net.Conn, error) {
	c, err := tl.Listener.Accept()
	if err != nil {
		return nil, err
	}

	tc := &TimeoutConn{
		Conn:         c,
		ReadTimeout:  tl.ReadTimeout,
		WriteTimeout: tl.WriteTimeout,
	}

	return tc, nil
}

// TimeoutConn wraps a net.Conn, and sets a deadline
// for every read and write operation.
type TimeoutConn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Read wrap raw listener function
func (tc *TimeoutConn) Read(b []byte) (n int, err error) {
	if tc.ReadTimeout != 0 {
		tc.Conn.SetReadDeadline(time.Now().Add(tc.ReadTimeout))
	}
	return tc.Conn.Read(b)
}

func (tc *TimeoutConn) Write(b []byte) (n int, err error) {
	if tc.WriteTimeout != 0 {
		tc.Conn.SetWriteDeadline(time.Now().Add(tc.WriteTimeout))
	}
	return tc.Conn.Write(b)
}

// Close wrap raw conn function
func (tc *TimeoutConn) Close() error {
	return tc.Conn.Close()
}
