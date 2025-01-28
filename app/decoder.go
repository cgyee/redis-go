package main

import (
	"io"
	"os"
)

type Decoder struct {
	r   io.Reader
	buf []byte
}

const (
	strMask  = 0xC0
	strLen6  = 0
	strLen14 = 1
	strLen32 = 3
)

func NewDecoder(r io.Reader) *Decoder {
	buf := make([]byte, 0)
	de := &Decoder{r, buf}
	return de
}

// Process -> Read a Byte, determine if it's a string or the the length < 4 bytes
// If yes return the length, if no return false
func (de *Decoder) decodeLength(f *os.File) (n int, isStr bool) {
	b := make([]byte, 1)
	_, err := f.Read(b)
	if err != nil {
		return
	}
	t := (b[0] & strMask) >> 6
	switch t {
	case strLen6:
		return int(b[0]), true
	case strLen14:
		return int(p[0] & 0x3F), true
	case strLen32:

		return

	}
	return

}

func (de *Decoder) Read() {
	de.r.Read()
}

func (de *Decoder) decodeString() {}

func (de *Decoder) Deocode() {}
