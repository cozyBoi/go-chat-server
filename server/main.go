package main

import (
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
var chatCache map[int]*queue        //cache 30 chats

func broadcast_msg(conn *websocket.Conn, msg []byte, roomId int) {
	connsArr := conns[roomId]
	for _, curr_conn := range connsArr {
		if curr_conn == conn {
			continue
		}
		curr_conn.WriteMessage(1, msg)
	}
}

func connClose(c *websocket.Conn, rid int) {
	for i, curr_conn := range conns[rid] {
		if curr_conn == c {
			conns[rid] = append(conns[rid][:i], conns[rid][i+1:]...) //... => unpack the slice
		}
	}
	c.Close()
}

func socketHandler(ctx echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, _ := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)

	roomIdStr := ctx.Param("id")
	roomId, _ := strconv.Atoi(roomIdStr)
	defer connClose(c, roomId)
	_, flag := conns[roomId]
	if !flag {
		conns[roomId] = make([]*websocket.Conn, 0, 10)
		chatCache[roomId] = NewQ()
	}

	conns[roomId] = append(conns[roomId], c)

	for {
		_, msg, _ := c.ReadMessage() //get msg type and msg
		if chatCache[roomId].IsFull() {
			chatCache[roomId].Pop()
		}
		chatCache[roomId].Push(string(msg))
		broadcast_msg(c, msg, roomId)
	}
}

var roomNumber int = 5
var chatRooms = []string{"1", "2", "3", "4", "5"}

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

func sendPrevChats(ctx echo.Context) error {
	var curr_error error
	roomIdStr := ctx.Param("id")
	roomId, _ := strconv.Atoi(roomIdStr)

	if chatCache[roomId] == nil {
		return ctx.NoContent(http.StatusOK)
	} else if chatCache[roomId].size == 0 {
		return ctx.NoContent(http.StatusOK)
	}

	for i := 0; i < chatCache[roomId].size; i++ {
		curr_idx := (chatCache[roomId].front + i) % 30
		//fmt.Println(chatCache[roomId].buf[curr_idx])
		if i == chatCache[roomId].size-1 {
			curr_error = ctx.String(http.StatusOK, chatCache[roomId].buf[curr_idx])
		} else {
			curr_error = ctx.String(http.StatusOK, chatCache[roomId].buf[curr_idx]+",")
		}
	}
	return curr_error
}

func changeRoomHandler(ctx echo.Context) error {
	return ctx.File("assets/comm.html")
}

func roomsCreate(ctx echo.Context) error {
	roomNumber++
	chatRooms = append(chatRooms, strconv.Itoa(roomNumber))
	return ctx.NoContent(http.StatusOK)
}

func main() {
	conns = make(map[int][]*websocket.Conn)
	chatCache = make(map[int]*queue)
	e := echo.New()

	e.Static("/", "assets")
	e.Static("/rooms", "assets")

	e.File("/", "assets/main.html")
	e.File("/rooms/:id", "assets/comm.html")

	e.GET("/rooms/:id/ws", socketHandler)
	e.GET("/rooms/:id/chats", sendPrevChats)
	e.GET("/rooms", roomsHandler)

	e.POST("/rooms", roomsCreate)

	e.Start(":9100")
}
