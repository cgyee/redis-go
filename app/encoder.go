package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
)

const (
	encInt8  = 0
	encInt16 = 1
	encInt32 = 2
	encObj   = 192
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	en := &Encoder{w}
	return en
}

func (en *Encoder) EncodeLength(l int) (n int) {
	if l <= bitMax6/8 {
		n, _ = en.w.Write([]byte{byte(l)})
	} else if l <= bitMax14/8 {
		n, _ = en.w.Write([]byte{byte(l>>8) | bitLen14, byte(l)})
	} else {
		b := make([]byte, 5)
		b[0] = bitLen32
		binary.BigEndian.PutUint32(b, uint32(l))
		n, _ = en.w.Write(b)
	}
	return n
}

func (en *Encoder) EncodeString(p []byte) (n int, e error) {
	l := en.EncodeLength(len(p))
	n, e = en.w.Write(p)
	if e != nil {
		fmt.Println("error EncodeString: ", e)
	}
	return n + l, nil
}

func (en *Encoder) EncodeInt(p []byte) (n int, err error) {
	s := string(p)
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return
	}

	switch {
	case math.MinInt8 <= v && v <= math.MaxInt8:
		n, err = en.w.Write([]byte{encObj, byte(int8(v))})
	case math.MinInt16 <= v && v <= math.MaxInt16:
		intBuf := make([]byte, 3)
		intBuf[0] = encInt8
		binary.BigEndian.PutUint16(intBuf[1:], uint16(n))
		n, err = en.w.Write(intBuf)
	case math.MinInt32 <= v && v <= math.MaxInt32:
		b := make([]byte, 5)
		b[0] = encObj | 2
		binary.LittleEndian.PutUint32(b, uint32(v))
		n, err = en.w.Write(b)
	}
	return n, err
}

func (en *Encoder) EncodeTimeStamp(ms bool, n int) (b []byte, e error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, int64(n))
	if err != nil {
		fmt.Println("binary.Write failed")
	}

	en.w.Write(buff.Bytes())
	if ms {
		en.w.Write([]byte{0xFC})
	} else {
		en.w.Write([]byte{0xFD})
	}
	i, e := en.w.Write(b)
	if e != nil {
		fmt.Printf("encode time failed:\n%v\n", e)
	}

	return b[:i], nil

}
