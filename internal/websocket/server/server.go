package server

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type WSer struct {
	Conns map[*websocket.Conn]bool
}

func NewWSer() *WSer {
	return &WSer{
		Conns: make(map[*websocket.Conn]bool),
	}
}

func (s *WSer) HandleConnections(ws *websocket.Conn) {
	defer func (){
		fmt.Println("WebSocket disconnected:",ws.RemoteAddr())
		ws.Close()
	}()
	fmt.Println("WebSocket Connected:",ws.RemoteAddr())
	s.Conns[ws] = true

}

func (s *WSer) HandleMessages(ws *websocket.Conn) {
	buff :=make ([]byte, 1024)
	for {
		n,err:= ws.Read(buff)
		
		if err!=nil{
			if err==io.EOF {
				fmt.Println("WebSocket closed by client:",ws.RemoteAddr())
				break
			}
			fmt.Println("Error reading message:",err)
			continue
		}
		message :=string (buff[:n])
		fmt.Printf("Received message from %s: %s\n",ws.RemoteAddr(),message)
	}
}

func Start() {

	wsServer := NewWSer()
	http.Handle("/ws",websocket.Handler(wsServer.HandleConnections))
	fmt.Println("WebSocket server started on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("WebSocket server error:", err)
	}

}