package conn

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type websocketConn struct {
	ws        *websocket.Conn
	closeSign chan struct{}
}

func WebSocket(ws *websocket.Conn) Conn {
	conn := &websocketConn{
		ws:        ws,
		closeSign: make(chan struct{}),
	}
	if ws == nil {
		conn.Close()
	}
	return conn
}

func (conn *websocketConn) Read() (msg []byte, err error) {
	_, msg, err = conn.ws.ReadMessage()
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "i/o timeout") {
			return nil, ErrReadTimeout
		} else {
			conn.Close()
			return nil, ErrClosed
		}
	}
	return msg, err
}

func (conn *websocketConn) Write(msg []byte) error {
	err := conn.ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "i/o timeout") {
			return ErrWriteTimeout
		} else {
			conn.Close()
			return ErrClosed
		}
	}
	return err
}

func (conn *websocketConn) Close() {
	if !conn.closed() {
		close(conn.closeSign)
		conn.ws.Close()
	}
}

func (conn *websocketConn) SetReadDeadline(t time.Time) error {
	if conn.closed() {
		return ErrClosed
	}

	return conn.ws.SetReadDeadline(t)
}

func (conn *websocketConn) SetWriteDeadline(t time.Time) error {
	if conn.closed() {
		return ErrClosed
	}

	return conn.ws.SetWriteDeadline(t)
}

func (conn *websocketConn) closed() bool {
	select {
	case <-conn.closeSign:
		return true
	default:
		return false
	}
}
