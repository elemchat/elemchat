package msg

import (
	"strconv"
	"time"

	"github.com/elemchat/elemchat/magic"
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

type Auth struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type Init struct {
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
}

type Dualover struct {
}

func (_ *Auth) _type() Type {
	return AUTH
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
