package msg

//go:generate stringer -type=Type type.go
type Type int

const (
	ILLEGAL Type = iota
	INIT
	WAIT_CHAT
	WAIT_MAGIC
	CHAT
	MAGIC
	EFFECT
	DUALOVER
)
