package main

import (
	"io"
)

const (
	db6BitMask  = 0x0
	db14bitMask = 0x40

	lead6Sz  = 1
	lead14Sz = 2
)

type byteWriter interface {
	io.ByteWriter
	io.Writer
}

type writer struct {
	i    int
	buff []byte
}

type Encoder struct {
	w byteWriter
}

type encoder struct {
	Encoder Encoder
}

func (w *writer) WriteByte(b byte) error {
	if w.i+1 > len(w.buff) {
		return io.ErrShortWrite
	}
	w.buff[w.i] = b
	w.i++
	return nil
}

func (w *writer) Write(b []byte) (n int, e error) {
	if w.i+len(b) > len(w.buff) {
		n = len(b) - w.i
		n = copy(w.buff[w.i:], b[:n])
		w.i += n
		return w.i, io.ErrShortWrite
	}
	n = copy(w.buff[w.i:], b)
	w.i += n
	return w.i, nil
}

func (en *Encoder) Encode6Bit(n int) (int, error) {
	err := en.w.WriteByte(byte(Bit6Mask | n))
	if err != nil {
		return 0, err

	}
	return 1, nil
}

func (en *Encoder) Encode14BitLen(n int) (int, error) {
	n, err := en.w.Write([]byte{byte(Bit14Mask | (n >> 8)), byte(n)})
	return n, err
}

func encode6BitLen(n int) []byte {
	return []byte{byte(db6BitMask | n)}
}

func encode14BitLen(n int) []byte {
	return []byte{
		byte(db14bitMask | (n >> 8)),
		byte(n),
	}
}

func encode32BitLen(n int) []byte {
	return []byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
}

func (en *encoder) Encode(t int, b []byte) (d []byte, e error) {
	switch t {
	case TypeString:
		d, e := en.Encoder.encodeString(b)
		if e != nil {
			return d, e
		}
		return d, nil
	}
	return d, nil
}

func (en *Encoder) stringWriter(sz int, p []byte) (b []byte, e error) {
	b = make([]byte, len(p)+sz)
	en = &Encoder{&writer{i: 0, buff: b}}
	n := 0
	switch sz {
	case lead6Sz:
		n, e = en.Encode6Bit(len(p))

	case lead14Sz:
		n, e = en.Encode14BitLen(len(p))
	}
	if e != nil {
		return b[:n], e
	}
	n, e = en.w.Write(p)
	if e != nil {
		return b[:n], e
	}
	return b[:n], nil
}
func (en *Encoder) encodeString(p []byte) (b []byte, e error) {
	if len(p) <= bit6Max/8 {
		b, _ = en.stringWriter(lead6Sz, p)
	} else if len(p) <= bit14Max/8 {
		b, _ = en.stringWriter(lead14Sz, p)
	}
	return b, nil
}

func encodeString(b []byte) (d []byte) {
	buff := make([]byte, 2048)
	en := Encoder{&writer{i: 0, buff: buff}}
	n := 0
	if len(b) <= bit6Max/8 {
		en.Encode6Bit(len(b))
	} else if len(b) <= bit14Max/8 {
		en.Encode14BitLen(len(b))

	} else {
		// TODO add functionality for additional objects i.e integers
		return d
	}
	n, _ = en.w.Write(b)
	return buff[:n]
}
