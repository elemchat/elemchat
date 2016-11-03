package magic

import "testing"

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
