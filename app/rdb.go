package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	TypeString = 0
)

const (
	bit6Max        = (1 << 6) - 1
	bit14Max       = (1 << 14) - 1
	bit32Max       = (1 << 32) - 1
	leadBitsMask   = 0x3F
	Bit6Mask       = 0x0
	Bit14Mask      = 0x40
	BitObjectMask  = 0xC0
	Bit6Length     = 0
	Bit14Length    = 1
	BitObject      = 2
	BitLengthInt8  = 10
	BitLengthInt16 = 11
	BitLengthInt32 = 15
)

const ()

type RData struct {
	key       []byte
	value     []byte
	expires   bool
	expireMS  bool
	expiresIn time.Time
	dataType  rune
}

type DbWriter struct {
	index int
}

var w *DbWriter

func InitFile() {
	file, err := os.Open("tempfile")

	if err != nil {
		file, err := os.Create("tempfile")
		if err != nil {
			fmt.Println(err)
			header := make([]byte, 0)
			header = append(header, []byte("REDIS")...)
			header = append(header, []byte("0011")...)
			header = append(header, 0xFA)
			header = append(header, 0xFE)
			header = append(header, 0x00)
			_, err = file.Write(header)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			n, _ := ReadAll(file)
			w.index = n

		}

	}
	defer file.Close()
	file.Close()

}

func ReadAll(f *os.File) (n int, b []byte) {
	for {
		buff := make([]byte, 4096)
		n, err := f.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return 0, make([]byte, 0)
		}
		if n < len(buff) {
			b = buff[n:]
			break
		} else {
			b = append(b, buff...)
		}
	}
	return len(b), b

}

func (DbWriter) Write(rd RData) {
	// idx := DbWriter.index
	// Todo move file opening to after transforming data
	_, err := os.Open("tempfile")
	if err != nil {
		fmt.Println(err)
	}
	d := make([]byte, 0)
	if rd.expires {
		if rd.expireMS {
			time := rd.expiresIn.UnixMilli()
			bitTime := []byte{
				byte(time),
				byte(time >> 8),
				byte(time >> 16),
				byte(time >> 24),
				byte(time >> 32),
				byte(time >> 40),
				byte(time >> 48),
				byte(time >> 56),
			}
			d = append(d, 0xFC)
			d = append(d, bitTime...)

		} else {
			time := rd.expiresIn.Unix()
			bitTime := []byte{
				byte(time),
				byte(time >> 8),
				byte(time >> 16),
				byte(time >> 24),
			}

			d = append(d, 0xFD)
			d = append(d, bitTime...)
		}
	}
	switch rd.dataType {
	case STRING:
		d = append(d, byte(0))

	}
}

// func encodeString(b []byte) (d []byte) {
// 	if len(b) <= bit6Max/8 {
// 		bits := Bit6Mask | len(b)
// 		d = append(d, byte(bits))
// 	} else if len(b) <= bit14Max/8 {
// 		bits := Bit14Mask | (len(b) >> 8)
// 		d = append(d, byte(bits), byte(len(b)))
// 	} else {
// 		// TODO add functionality for additional objects i.e integers
// 		return d
// 	}
// 	d = append(d, b...)
// 	return d

// }

func decodeString(b []byte) (d []byte) {
	strBitLen := (b[0] & 0xC0) >> 6
	l := 0
	if strBitLen == Bit6Length {
		b[0] = CleanLeadBits(b[0])
		l = int(b[0])
		b = b[1:]
	} else if strBitLen == Bit14Length {
		b[0] = CleanLeadBits(b[0])
		l = int(0x3F&b[0])<<8 + int(b[1])
		b = b[2:]
	} else {
		for i := 1; i < 5; i++ {
			l = l<<8 + int(b[i])
		}
		b = b[5:]

	}
	return b[:l]

}

func UintToInt(b []byte) int {
	switch len(b) {
	// We only need the last 6 bits to get our number
	case 1:
		return int(b[0] & 0b111111)
	case 2:
		fmt.Println("16 bits!!")
		return int(((b[0] & 0b111111) << 8) + b[1])
	case 4:
		return int((binary.BigEndian.Uint32([]byte{b[0] & 0x3F, b[1], b[2], b[3]})))
	}
	return -1
}

func CleanLeadBits(b byte) byte {
	return b & 0x3F
}
