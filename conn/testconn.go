package conn

import (
	"sync"
	"time"
)

type testConn struct {
	rd <-chan []byte
	wt chan<- []byte

	rlock *sync.Mutex
	wlock *sync.Mutex

	readDeadline  <-chan time.Time
	writeDeadline <-chan time.Time

	closed      chan struct{}
	wtCloseWait *sync.WaitGroup
}

func TestPair() (c1 Conn, c2 Conn) {
	a := make(chan []byte)
	b := make(chan []byte)

	c1 = &testConn{
		rd:          a,
		wt:          b,
		rlock:       &sync.Mutex{},
		wlock:       &sync.Mutex{},
		closed:      make(chan struct{}),
		wtCloseWait: &sync.WaitGroup{},
	}
	c2 = &testConn{
		rd:          b,
		wt:          a,
		rlock:       &sync.Mutex{},
		wlock:       &sync.Mutex{},
		closed:      make(chan struct{}),
		wtCloseWait: &sync.WaitGroup{},
	}

	return c1, c2
}

func (conn *testConn) isClosed() bool {
	select {
	case <-conn.closed:
		return true
	default:
		return false
	}
}

func (conn *testConn) Read() (msg []byte, err error) {
	if conn.isClosed() {
		return nil, ErrClosed
	}

	conn.rlock.Lock()
	defer conn.rlock.Unlock()

	select {
	case <-conn.closed:
		return nil, ErrClosed
	case <-conn.readDeadline:
		return nil, ErrReadTimeout
	case data, ok := <-conn.rd:
		if !ok {
			if !conn.isClosed() {
				conn.Close()
			}

			return data, ErrClosed
		}
		return data, nil
	}
}

func (conn *testConn) Write(msg []byte) error {
	if conn.isClosed() {
		return ErrClosed
	}

	conn.wlock.Lock()
	defer conn.wlock.Unlock()

	select {
	case <-conn.closed:
		return ErrClosed
	case <-conn.writeDeadline:
		return ErrWriteTimeout
	case conn.wt <- msg:
		return nil
	}
}

func (conn *testConn) Close() {
	if !conn.isClosed() {
		close(conn.closed)

		conn.wlock.Lock()
		close(conn.wt)
		conn.wlock.Unlock()
	}
}
func (conn *testConn) SetReadDeadline(t time.Time) error {
	if conn.isClosed() {
		return ErrClosed
	}

	conn.readDeadline = time.After(t.Sub(time.Now()))
	return nil
}
func (conn *testConn) SetWriteDeadline(t time.Time) error {
	if conn.isClosed() {
		return ErrClosed
	}

	conn.writeDeadline = time.After(t.Sub(time.Now()))
	return nil
}
