package wizard

import (
	"reflect"
	"sync"
	"testing"

	"github.com/elemchat/elemchat/codec"
	"github.com/elemchat/elemchat/conn"
	"github.com/elemchat/elemchat/msg"
	"github.com/elemchat/elemchat/wizard/attr"
)

func TestNew_Send(t *testing.T) {
	s, c := conn.TestPair()
	defer s.Close()
	defer c.Close()

	recv := make(chan Message)
	w := New("wizard", attr.Attr{Blood: 10}, c, codec.JsonCodec(), recv)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		data, err := s.Read()
		if err != nil {
			t.Error(err)
			return
		}
		m, err := codec.JsonCodec().Decode(data)
		if err != nil {
			t.Error(err)
			return
		}

		if m == nil {
			t.Error("expect not nil got nil")
			return
		}
		if m, ok := m.(*msg.Chat); !ok {
			t.Error("expect *msg.Chat got", reflect.TypeOf(m))
			return
		} else {
			if m.Text != "hello wizard!" {
				t.Error("expect hello wizard! got", m.Text)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		w.Send(&msg.Chat{Text: "hello wizard!"})
	}()
	wg.Wait()
}

func TestNew_recvLoop(t *testing.T) {
	s, c := conn.TestPair()
	defer s.Close()
	defer c.Close()

	recv := make(chan Message)
	w := New("wizard", attr.Attr{Blood: 10}, c, codec.JsonCodec(), recv)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// recv
	go func() {
		defer wg.Done()

		m := <-recv
		if m == nil {
			t.Error("expect not nil got nil")
			return
		}
		if m.Wizard() != w {
			t.Error("expect", w, "got", m.Wizard())
			return
		}

		if m.Msg() == nil {
			t.Error("expect not nil got nil")
			return
		}

		if m, ok := m.Msg().(*msg.Chat); !ok {
			t.Error("expect *msg.Chat got", reflect.TypeOf(m))
			return
		} else {
			if m.Text != "hello recvLoop!" {
				t.Error("expect hello recvLoop! got", m.Text)
				return
			}
		}
	}()
	// send
	go func() {
		defer wg.Done()

		data, err := codec.JsonCodec().Encode(&msg.Chat{Text: "hello recvLoop!"})
		if err != nil {
			t.Error(err)
			return
		}

		err = s.Write(data)
		if err != nil {
			t.Error(err)
			return
		}

	}()
	wg.Wait()
	close(recv)
}
