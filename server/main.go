package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
} // use default options

var chat_log []string

var conns map[int][]*websocket.Conn //chatroom id => map => conn[]

func broadcast_msg(conn *websocket.Conn, msg []byte, roomId int) {
	connsArr := conns[roomId]
	for _, curr_conn := range connsArr {
		if curr_conn == conn {
			continue
		}
		curr_conn.WriteMessage(1, msg)
	}
}

func socketHandler(ctx echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, _ := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)

	defer c.Close()

	roomIdStr := ctx.Param("id")
	roomId, _ := strconv.Atoi(roomIdStr)
	_, flag := conns[roomId]
	if !flag {
		conns[roomId] = make([]*websocket.Conn, 0, 10)
	}

	conns[roomId] = append(conns[roomId], c)

	for {
		_, msg, _ := c.ReadMessage() //get msg type and msg
		chat_log = append(chat_log, string(msg))
		broadcast_msg(c, msg, roomId)
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
	conns = make(map[int][]*websocket.Conn)
	e := echo.New()

	e.Static("/", "assets")
	e.Static("/rooms", "assets")

	e.File("/", "assets/main.html")
	e.File("/rooms/:id", "assets/comm.html")

	e.GET("/rooms/:id/ws", socketHandler)
	e.GET("/rooms", roomsHandler)

	e.POST("/rooms", roomsCreate)

	e.Start(":9100")
}
