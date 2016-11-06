package attr

import "testing"

func TestAttr_Diff(t *testing.T) {
	before := Attr{Blood: 10}
	after := Attr{Blood: 8}
	diff := after.Diff(before)
	if len(diff) != 1 {
		t.Error("expect", 1, "got", len(diff))
	}
	if diff[BLOOD] != -2 {
		t.Error("expect", -2, "got", diff[BLOOD])
	}
}

func TestAttr_Effect(t *testing.T) {
	diff := make(Diff)
	diff[BLOOD] = -2
	attr := Attr{Blood: 10}
	(&attr).Effect(diff)
	if attr.Blood != 8 {
		t.Error("expect", 8, "got", attr.Blood)
	}
}
