package main

import (
	"github.com/DmitriyMosk/go-fun/websocket"
	"net/http"
)

func main() {
	http.HandleFunc("/mysocket", websocket.WSHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return
	}
}
