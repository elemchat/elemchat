package wizard

import (
	"fmt"
	"time"

	"github.com/elemchat/elemchat/codec"
	"github.com/elemchat/elemchat/conn"
	"github.com/elemchat/elemchat/msg"
	"github.com/elemchat/elemchat/wizard/attr"
)

type Message interface {
	Wizard() *Wizard
	Time() time.Time
	Msg() msg.Message
}
type message struct {
	wizard *Wizard
	time   time.Time
	msg    msg.Message
}

func (m *message) Wizard() *Wizard {
	return m.wizard
}
func (m *message) Time() time.Time {
	return m.time
}
func (m *message) Msg() msg.Message {
	return m.msg
}

type Wizard struct {
	conn   conn.Conn
	codec  codec.Codec
	config *Config // read-only
	recv   chan<- Message
	closed chan struct{}

	Name string

	Attr attr.Attr
}

// Wizard do NOT close chan recv
func New(name string, attr attr.Attr,
	conn conn.Conn, codec codec.Codec, recv chan<- Message) *Wizard {
	w := &Wizard{
		conn:   conn,
		codec:  codec,
		config: DefaultConfig(),
		recv:   recv,
		closed: make(chan struct{}),
		Name:   name,
		Attr:   attr,
	}

	go w.recvLoop()
	return w
}

func NewWithConfig(config *Config,
	name string,
	conn conn.Conn, codec codec.Codec,
	recv chan<- Message) *Wizard {
	w := &Wizard{
		conn:   conn,
		codec:  codec,
		config: config,
		recv:   recv,
		closed: make(chan struct{}),
		Name:   name,
		Attr:   config.DefaultAttr,
	}

	go w.recvLoop()
	return w
}

func (w *Wizard) Closed() bool {
	select {
	case <-w.closed:
		return true
	default:
		return false
	}
}

// NOTE:
// While call .Close(),caller show receive all message in .recv chan
// before close recv chan.
// If do not,DEADLOCK maybe occur!
func (w *Wizard) Close(wait bool) {
	w.close()
	if wait {
		<-w.closed
	}
}

func (w *Wizard) close() {
	if !w.Closed() {
		w.conn.Close()
	}
}

func (w *Wizard) Send(msg msg.Message) bool {
	if w.Closed() {
		return false
	}
	// encode
	if w.codec == nil {
		w.close()
		return false
	}
	data, err := w.codec.Encode(msg)
	if err != nil {
		return false
	}

	// write
	w.conn.SetWriteDeadline(time.Now().Add(w.config.WriteTimeout))
	err = w.conn.Write(data)
	if err != nil {
		// handle err
		if err == conn.ErrClosed {
			w.close()
		}
		if err == conn.ErrWriteTimeout {
			w.close()
		}

		return false
	}
	return true
}

func (w *Wizard) recvLoop() {
	defer close(w.closed)
	defer func() {
		if r := recover(); r != nil {
			if fmt.Sprint(r) == "send on closed channel" {
				w.close()
			} else {
				panic(r)
			}
		}
	}()
	for {

		// read
		if w.conn == nil {
			return
		}
		w.conn.SetReadDeadline(time.Now().Add(w.config.ReadTimeout))
		data, err := w.conn.Read()
		if err != nil {
			// handle err
			if err == conn.ErrClosed {
				w.close()
				return
			}
			if err == conn.ErrReadTimeout {
				w.close()
				return
			}
			continue
		}

		// decode
		if w.codec == nil {
			w.close()
			return
		}
		msg, err := w.codec.Decode(data)
		if err != nil {
			// handle err
			continue
		}

		// send to recv chan
		w.recv <- &message{
			wizard: w,
			time:   time.Now(),
			msg:    msg,
		}
	}
}
