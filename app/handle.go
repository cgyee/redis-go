package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Redis struct {
	commands map[string]func(b []byte) (int, []byte)
}

type RedisCommand struct {
	name    string
	handler func(b []byte) (int, []byte)
}

var r *Redis
var o sync.Once

func init() {
	r = Register()
}

func (r *Redis) register(c RedisCommand) {
	_, ok := r.commands[c.name]
	if ok {
		return
	}
	r.commands[c.name] = c.handler
}

func Register() *Redis {
	o.Do(func() {
		r = &Redis{commands: map[string]func(b []byte) (int, []byte){}}
	})
	r.register(RedisCommand{name: "ping", handler: Ping})
	r.register(RedisCommand{name: "get", handler: Get})
	r.register(RedisCommand{name: "echo", handler: Echo})
	r.register(RedisCommand{name: "set", handler: Set})
	r.register(RedisCommand{name: "config", handler: Config})
	return r
}

func Echo(b []byte) (n int, packet []byte) {
	_, r := ReadRESP(b)
	packet = AppendBulkString(make([]byte, 0), r.Data)
	return len(packet), packet
}

func Ping(b []byte) (n int, packet []byte) {
	p := AppendString(make([]byte, 0), "PONG")
	return len(p), p
}

func Get(b []byte) (n int, packet []byte) {
	_, r := ReadRESP(b)
	data := store.Get(string(r.Data))
	p := []byte("$-1\r\n")
	if data.Key != "" {
		p = AppendBulkString(make([]byte, 0), (data.Value))
	}
	return len(p), p
}

func Config(b []byte) (n int, packet []byte) {
	n, resp := ReadRESP(b)
	if strings.ToLower(string(resp.Data)) == "get" {
		_, resp = ReadRESP(b[n:])
		data := store.Get(string(resp.Data))
		res := []byte("$-1\r\n")
		if data.Key != "" {
			packet = AppendArray(packet, 2)
			packet = AppendBulkString(packet, []byte(data.Key))
			res = AppendBulkString(packet, (data.Value))
		}
		return len(res), res
	}
	return 0, make([]byte, 0)
}

func Set(b []byte) (n int, packet []byte) {
	n = 0
	kv := make([]RESP, 2)
	for j := 0; j < len(kv); j++ {
		l, resp := ReadRESP(b[n:])
		kv[j] = resp
		n += l
	}
	ok := false
	fmt.Println(n < len(b), n, len(b))
	for n < len(b) {
		l, resp := ReadRESP(b[n:])
		kv = append(kv, resp)
		n += l

	}
	if len(kv) > 2 && string(kv[2].Data) == "px" {
		expire, err := strconv.Atoi(string(kv[3].Data))
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

func ReadCommand(b []byte) (cmd Cmd) {
	switch b[0] {
	case ARRAY, BULK:
		crlf, err := getCRLFOffset(b)
		if err != nil {
			return cmd
		}
		endCRLF := crlf + 2

		n, r := ReadRESP(b[endCRLF:])
		fmt.Println(string(r.Data), string(b[endCRLF:]))
		return Cmd{name: strings.ToLower(string(r.Data)), raw: r.Raw, data: b[endCRLF+n:]}

	}
	return cmd
}

func Handle(b []byte) (n int, data []byte) {
	c := ReadCommand(b)
	if handler, ok := r.commands[c.name]; ok {
		n, data := handler(c.data)
		return n, data
	}
	return 0, make([]byte, 0)
}
