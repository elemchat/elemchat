package codec

import (
	"fmt"
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
	if string(json) != `{"msg":{"Text":"hello codec!"},"type":"CHAT"}` {
		t.Error(fmt.Sprintf("expect "+
			`{"msg":{"Text":"hello codec!"},"type":"CHAT"}`+
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
