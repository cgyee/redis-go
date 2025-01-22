package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestEncode(t *testing.T) {

	foo := "foo"
	got := encodeString([]byte(foo))
	m := got[0]
	if m&0xC0 != 0x0 {
		t.Errorf("encodeString(foo) = %b; want 0000 0000 bits ", got)
	}
	if string(got[1:]) != foo {
		t.Errorf("encodeString(foo) = %v; want foo ", string(got))
	}

	raw := strconv.Itoa(UintToInt([]byte{got[0]})) + string(got[1:])
	if raw != "3"+foo {
		t.Errorf("encodeString(foo) -> UintToInt(b[0]) = %v; want 3foo ", raw)
	}

	foo64 := "fooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo"
	got = encodeString([]byte(foo64))
	if got[0]&0xC0 != 0x40 {
		t.Errorf(("encodeString(foo64) = %b; want 0100 0000 bits"), got)
	}
	// l := strconv.Itoa(UintToInt([]byte{got[0] & 0x3F, got[1]}))
	tmp := int(0x3F&got[0])<<8 + int(got[1])
	l := strconv.Itoa(tmp)
	raw = l + string(got[2:])
	if raw != strconv.Itoa(len((foo64)))+foo64 {
		t.Errorf("encodeString(foo64) = %v; want %d%v\n", raw, len(foo64), foo64)
	}

}

func TestDecode(t *testing.T) {
	foo := "foo"
	fooEnc := encodeString([]byte(foo))
	got := decodeString(fooEnc)
	if string(got) != foo {
		t.Errorf("decodeString(\"foo\") = %v; want \"foo\"", string(got))
	}

	foo64 := "fooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo"
	fooEnc = encodeString([]byte(foo64))
	got = decodeString(fooEnc)
	if string(got) != foo64 {
		t.Errorf("decodeString(\"foo64\") = %v; want \"%v\"", string(got), foo64)

	}
}

func TestEncoder(t *testing.T) {
	foo := "foo"
	en := encoder{}
	got, _ := en.Encode(TypeString, []byte(foo))
	if got[0]&0xC0 != 0x0 {
		t.Errorf("en.Encode(foo) = %b; want 0000 0000 bits ", got)
	}
	if string(got[1:]) != foo {
		t.Errorf("en.Encode(foo) = %v; want foo ", string(got))
	}
	raw := strconv.Itoa(UintToInt([]byte{got[0]})) + string(got[1:])
	if raw != "3"+foo {
		t.Errorf("encodeString(foo) -> UintToInt(b[0]) = %v; want 3foo ", raw)
	}

	foo64 := "fooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo"
	got, err := en.Encode(TypeString, []byte(foo64))
	if err != nil {
		fmt.Println(err)
	}
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

// int((binary.BigEndian.Uint16([]byte{0x0, 0x0, got[0] & 0x63})))
