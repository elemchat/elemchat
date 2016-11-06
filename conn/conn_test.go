package conn

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type ConnPairFunc func() (Conn, Conn)

func TestConn(t *testing.T) {
	t.Run("TestPair", func(t *testing.T) {
		subtest_Conn(t, func() (Conn, Conn) {
			return TestPair()
		})
	})
}

func subtest_Conn(t *testing.T, pairFn ConnPairFunc) {
	if pairFn == nil {
		t.Error("ConnPairFunc equal nil")
		return
	}

	wait := func(fs func(s Conn), fc func(c Conn)) {
		s, c := pairFn()
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			fs(s)
		}()

		go func() {
			defer wg.Done()
			fc(c)
		}()
		wg.Wait()
	}

	wait(func(s Conn) {
		s.SetReadDeadline(time.Now().Add(time.Second * 1))
		msg, err := s.Read()
		if err != nil {
			t.Error(err)
			return
		}
		if string(msg) != "normal" {
			t.Error("no normal")
		}
	}, func(c Conn) {
		c.SetWriteDeadline(time.Now().Add(time.Second * 1))
		err := c.Write([]byte("normal"))
		if err != nil {
			t.Error(err)
			return
		}
	})

	wait(func(s Conn) {
		s.SetReadDeadline(time.Now().Add(time.Second * 1))
		_, err := s.Read()
		if err != ErrReadTimeout {
			t.Error(fmt.Sprintf("expect ErrReadTimeout got %v", err))
			return
		}

	}, func(c Conn) {
	})

	wait(func(s Conn) {
		s.SetWriteDeadline(time.Now().Add(time.Second * 1))
		err := s.Write([]byte("write timeout"))
		if err != ErrWriteTimeout {
			t.Error(fmt.Sprintf("expect ErrWriteTimeout got %v", err))
			return
		}
	}, func(c Conn) {
	})

	wait(func(s Conn) {
		s.Close()
		err := s.SetWriteDeadline(time.Now())
		if err != ErrClosed {
			t.Error(fmt.Sprintf("expect ErrClosed got %v", err))
			return
		}

		err = s.SetReadDeadline(time.Now())
		if err != ErrClosed {
			t.Error(fmt.Sprintf("expect ErrClosed got %v", err))
			return
		}
	}, func(c Conn) {
	})

	wait(func(s Conn) {
		_, err := s.Read()
		if err != ErrClosed {
			t.Error(fmt.Sprintf("expect ErrClosed got %v", err))
			return
		}

		_, err = s.Read()
		if err != ErrClosed {
			t.Error(fmt.Sprintf("expect ErrClosed got %v", err))
			return
		}

		err = s.Write([]byte("ErrClosed"))
		if err != ErrClosed {
			t.Error(fmt.Sprintf("expect ErrClosed got %v", err))
			return
		}
	}, func(c Conn) {
		c.Close()
	})
}
