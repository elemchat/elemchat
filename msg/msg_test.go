package msg

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_TimeJson(t *testing.T) {
	tm := time.Unix(0, 0).Add(1478432655765 * time.Millisecond)
	wm := (&WaitMagic{}).SetDeadline(tm)
	data, err := json.Marshal(wm)
	if err != nil {
		t.Error(err)
		return
	}
	if string(data) != `{"deadline":1478432655765}` {
		t.Error("expect", `{"deadline":1478432655765}`, "got", string(data))
		return
	}

	unwm := &WaitMagic{}
	err = json.Unmarshal(data, unwm)
	if err != nil {
		t.Error(err)
		return
	}
	if unwm.GetDeadline().String() !=
		"2016-11-06 19:44:15.765 +0800 CST" {
		t.Error("expect", "2016-11-06 19:44:15.765 +0800 CST",
			"got", unwm.GetDeadline())
	}

}
