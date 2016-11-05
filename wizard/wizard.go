package wizard

import (
	"fmt"
	"time"

	"github.com/elemchat/elemchat/codec"
	"github.com/elemchat/elemchat/conn"
	"github.com/elemchat/elemchat/msg"
)

type AttrValue int
type Attr struct {
	Blood AttrValue
}

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
	conn  conn.Conn
	codec codec.Codec
	recv  chan<- Message

	Name string

	Attr Attr
}

// Wizard do NOT close chan recv
func New(name string, attr Attr,
	conn conn.Conn, codec codec.Codec, recv chan<- Message) *Wizard {
	w := &Wizard{
		conn:  conn,
		codec: codec,
		recv:  recv,
		Name:  name,
		Attr:  attr,
	}

	go w.recvLoop()
	return w
}

func (w *Wizard) Closed() bool {
	if w.conn == nil {
		return true
	}
	return false
}

func (w *Wizard) Close() {
	w.close()
}

func (w *Wizard) close() {
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
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
	err = w.conn.Write(data)
	if err != nil {
		// handle err
		if err == conn.ErrClosed {
			w.close()
		}

		return false
	}
	return true
}

func (w *Wizard) recvLoop() {
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
		data, err := w.conn.Read()
		if err != nil {
			// handle err
			if err == conn.ErrClosed {
				w.close()
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
