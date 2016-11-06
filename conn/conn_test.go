package conn

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type ConnPairFunc func() (Conn, Conn)

func TestConn(t *testing.T) {
	t.Run("TestPair", func(t *testing.T) {
		subtest_Conn(t, func() (Conn, Conn) {
			return TestPair()
		})
	})
	t.Run("WebSocket", func(t *testing.T) {
		subtest_Conn(t, func() (Conn, Conn) {
			connChan := make(chan *websocket.Conn)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn_s, err := websocket.Upgrade(w, r, nil, 0, 0)
				if err != nil {
					t.Error(err)
				}

				connChan <- conn_s
			}))
			server_url, err := url.Parse(server.URL)

			dialer := &websocket.Dialer{}
			if err != nil {
				t.Fatal(err)
			}

			conn_c, _, err := dialer.Dial(fmt.Sprintf("ws://%s/", server_url.Host), nil)
			if err != nil {
				t.Fatal(err)
				return nil, nil
			}

			select {
			case conn_s := <-connChan:
				if conn_s != nil && conn_c != nil {
					return WebSocket(conn_s), WebSocket(conn_c)
				}
			default:
			}
			t.Fatal("get conn server failure")
			return nil, nil
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
		defer s.Close()
		defer c.Close()
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
			t.Error(fmt.Sprintf("expect ErrReadTimeout got %v", err, reflect.TypeOf(err)))
			return
		}

	}, func(c Conn) {
	})

	wait(func(s Conn) {
		s.SetWriteDeadline(time.Now().Add(time.Second * 1))
		var err error
		for {
			err = s.Write([]byte("write timeout"))
			if err != nil {
				break
			}
		}
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
