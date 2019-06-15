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

type Request struct {
	User *User 			`json:"user"`
	Message string `json:"message"`
}

type Response struct {
	User string 			`json:"user"`
	Message string `json:"message"`
	Client *websocket.Conn `json:"-"`
}

var (
	entering = make (chan *websocket.Conn)
	leaving = make (chan *websocket.Conn)
	messages = make (chan Response)
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}



func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket Connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrader Error: ", err)
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
			Message: "",
		}
		if(messageType == -1) {
			leaving <- conn
			break
		}
		if err != nil {
			log.Println("Error conn.ReadMessage: ", err)
			leaving <- conn
			break
		}
		err = json.Unmarshal(p, Request)
		if err != nil {
			log.Println("Error JSON.Unmarshal",err, p, Request)
			break
		}
		resp := Response{
			User: Request.User.User,
			Message: Request.Message,
			Client: conn,
		}
		switch Request.User.Action {
		case "Enter Room":
			messages<- resp
		case "Message":
			messages<- resp
		case "Leave Room":
			messages<- resp
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	setupRoutes()
	fmt.Println("Go Websockets")
	go broadcaster()
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func broadcaster() {
	var clients = make(map[*websocket.Conn]bool)
	for {
		select {
		case resp := <-messages:
				for cli := range clients {
					if (resp.Client != cli) {
						byteResp, err := json.Marshal(resp);
						if err != nil {
							log.Println(err)
							continue
						}
						if err := cli.WriteMessage(1, byteResp); err!=nil{
							log.Println(err)
						}
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