package main

import (
	"bytes"
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
}

type Writer struct {
	idx  int
	en   *Encoder
	buf  *bytes.Buffer
	mode os.FileMode
}

type db struct {
	f        *os.File
	fileName string
	w        Writer
}

func NewRDB(fileName string) *db {
	b := new(bytes.Buffer)
	en := &Encoder{w: b}
	w := Writer{en: en, buf: b}
	rdb := &db{w: w, fileName: fileName}
	return rdb

}

func (w *Writer) writeHeader(f *os.File) (n int, err error) {
	header := "REDIS"
	n, err = f.Write([]byte(header))
	if err != nil {
		fmt.Printf("error during writeHeader:\n%v\n", err)
		return n, err
	}
	return n, nil

	//  metadata section
}

func (rdb *db) openFile() (f *os.File, err error) {
	f, err = os.OpenFile(rdb.fileName, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	return f, nil
}

func (rdb *db) createFile() (f *os.File, err error) {
	f, err = os.OpenFile(rdb.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	return f, nil

}

func (rdb *db) overwriteFile() (f *os.File, err error) {
	f, err = os.OpenFile(rdb.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		switch err {
		case os.ErrExist:
			err = f.Truncate(0)
			return
		case os.ErrNotExist:
			f, err = os.Create(rdb.fileName)
			return
		default:
			return
		}
	}
	err = f.Truncate(0)
	return f, err
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

func (rdb *db) WriteString(data []byte) (n int, err error) {
	n, err = rdb.w.en.EncodeString(data)
	if err != nil {
		return
	}
	b := rdb.w.buf.Bytes()
	fmt.Println(string(b))
	n, err = rdb.write(b)
	return
}

func (rdb *db) write(data []byte) (n int, err error) {
	n, err = rdb.f.Write(data)
	if err != nil {
		switch err {
		case os.ErrNotExist:
			rdb.f, _ = rdb.createFile()
		case os.ErrClosed:
			rdb.f, _ = rdb.openFile()
		default:
			return 0, err
		}

		n, err = rdb.f.Write(data)
		if err != nil {
			return
		}
	}
	return n, nil
}
