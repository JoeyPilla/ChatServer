package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

//start at 1 because -1 happens when a new client connects
var counter int = 1 
var countChannel = make(chan int)

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

func handleCount(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket Connection
	counterSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrader Error: ", err)
	}
	log.Println("Counter Connected")

	for {
		//freeze until value is sent to countChannel
		count := <-countChannel
		s := strconv.Itoa(count)
		if err:= counterSocket.WriteMessage(1, []byte(s)); err!=nil {
			log.Println(err)
			counterSocket.Close()
			break
		}
	}
}

func handleChatSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket Connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrader Error: ", err)
	}
	log.Println("Client Connected")

	//tell broadcaster a new client has entered
	entering<- ws
	//open a new go routine for new client
	go clientWriter(ws)
}

func clientWriter(conn *websocket.Conn) {
	//make sure cleans up when client closes
	defer func() {
		leaving <- conn
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		Request := &Request{
			User: &User{},
			Message: "",
		}
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

		resp := Response{
			User: Request.User.User,
			Message: Request.Message,
			Client: conn,
		}

		//TODO unnecessary??? 
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
	http.HandleFunc("/count", handleCount)
	http.HandleFunc("/ws", handleChatSocket)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	setupRoutes()
	fmt.Println("Go Chat Server Launched")
	go broadcaster()
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func broadcaster() {
	var clients = make(map[*websocket.Conn]bool)
	for {
		select {
		case resp := <-messages:
			counter++
			countChannel <- counter
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
			//ignore message for new client
			counter--
			clients[cli] = true
		case cli := <-leaving:
			//clean up
			delete(clients, cli)
			cli.Close()
		}
	}
}