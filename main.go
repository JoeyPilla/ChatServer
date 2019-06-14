package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type User struct {
	User   string `json:"user"`
	Action string `json:"action"`
}



func reader(conn *websocket.Conn) {
	//Start a go routine to listen
	defer conn.Close()
	defer delete(clients, conn.RemoteAddr())
	for {
		messageType, p, err := conn.ReadMessage()
		data := &User{}
		fmt.Printf("\rMessage Type: %d, Message String: %s\n", messageType, string(p))
		if(messageType === -1) {
			break
		}
		if err != nil {
			log.Println("Error conn.ReadMessage: ", err)
			break
		}
		err = json.Unmarshal(p, data)
		if err != nil {
			log.Println("Error JSON.Unmarshal",err, p, data)
			break
		}
		fmt.Println(data.User)
		fmt.Println(data.Action)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			break
		}
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

var clients = make(map[net.Addr]bool)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	clients[ws.RemoteAddr()] = true
	fmt.Println(clients)
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)
}

func main() {
	setupRoutes()
	fmt.Println("Go Websockets")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
