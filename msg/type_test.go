package msg

import (
	"reflect"
	"testing"
)

type case_type_string struct {
	Type   Type
	String string
	Msg    Message
}

var cases = []case_type_string{
	{ILLEGAL, "illegal", nil},
	{INIT, "init", &Init{}},
	{WAIT_CHAT, "waitChat", &WaitChat{}},
	{WAIT_MAGIC, "waitMagic", &WaitMagic{}},
	{CHAT, "chat", &Chat{}},
	{MAGIC, "magic", &Magic{}},
	{EFFECT, "effect", &Effect{}},
	{DUALOVER, "dualover", &Dualover{}},
}

func TestCasesLength(t *testing.T) {
	if len(cases) != int(typeLength) {
		t.Error("cases length is not equal type's lenght")
		return
	}
}

func TestType_String(t *testing.T) {
	for _, c := range cases {
		if c.Type.String() != c.String {
			t.Error("expect", c.String, "got", c.Type.String())
		}
	}
}

func TestType_ToType(t *testing.T) {
	for _, c := range cases {
		if ToType(c.String) != c.Type {
			t.Error("expect", c.Type, "got", ToType(c.String))
		}
	}

}
func TestType_Decode(t *testing.T) {
	for _, c := range cases {
		Decode(c.Type, func(m Message) error {
			if c.Type == ILLEGAL {
				return nil
			}
			if reflect.TypeOf(m) != reflect.TypeOf(c.Msg) {
				t.Error("expect", reflect.TypeOf(c.Msg), "got", reflect.TypeOf(m))
			}
			return nil
		})
	}
}
