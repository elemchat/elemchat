package msg

import (
	"strconv"
	"time"

	"github.com/elemchat/elemchat/magic"
	"github.com/elemchat/elemchat/wizard/attr"
)

var (
	_ Message = &Auth{}
	_ Message = &Init{}
	_ Message = &WaitChat{}
	_ Message = &WaitMagic{}
	_ Message = &Chat{}
	_ Message = &Magic{}
	_ Message = &Effect{}
	_ Message = &Dualover{}
)

type time_Time time.Time

func (t time_Time) MarshalJSON() ([]byte, error) {
	dur := time.Duration(time.Time(t).UnixNano())
	return []byte(strconv.Itoa(int(dur / time.Millisecond))), nil
}

func (t *time_Time) UnmarshalJSON(date []byte) error {
	unixms, err := strconv.Atoi(string(date))
	if err != nil {
		return err
	}

	dur := time.Duration(unixms) * time.Millisecond
	*t = time_Time(time.Unix(0, 0).Add(dur))
	return nil
}

type Message interface {
	_type() Type
}

func GetType(msg Message) Type {
	return msg._type()
}

type Ping struct{}
type Pong struct{}

type Match struct {
	FieldURI string `json"fieldURI"`
	Password string `json"password"`
}

type Auth struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type TimeSyncReq struct {
	T0 time_Time `json:"t0"`
}

func (tsr TimeSyncReq) GetT0() time.Time {
	return time.Time(tsr.T0)
}

func (tsr *TimeSyncReq) SetT0(t time.Time) *TimeSyncReq {
	tsr.T0 = time_Time(t)
	return tsr
}

type TimeSyncResp struct {
	T0 time_Time `json:"t0"`
	T1 time_Time `json:"t1"`
	T2 time_Time `json:"t2"`
}

func (tsr TimeSyncResp) GetT0() time.Time {
	return time.Time(tsr.T0)
}

func (tsr *TimeSyncResp) SetT0(t time.Time) *TimeSyncResp {
	tsr.T0 = time_Time(t)
	return tsr
}

func (tsr TimeSyncResp) GetT1() time.Time {
	return time.Time(tsr.T1)
}

func (tsr *TimeSyncResp) SetT1(t time.Time) *TimeSyncResp {
	tsr.T1 = time_Time(t)
	return tsr
}

func (tsr TimeSyncResp) GetT2() time.Time {
	return time.Time(tsr.T2)
}

func (tsr *TimeSyncResp) SetT2(t time.Time) *TimeSyncResp {
	tsr.T2 = time_Time(t)
	return tsr
}

type Init struct {
	Attrs map[string]attr.Attr `json:"attrs"`
}

type WaitChat struct {
	Deadline time_Time `json:"deadline"`
}

func (wc WaitChat) GetDeadline() time.Time {
	return time.Time(wc.Deadline)
}
func (wc *WaitChat) SetDeadline(t time.Time) *WaitChat {
	wc.Deadline = time_Time(t)
	return wc
}

type WaitMagic struct {
	Deadline time_Time `json:"deadline"`
}

func (wm WaitMagic) GetDeadline() time.Time {
	return time.Time(wm.Deadline)
}
func (wm *WaitMagic) SetDeadline(t time.Time) *WaitMagic {
	wm.Deadline = time_Time(t)
	return wm
}

type Chat struct {
	Text string `json:"text"`
}

type Magic struct {
	Magic magic.Magic `json:"magic"`
}
type Effect struct {
	Diff attr.Diff `json:"diff"`
}

type Dualover struct {
}

func (_ *Ping) _type() Type {
	return PING
}
func (_ *Pong) _type() Type {
	return PONG
}
func (_ *Match) _type() Type {
	return MATCH
}
func (_ *Auth) _type() Type {
	return AUTH
}
func (_ *TimeSyncReq) _type() Type {
	return TIME_SYNC_REQ
}
func (_ *TimeSyncResp) _type() Type {
	return TIME_SYNC_RESP
}
func (_ *Init) _type() Type {
	return INIT
}
func (_ *WaitChat) _type() Type {
	return WAIT_CHAT
}
func (_ *WaitMagic) _type() Type {
	return WAIT_MAGIC
}
func (_ *Chat) _type() Type {
	return CHAT
}
func (_ *Magic) _type() Type {
	return MAGIC
}
func (_ *Effect) _type() Type {
	return EFFECT
}
func (_ *Dualover) _type() Type {
	return DUALOVER
}
