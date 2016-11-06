package msg

import "errors"

var (
	ErrIllegal = errors.New("illegal type")
)

type Type int

const (
	ILLEGAL Type = iota
	AUTH
	TIME_SYNC_REQ
	TIME_SYNC_RESP
	INIT
	WAIT_CHAT
	WAIT_MAGIC
	CHAT
	MAGIC
	EFFECT
	DUALOVER

	// length of Type
	typeLength
)

func (t Type) String() string {
	switch t {
	case ILLEGAL:
		return "illegal"
	case AUTH:
		return "auth"
	case TIME_SYNC_REQ:
		return "timeSyncReq"
	case TIME_SYNC_RESP:
		return "timeSyncResp"
	case INIT:
		return "init"
	case WAIT_CHAT:
		return "waitChat"
	case WAIT_MAGIC:
		return "waitMagic"
	case CHAT:
		return "chat"
	case MAGIC:
		return "magic"
	case EFFECT:
		return "effect"
	case DUALOVER:
		return "dualover"
	default:
		return "illegal"
	}
}

func ToType(t string) Type {
	switch t {
	case "illegal":
		return ILLEGAL
	case "auth":
		return AUTH
	case "timeSyncReq":
		return TIME_SYNC_REQ
	case "timeSyncResp":
		return TIME_SYNC_RESP
	case "init":
		return INIT
	case "waitChat":
		return WAIT_CHAT
	case "waitMagic":
		return WAIT_MAGIC
	case "chat":
		return CHAT
	case "magic":
		return MAGIC
	case "effect":
		return EFFECT
	case "dualover":
		return DUALOVER
	}
	return ILLEGAL
}

type DecodeFunc func(Message) error

func Decode(t Type, fn DecodeFunc) (Message, error) {
	unmarshalMsg := func(m Message) (Message, error) {
		err := fn(m)
		if err != nil {
			return nil, err
		}
		return m, nil
	}
	switch t {
	case ILLEGAL:
		return nil, ErrIllegal
	case AUTH:
		return unmarshalMsg(&Auth{})
	case TIME_SYNC_REQ:
		return unmarshalMsg(&TimeSyncReq{})
	case TIME_SYNC_RESP:
		return unmarshalMsg(&TimeSyncResp{})
	case INIT:
		return unmarshalMsg(&Init{})
	case WAIT_CHAT:
		return unmarshalMsg(&WaitChat{})
	case WAIT_MAGIC:
		return unmarshalMsg(&WaitMagic{})
	case CHAT:
		return unmarshalMsg(&Chat{})
	case MAGIC:
		return unmarshalMsg(&Magic{})
	case EFFECT:
		return unmarshalMsg(&Effect{})
	case DUALOVER:
		return unmarshalMsg(&Dualover{})
	}
	return nil, ErrIllegal
}
