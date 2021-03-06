package magic

import (
	"encoding/json"
	"fmt"
	"testing"
)

type case_magic_string struct {
	Magic  Magic
	String string
}

var cases = []case_magic_string{
	{ILLEGAL, "illegal"},
	{FIRE_BALL, "fireBall"},
	{WATER_WALL, "waterWall"},
	{GRASS_GROWTH, "grassGrowth"},
}

func TestMagic_String(t *testing.T) {
	for _, c := range cases {
		if c.Magic.String() != c.String {
			t.Error("expect", c.String, "got", c.Magic.String())
		}
	}
}

func TestToMagic(t *testing.T) {
	for _, c := range cases {
		if ToMagic(c.String) != c.Magic {
			t.Error("expect", c.Magic, "got", ToMagic(c.String))
		}
	}
}

func TestMagic_MarshalJson(t *testing.T) {
	for _, c := range cases {
		data, err := json.Marshal(&c.Magic)
		if err != nil {
			t.Error(err)
			continue
		}

		if string(data) != fmt.Sprintf("%q", c.String) {
			t.Error("expect", fmt.Sprintf("%q", c.String), "got", string(data))
		}
	}
}

func TestMagic_UnmarshalJson(t *testing.T) {
	for _, c := range cases {
		var m Magic
		err := json.Unmarshal([]byte(fmt.Sprintf("%q", c.String)), &m)
		if err != nil {
			t.Error(err)
			continue
		}

		if m != c.Magic {
			t.Error("expect", c.Magic, "got", m)
		}
	}
}
