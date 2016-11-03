package msg

import (
	"time"

	"github.com/elemchat/elemchat/magic"
)

var (
	_ Message = &Init{}
	_ Message = &WaitChat{}
	_ Message = &WaitMagic{}
	_ Message = &Chat{}
	_ Message = &Magic{}
	_ Message = &Effect{}
	_ Message = &Dualover{}
)

type Message interface {
	_type() Type
}

func GetType(msg Message) Type {
	return msg._type()
}

type Init struct {
}

type WaitChat struct {
	Deadline time.Time `json:"deadline"`
}

type WaitMagic struct {
	Deadline time.Time `json:"deadline"`
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
