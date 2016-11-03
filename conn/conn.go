package conn

import (
	"errors"
	"time"
)

var (
	ErrReadTimeout  = errors.New("read timeout")
	ErrWriteTimeout = errors.New("write timeout")
	ErrClosed       = errors.New("closed")
)

type Conn interface {
	Read() (msg []byte, err error)
	Write(msg []byte) error
	Close()
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}
