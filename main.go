package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

//EchoServer  Echo everything back.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

// TimeServer Send clock ticks.
func TimeServer(ws *websocket.Conn) {
	for {
		if _, err := ws.Write([]byte(time.Now().String())); err != nil {
			log.Println("Write: ", err)
			return
		}
		time.Sleep(1e+9)
	}
}

// MsgServer recieves a msg and sends it back to all users.
func MsgServer(ws *websocket.Conn) {
	var message string
	for {
		websocket.Message.Receive(ws, &message)
		if _, err := ws.Write([]byte(message)); err != nil {
			log.Println("SendingBack: ", err)
		}
	}
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/echo", websocket.Handler(EchoServer))
	http.Handle("/time", websocket.Handler(TimeServer))
	http.Handle("/msg", websocket.Handler(MsgServer))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}
}
