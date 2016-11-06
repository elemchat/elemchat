package wizard

import (
	"time"

	"github.com/elemchat/elemchat/wizard/attr"
)

type Config struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DefaultAttr  attr.Attr
}

func DefaultConfig() *Config {
	return &Config{
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
		DefaultAttr: attr.Attr{
			Blood: 20,
		},
	}
}
