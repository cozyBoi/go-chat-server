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

var roomNumber int = 5
var chatRooms = [5]string{"1", "2", "3", "4", "5"}

func roomsHandler(ctx echo.Context) error {
	var curr_error error
	for i, room := range chatRooms {
		if i == roomNumber-1 {
			curr_error = ctx.String(http.StatusOK, room)
		} else {
			curr_error = ctx.String(http.StatusOK, room+",")
		}
	}
	return curr_error
}

func changeRoomHandler(ctx echo.Context) error {
	fmt.Println("hi")
	return ctx.File("assets/comm.html")
}

func roomsCreate(ctx echo.Context) error {
	var curr_error error
	curr_error = ctx.String(http.StatusOK, "roomHandler!")
	return curr_error
}

func main() {
	e := echo.New()

	e.Static("/", "assets")
	e.Static("/rooms", "assets")
	e.File("/", "assets/main.html")
	e.File("/rooms/1", "assets/comm.html")
	e.GET("/ws", socketHandler)
	e.GET("/rooms", roomsHandler)
	//js 에서는 버튼클릭하면 GET /rooms/:id 보내게
	//ws://192.168.0.37:9100/rooms/:id/ws 와 웹소켓을 열어야함

	e.POST("/rooms", roomsCreate)

	e.Start(":9100")
}
