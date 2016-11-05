package field

import (
	"testing"
	"time"

	"github.com/elemchat/elemchat/codec"
	"github.com/elemchat/elemchat/conn"
	"github.com/elemchat/elemchat/msg"
	"github.com/elemchat/elemchat/wizard"
)

func TestNew_Enter(t *testing.T) {
	field := New()
	s, c := conn.TestPair()
	field.Enter(func(recv chan<- wizard.Message) *wizard.Wizard {
		return wizard.New("mofon",
			wizard.Attr{Blood: 10}, s, codec.JsonCodec(), recv)

	})
	s2, c2 := conn.TestPair()
	field.Enter(func(recv chan<- wizard.Message) *wizard.Wizard {
		return wizard.New("wizard",
			wizard.Attr{Blood: 10}, s2, codec.JsonCodec(), recv)

	})
	send(c, &msg.Chat{Text: "hello field!"})
	send(c2, &msg.Chat{Text: "hello field!"})
	field.Close()
	time.Sleep(1 * time.Second)

}

func send(conn conn.Conn, m msg.Message) {
	data, err := codec.JsonCodec().Encode(m)
	if err != nil {
		panic(err)
	}
	err = conn.Write(data)
	if err != nil {
		panic(err)
	}
}
