package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

const (
	TypeString = 0
)

const (
	bitLen6        = 0x0
	bitLen14       = 0x40
	bitLen32       = 0xC0
	bitLenObj      = 0xC0
	bitMax6        = (1 << 6) - 1
	bitMax14       = (1 << 14) - 1
	bitMax32       = (1 << 32) - 1
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

type RData struct {
	key       []byte
	value     []byte
	expires   bool
	expireMS  bool
	expiresIn time.Time
	dataType  rune
}

type DBWriter struct {
	index     int
	createdAt int64
	mem       int
	en        encoder
}

type DBReader struct {
	idx int
}

type db struct {
	f        *os.File
	fileName string
	en       encoder
}

// func (rdb *db) Write(t rune, p []byte) (n int, err error) {
// 	// en := &rdb.en
// 	switch t {
// 	case BULK:
// 	}

// }

func (db *DBWriter) writeHeader(f *os.File) (n int, err error) {
	header := make([]byte, 0)
	header = append(header, []byte("REDIS")...)

	n, err = f.Write(header)
	if err != nil {
		fmt.Printf("error during writeHeader:\n%v\n", err)
		return n, err
	}
	return n, nil

	//  metadata section
}

func (rdb *db) openFile() (*os.File, error) {
	f, err := os.OpenFile(rdb.fileName, os.O_RDWR, os.FileMode(os.O_CREATE))
	if err != nil {
		switch err {
		case os.ErrExist:
			f, err = os.OpenFile(rdb.fileName, os.O_RDWR, os.ModeAppend)
		}
	}
	f.Chmod(os.ModeAppend)
	return f, nil
}

func (rdb *db) overwriteFile() (*os.File, error) {
	f, err := os.Open(rdb.fileName)
	if err != nil {
		if err == os.ErrNotExist {
			f, _ = rdb.openFile()
		}
	} else {
		f, _ = os.OpenFile(rdb.fileName, os.O_RDWR, os.FileMode(os.O_TRUNC))
		if err = rdb.f.Truncate(0); err != nil {
			fmt.Printf("overwriteFile error truncating: \n%v\n", err)
		}
		f.Chmod(os.ModeAppend)
	}
	return f, nil
}

func (rdb *db) writeClose() (n int, err error) {
	n, err = rdb.f.Write([]byte{0xFF})
	if err != nil {
		switch err {
		case os.ErrClosed:
		}
	}
	return n, err
}

func (db *DBWriter) Write(rd RData) {
}

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
