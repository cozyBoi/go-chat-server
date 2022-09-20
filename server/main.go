package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

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

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/socket")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

//http://192.168.0.37:9100/

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var printL = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        d.className = "bubbleleft"
        output.appendChild(d);  
        output.scroll(0, output.scrollHeight);
    };
    var printR = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        d.className = "bubbleright"
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            //print("OPEN");
        }
        ws.onclose = function(evt) {
            //print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            printL("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            printL("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        //input.value.className = "bubble"
        printR("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
<style>

.conversation-container {
    margin: 0 auto;
    max-width: 400px;
    height: 600px;
    padding: 0 20px;
    border: 3px solid #f1f1f1;
    overflow: scroll;
  }
  .input-container {
    margin: 0 auto;
    max-width: 400px;
    height: 800px;
  }
  .bubble-container {
      text-align: left;
  }
  .bubbleright {
    border: 2px solid #f1f1f1;
    background-color: #7fbefc;
    border-radius: 5px;
    padding: 10px;
    margin: 10px 0;
      width: 230px;
      float: right;
  }
  .bubbleleft {
    border: 2px solid #f1f1f1;
    background-color: #fafdfd;
    border-radius: 5px;
    padding: 10px;
    margin: 10px 0;
      width: 230px;
      float: left;
  }
  .bubble {
      background-color: #abf1ea;
      border: 2px solid #87E0D7;
      float: left;
  }
  .name {
    padding-right: 8px;
      font-size: 11px;
  }
  ::-webkit-scrollbar {
    width: 10px;
  }
  ::-webkit-scrollbar-track {
    background: #f1f1f1;
  }
  ::-webkit-scrollbar-thumb {
    background: #888;
  }
  ::-webkit-scrollbar-thumb:hover {
    background: #555;
  }
</style>
</head>
    <body>
        <table>
            <tr><td valign="top" width="50%">
                <p>Click "Open" to create a connection to the server, 
                    "Send" to send a message to the server and "Close" to close the connection. 
                    You can change the message and send multiple times.
                <p>
                <form>
                    <button id="open">Open</button>
                    <button id="close">Close</button>
                    <p><input id="input" type="text" value="Hello world!">
                    <button id="send">Send</button>
                </form>
                </td><td valign="top" width="50%">
                <div id="output" class="conversation-container" style="max-height: 70vh;overflow-y: scroll;"></div>
            </td></tr>
        </table>
    </body>
</html>
`))
