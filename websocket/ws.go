package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
)

type Socket struct {
	key string
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	_, err := fmt.Fprintln(w, "Hi")
	if err != nil {
		return
	}

	/*
		Connection update metrics:

		-> header["Connection"] is equal "Upgrade"
		-> header["Upgrade"] is equal "websocket"

		-> Sec-Websocket-Key: <websocket key> NOT NULL
	*/
	if r.Header.Get("Upgrade") != "websocket" {
		return
	}
	if r.Header.Get("Connection") != "Upgrade" {
		return
	}

	sock := Socket{}
	sock.key = r.Header.Get("Sec-Websocket-Key")

	if sock.key == "" {
		return
	}

	result := sock.Handshake(w)
	fmt.Println(result)
}

func (sock Socket) CalculateAcccept() string {
	GUID := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hashKey := sha1.Sum([]byte(sock.key + GUID))

	return base64.StdEncoding.EncodeToString(hashKey[:])
}

func (sock Socket) Handshake(w http.ResponseWriter) bool {
	hj, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println("Hijacker fault open")
		return false
	}

	conn, buffrw, err := hj.Hijack()
	if err != nil {
		fmt.Println("Hijacker fault open 2")
		return false
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Fatal error")
		}
	}(conn)

	accept_key := sock.CalculateAcccept()

	buffrw.WriteString("HTTP/1.1 101 Switching Protocols\\r\\n")
	buffrw.WriteString("Upgrade: websocket\r\n")
	buffrw.WriteString("Connection: Upgrade\r\n")
	buffrw.WriteString("Sec-Websocket-Accept: " + accept_key + "\r\n\r\n")
	buffrw.Flush()

	buf := make([]byte, 1024)
	for {
		n, err := buffrw.Read(buf)
		if err != nil {
			return false
		}
		fmt.Println(buf[:n])
	}
}
