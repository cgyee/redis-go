package main

import (
	"fmt"
	"strconv"
	"strings"
)

type RESPERROR struct{ message string }

func (m *RESPERROR) Error() string {
	return fmt.Sprintf("Error: %v", m.message)
}

type RESP struct {
	Type  byte
	Raw   []byte
	Data  []byte
	Count int
}

func (r RESP) Bytes() []byte {
	return r.Data
}

const (
	INTEGER = ':'
	STRING  = '+'
	BULK    = '$'
	ARRAY   = '*'
	MAP     = '%'
	ERROR   = '-'
)

// This finds the first carriage return(\r\n) and return the number of bytes from b[1:]
// to \r\n
func getCRLFOffset(b []byte) (i int, err error) {
	if len(b) == 0 {
		return 0, &RESPERROR{"Error finding CR: "}
	}
	i = 1
	for ; ; i++ {
		if b[i] == '\n' {
			if b[i-1] != '\r' {
				return 0, &RESPERROR{"Error finding CR: "}
			}
			return i - 1, nil
		}
	}
}

func Read(b []byte) (n int, buff []byte) {
	char := b[0]
	switch char {
	case ARRAY:
		crlf, err := getCRLFOffset(b)

		if err != nil {
			return 0, buff
		}
		count, err := strconv.Atoi(string(b[1:crlf]))
		if err != nil {
			fmt.Println(err)
		}

		data := b[crlf+2:]
		i := 0
		for ; i < count; i++ {
			ll, rr := ReadRESP(data[0:])
			cmd := strings.ToLower(string(rr.Data))
			switch cmd {
			case "echo":
				n, r := ReadRESP(data[ll+2:])
				t := AppendBulkString(make([]byte, 0), r.Data)
				return n, t
			case "ping":
				r := AppendString(make([]byte, 0), "PONG")
				return len(r), r
			}
		}
	}
	return n, buff
}

// Read the RESP protocol this runs after we determine the command
// We don't do and transformations here so we return the bytes
func ReadRESP(b []byte) (n int, resp RESP) {
	if len(b) == 0 {
		return 0, RESP{}
	}
	resp.Type = b[0]
	switch resp.Type {
	case INTEGER, BULK, ARRAY, STRING:
	default:
		return 0, RESP{}
	}
	// clrf start at index == \r, endCRLF starts at index \n + 1
	crlf, _ := getCRLFOffset(b)
	endCRLF := crlf + 2
	resp.Raw = b[:endCRLF]
	resp.Data = b[1:crlf]
	fmt.Println(string(resp.Data))
	switch resp.Type {
	case ARRAY:
		count, err := strconv.Atoi(string(resp.Data))
		if err != nil {
			fmt.Println("Error parsing array length: ", err)
			return 0, RESP{}
		}
		var sl int
		// [*, 2, \r, \n, $]
		data := b[endCRLF:]
		for el := 0; el < count; el++ {
			l, r := ReadRESP(data)
			if r.Type == 0 {
				return 0, RESP{}
			}
			data = data[l:]
			sl += l
		}
		resp.Raw = b[0:sl]
		resp.Data = b[1:sl]
		return len(resp.Raw), resp

	case BULK:
		count, err := strconv.Atoi(string(resp.Data))
		if err != nil {
			fmt.Println("Error parsing bulk string length: ", err)
			return 0, RESP{}
		}

		if count < 0 {
			resp.Data = nil
			count = 0
			return len(resp.Raw), RESP{}
		}

		if len(b) < crlf+count+2 {
			return 0, RESP{}
		}

		if b[endCRLF+count] != '\r' || b[endCRLF+count+1] != '\n' {
			return 0, RESP{}
		}
		resp.Raw = b[0 : crlf+count+2]
		resp.Data = b[endCRLF : endCRLF+count]
		count = 0
		return len(resp.Raw), resp
	}
	return 0, RESP{}
}

func AppendPrefix(b []byte, c byte, n int64) []byte {
	if 0 <= n && n <= 9 {
		return append(b, c, byte('0'+n), '\r', '\n')
	}
	b = append(b, c)
	b = strconv.AppendInt(b, n, 10)
	return append(b, '\r', '\n')
}

func AppendString(b []byte, s string) []byte {
	b = append(b, '+')
	b = append(b, strings.TrimSpace(s)...)
	return append(b, '\r', '\n')
}

func AppendInt(b []byte, n int64) []byte {
	b = append(b, ':')
	b = strconv.AppendInt(b, n, 10)
	return append(b, '\r', '\n')
}

func AppendBulkString(b []byte, bulk []byte) []byte {
	b = AppendPrefix(b, '$', int64(len(bulk)))
	b = append(b, bulk...)
	return append(b, '\r', '\n')
}
func AppendArray(b []byte, n int64) []byte {
	return AppendPrefix(b, 'n', n)
}
