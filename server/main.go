package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Broad bool   `json:broad`
	Data  string `json:data`
}

var addr = flag.String("addr", "0.0.0.0:9100", "http service address")

var upgrader = websocket.Upgrader{
	EnableCompression: true,
} // use default options

var chat_log []string

var conns []*websocket.Conn

func broadcast_msg(conn *websocket.Conn, msg []byte) {
	for _, curr_conn := range conns {
		if curr_conn == conn {
			continue
		}
		curr_conn.WriteMessage(1, msg)
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //araboza
	c, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("connect socket")
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	conns = append(conns, c)

	for { //busy waiting?
		_, msg, _ := c.ReadMessage() //get msg type and msg
		chat_log = append(chat_log, string(msg))
		broadcast_msg(c, msg)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/socket", socketHandler)
	//http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
