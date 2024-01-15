package bifrost

import (
	"bytes"

	"github.com/pkg/term"
)

// KeyType describes a key
type KeyType uint16

// Key holds the input value from terminal stream
type Key struct {
	Value []byte
	Type  KeyType
}

// Key constants
const (
	Enter KeyType = iota + 1
	Esc
	Space
	Tab
	CtrlA
	CtrlB
	CtrlC
	Delete
	CtrlBackslash
	Backspace
	UpArrow
	DownArrow
	LeftArrow
	RightArrow
)

func pollKeyEvents() Key {
	t, err := term.Open("/dev/tty")
	if err != nil {
		return Key{}
	}

	term.RawMode(t)
	buff := make([]byte, 2048)
	size, err := t.Read(buff)
	t.Restore()
	t.Close()
	if err != nil {
		return Key{}
	}

	switch {
	case bytes.Equal(buff[0:size], []byte{13}) || bytes.Equal(buff[0:size], []byte{10}):
		return Key{Type: Enter}
	case bytes.Equal(buff[0:size], []byte{27}):
		return Key{Type: Esc}
	case bytes.Equal(buff[0:size], []byte{32}):
		return Key{Type: Space}
	case bytes.Equal(buff[0:size], []byte{9}):
		return Key{Type: Tab}
	case bytes.Equal(buff[0:size], []byte{1}):
		return Key{Type: CtrlA}
	case bytes.Equal(buff[0:size], []byte{2}):
		return Key{Type: CtrlB}
	case bytes.Equal(buff[0:size], []byte{3}):
		return Key{Type: CtrlC}
	case bytes.Equal(buff[0:size], []byte{28}):
		return Key{Type: CtrlBackslash}
	case bytes.Equal(buff[0:size], []byte{8}) || bytes.Equal(buff[0:size], []byte{127}):
		return Key{Type: Backspace}
	case bytes.Equal(buff[0:size], []byte{27, 91, 51, 126}):
		return Key{Type: Delete}
	case bytes.Equal(buff[0:size], []byte{27, 91, 65}):
		return Key{Type: UpArrow}
	case bytes.Equal(buff[0:size], []byte{27, 91, 66}):
		return Key{Type: DownArrow}
	case bytes.Equal(buff[0:size], []byte{27, 91, 67}):
		return Key{Type: RightArrow}
	case bytes.Equal(buff[0:size], []byte{27, 91, 68}):
		return Key{Type: LeftArrow}
	default:
		return Key{Value: buff[0:size]}
	}
}
