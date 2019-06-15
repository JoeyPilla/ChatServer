package main

import (
	"encoding/json"
	"fmt"
	"log"
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

type Message struct {
	Message   string `json:"message"`
}

type Request struct {
	User *User 			`json:"user"`
	Message *Message `json:"message"`
}



var (
	entering = make (chan *websocket.Conn)
	leaving = make (chan *websocket.Conn)
	messages = make (chan string)
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}



func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket Connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected")
	entering<- ws
	go clientWriter(ws)
}

func clientWriter(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		Request := &Request{
			User: &User{},
			Message: &Message{},
		}
		fmt.Printf("\rMessage Type: %d, Message String: %s\n", messageType, string(p))
		if(messageType == -1) {
			break
		}
		if err != nil {
			log.Println("Error conn.ReadMessage: ", err)
			break
		}
		err = json.Unmarshal(p, Request)
		if err != nil {
			log.Println("Error JSON.Unmarshal",err, p, Request)
			break
		}

		switch Request.User.Action {
		case "Enter Room":
			messages<- Request.Message.Message
		case "Message":
			messages<- Request.Message.Message
		case "Leave Room":
			messages<- Request.Message.Message
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)
}

func main() {
	setupRoutes()
	fmt.Println("Go Websockets")
	go broadcaster()
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func broadcaster() {
	var clients = make(map[*websocket.Conn]bool)
	for {
		select {
		case msg := <-messages:
				for cli := range clients {
					 if err := cli.WriteMessage(1, []byte(msg)); err!=nil{
						 log.Println(err)
					 }
				}
		case cli := <-entering:
				clients[cli] = true
		case cli := <-leaving:
				delete(clients, cli)
				cli.Close()
		}
	}
}