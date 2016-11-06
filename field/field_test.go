package field

import (
	"reflect"
	"testing"

	"github.com/elemchat/elemchat/codec"
	"github.com/elemchat/elemchat/conn"
	"github.com/elemchat/elemchat/msg"
	"github.com/elemchat/elemchat/wizard"
	"github.com/elemchat/elemchat/wizard/attr"
)

func TestNew_Enter(t *testing.T) {
	msgch := make(chan wizard.Message)
	s, c := conn.TestPair()
	field := New(func(message wizard.Message) {
		msgch <- message
	})

	field.Enter(func(recv chan<- wizard.Message) *wizard.Wizard {
		return wizard.New("wizard",
			attr.Attr{Blood: 10}, s, codec.JsonCodec(), recv)

	})

	send(c, &msg.Chat{Text: "hello field!"})
	message := <-msgch
	if message.Wizard().Name != "wizard" {
		t.Error("expect wizard got", message.Wizard().Name)
		return
	}
	if msg.GetType(message.Msg()) != msg.CHAT {
		t.Error("expect", msg.CHAT, "got", msg.GetType(message.Msg()))
		return
	}
	if m, ok := message.Msg().(*msg.Chat); !ok {
		t.Error("expect *msg.Chat got", reflect.TypeOf(message.Msg()))
		return
	} else {
		if m.Text != "hello field!" {
			t.Error("expect hello field! got", m.Text)
			return
		}
	}
	field.Close()
	field.WithLock(func(f *Field) {
		if len(f.Wizards) != 0 {
			t.Error("expect", 0, "got", len(f.Wizards))
		}
		for w, _ := range f.Wizards {
			if w.Closed() {
				t.Error("expect", w.Name, "is Closed", "got", "not")
			}
		}
	})

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
