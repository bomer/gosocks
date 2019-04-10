package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// rooms := map[string]
var all map[string]*websocket.Conn
var allconnections []*websocket.Conn
var history []string

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

func remove(slice []*websocket.Conn, s int) []*websocket.Conn {
	return append(slice[:s], slice[s+1:]...)
}

// MsgServer recieves a msg and sends it back to all users.
func MsgServer(ws *websocket.Conn) {
	var message string
	// all.(ws)
	alive := true
	allconnections = append(allconnections, ws)
	for {
		websocket.Message.Receive(ws, &message)
		for i, aWs := range allconnections {

			if _, err := aWs.Write([]byte(message)); err != nil {
				log.Println("Error!!: ", err)
				aWs.Close()
				allconnections = remove(allconnections, i)
				alive = false

			}

		}
		log.Println("I hope I didn't crash: ")
		if !alive {
			break
		}
	}
}

func main() {

	//Init
	all = make(map[string]*websocket.Conn)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/echo", websocket.Handler(EchoServer))
	http.Handle("/time", websocket.Handler(TimeServer))
	http.Handle("/msg", websocket.Handler(MsgServer))

	err := http.ListenAndServe(":8002", nil)
	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}
}
