package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestEncoder(t *testing.T) {
	b := new(bytes.Buffer)
	en := NewEncoder(b)
	foo := "foo"
	en.EncodeString([]byte(foo))
	got := b.Bytes()
	if got[0]&0xC0 != 0x0 {
		t.Errorf("en.Encode(foo) = %b; want 0000 0000 bits ", got)
	}
	if string(got[1:]) != foo {
		t.Errorf("en.Encode(foo) = %v; want foo ", string(got))
	}
	raw := strconv.Itoa(int(got[0]&0x3F)) + string(got[1:])
	if raw != "3"+foo {
		t.Errorf("encodeString(foo) -> UintToInt(b[0]) = %v; want 3foo ", raw)
	}

	foo64 := "fooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo"
	b = new(bytes.Buffer)
	en = NewEncoder(b)
	en.EncodeString([]byte(foo64))
	got = b.Bytes()
	if got[0]&0xC0 != 0x40 {
		t.Errorf(("encodeString(foo64) = %b; want 0100 0000 bits"), got)
	}
	tmp := int(0x3F&got[0])<<8 + int(got[1])
	l := strconv.Itoa(tmp)
	raw = l + string(got[2:])
	if raw != strconv.Itoa(len((foo64)))+foo64 {
		t.Errorf("encodeString(foo64) = \n%v; \nwant \n%d%v\n", raw, len(foo64)+2, foo64)
	}
}

func TestDb(t *testing.T) {
	r := NewRDB("test.txt")
	_, got := r.openFile()
	if got != nil && got != os.ErrNotExist {
		fmt.Println("db.openFile successfully did not find the file")
	} else {
		fmt.Println("db.openFile failed: ", got)
	}

	r.f, got = r.createFile()
	if got != nil {
		fmt.Println("db.createFile failed, file not created: ", got)
	} else {
		fmt.Println("db.createFile file successfully created: ")
	}

	// r.f, got = r.createFile()
	// if got != nil {
	// 	if got == os.ErrExist {
	// 		fmt.Println("db.createFile successfully returned exist error: ")
	// 	} else {
	// 		fmt.Print("db.createFile failed: ", got)
	// 	}
	// }
	defer r.f.Close()
	got_n, err := r.w.writeHeader(r.f)
	if err != nil {
		fmt.Println("Write Header failed: ", err)
	}
	if got_n != len([]byte("REDIS")) {
		fmt.Println("Write len != input: ", got_n)
	}
	got_n, err = r.writeClose()
	if err != nil {
		fmt.Println("WriteClose failed: ", err)
	}
	if got_n != 1 {
		fmt.Printf("WriteClose = %v,  want %v", got_n, 1)
	}

	foo := []byte("fooWrite")
	got_n, err = r.WriteString(foo)
	if err != nil {
		fmt.Printf("WriteString failed = %v, want %v", got_n, len(foo))
	}
	b := make([]byte, 1024)
	f, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	n, err := f.Read(b)
	if err != nil {
		fmt.Println("read failed ", got_n, err)
	}
	fmt.Printf("got %v\n", string(b[:n]))

}
