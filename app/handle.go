package main

import (
	"fmt"
	"strconv"
	"strings"
)

func Read(b []byte) (n int, buff []byte) {
	char := b[0]
	switch char {
	case ARRAY:
		crlf, err := getCRLFOffset(b)
		if err != nil {
			return 0, buff
		}

		endCRLF := crlf + 2
		count, err := strconv.Atoi(string(b[1:crlf]))
		if err != nil {
			fmt.Println(err)
		}

		data := b[endCRLF:]

		lgth, rr := ReadRESP(data[:])
		cmd := strings.ToLower(string(rr.Data))
		switch cmd {
		case "echo":
			_, r := ReadRESP(data[lgth:])
			t := AppendBulkString(make([]byte, 0), r.Data)
			return len(t), t
		case "ping":
			r := AppendString(make([]byte, 0), "PONG")
			return len(r), r
		case "get":
			_, rk := ReadRESP(data[lgth:])
			r := store.Get(string(rk.Data))
			res := []byte("$-1\r\n")
			if r.Key != "" {
				res = AppendBulkString(make([]byte, 0), (r.Value))
			}
			return len(res), res
		case "set":
			n = 0
			kv := make([]RESP, count-1)
			data = data[lgth:]
			for j := 0; j < len(kv); j++ {
				l, r := ReadRESP(data[n:])
				kv[j] = r
				n += l
			}
			ok := false
			if len(kv) > 2 && strings.ToLower(string(kv[2].Data)) == "px" {
				r := kv[3]
				expire, err := strconv.Atoi(string(r.Data))
				if err != nil {
					fmt.Println("Error parsing expiry: ", err)
					return 0, make([]byte, 0)
				}
				ok = store.SetEx(string(kv[0].Data), string(kv[1].Type), kv[1].Data, int64(expire))
			} else {
				ok = store.Set(string(kv[0].Data), string(kv[1].Type), kv[1].Data)
			}
			if !ok {
				return 0, make([]byte, 0)
			}
			res := AppendString(make([]byte, 0), "OK")
			return len(res), res
		}
	}
	return n, buff
}
