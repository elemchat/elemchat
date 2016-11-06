package field

import (
	"sync"

	"github.com/elemchat/elemchat/wizard"
)

type Field struct {
	sync.Mutex
	Wizards map[*wizard.Wizard]struct{}
	recv    chan wizard.Message
	closed  bool
}

type HandleFunc func(wizard.Message)

func New(handle HandleFunc) *Field {
	f := &Field{
		Mutex:   sync.Mutex{},
		Wizards: make(map[*wizard.Wizard]struct{}),
		recv:    make(chan wizard.Message),
		closed:  false,
	}

	go f.loop(handle)
	return f
}

// WithLock is MUST used on anything of Field
func (f *Field) WithLock(fn func(*Field)) {
	f.Lock()
	fn(f)
	f.Unlock()
}

func (f *Field) Close() {
	f.WithLock(func(f *Field) {
		if f.closed {
			return
		}

		go func() {
			for _ = range f.recv {
				//do nothing
			}
		}()

		for w, _ := range f.Wizards {
			w.Close(true)
			delete(f.Wizards, w)
		}
		close(f.recv)
		f.closed = true
	})
}

func (f *Field) Enter(fn func(recv chan<- wizard.Message) *wizard.Wizard) {
	f.WithLock(func(f *Field) {
		if f.closed {
			return
		}

		w := fn(f.recv)
		if w != nil {
			f.Wizards[w] = struct{}{}
		}
	})
}

func (f *Field) loop(handle HandleFunc) {
	if handle == nil {
		f.Close()
		return
	}
	for message := range f.recv {
		if message.Msg() == nil {
			continue
		}

		handle(message)
	}
}
