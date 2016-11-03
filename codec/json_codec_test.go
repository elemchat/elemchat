package codec

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/elemchat/elemchat/msg"
)

func TestJsonCodec_Encode(t *testing.T) {
	codec := JsonCodec()
	json, err := codec.Encode(&msg.Chat{Text: "hello codec!"})
	if err != nil {
		t.Error(err)
		return
	}
	if string(json) != `{"msg":{"text":"hello codec!"},"type":"chat"}` {
		t.Error(fmt.Sprintf("expect "+
			`{"msg":{"text":"hello codec!"},"type":"chat"},"type":"CHAT"}`+
			" got %s",
			string(json)))
		return
	}

	_, err = codec.Encode(nil)
	if err != ErrMessageNil {
		t.Error("expect ErrMessageNil got ", err)
		return
	}
}

func TestJsonCodec_Decode(t *testing.T) {
	codec := JsonCodec()
	m, err := codec.Decode(
		[]byte(`{"msg":{"text":"hello codec!"},"type":"chat"}`))
	if err != nil {
		t.Error(err)
		return
	}
	if m == nil {
		t.Error("got nil expect not")
		return
	}

	if msg.GetType(m) != msg.CHAT {
		t.Error("expect", msg.CHAT, "got", msg.GetType(m))
		return
	}
	switch m := m.(type) {
	case *msg.Chat:
		if m.Text != "hello codec!" {
			t.Error("expect hello codec! got", m.Text)
			return
		}
	default:
		t.Error("expect m.(type) is *msg.Chat;got", reflect.TypeOf(m))
		return
	}

}
