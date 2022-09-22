package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

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

func socketHandler(ctx echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //araboza
	c, _ := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	fmt.Println("connect socket")
	defer c.Close()

	conns = append(conns, c)

	for {
		_, msg, _ := c.ReadMessage() //get msg type and msg
		chat_log = append(chat_log, string(msg))
		broadcast_msg(c, msg)
	}
}

func main() {
	e := echo.New()
	e.Static("/", "assets")
	e.File("/", "assets/main.html")
	e.GET("/ws", socketHandler)
	e.Start(":9100")
}
