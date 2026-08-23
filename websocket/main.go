package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 代表一個 websocket 連線
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub 管理所有 client
type Hub struct {
	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client

	mutex sync.Mutex
}

func NewHub() *Hub {

	return &Hub{
		clients: make(map[*Client]bool),

		broadcast: make(chan []byte),

		register: make(chan *Client),

		unregister: make(chan *Client),
	}
}

// Hub 主循環
func (oSelf *Hub) Run() {

	for {

		select {

		case client := <-oSelf.register:

			oSelf.mutex.Lock()

			oSelf.clients[client] = true

			oSelf.mutex.Unlock()

			log.Println("client connected")

		case client := <-oSelf.unregister:

			oSelf.mutex.Lock()

			if _, ok := oSelf.clients[client]; ok {

				delete(oSelf.clients, client)

				close(client.send)
			}

			oSelf.mutex.Unlock()

		case message := <-oSelf.broadcast:

			oSelf.mutex.Lock()

			for client := range oSelf.clients {

				select {

				case client.send <- message:

				default:

					close(client.send)

					delete(oSelf.clients, client)
				}
			}

			oSelf.mutex.Unlock()
		}
	}
}

// 讀取 client 消息
func (oSelf *Client) readPump(hub *Hub) {

	defer func() {

		hub.unregister <- oSelf

		oSelf.conn.Close()

	}()

	for {

		_, message, err := oSelf.conn.ReadMessage()

		if err != nil {

			break
		}

		log.Println("receive:", string(message))

		// 廣播
		hub.broadcast <- message
	}
}

// 發送消息
func (oSelf *Client) writePump() {

	defer oSelf.conn.Close()

	for message := range oSelf.send {

		err := oSelf.conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		if err != nil {

			break
		}
	}
}

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {

		return
	}

	client := &Client{

		conn: conn,

		send: make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()

	go client.readPump(hub)
}

func main() {

	hub := NewHub()

	go hub.Run()

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			wsHandler(hub, w, r)
		},
	)

	log.Println("server :4031")

	err := http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {

		panic(err)
	}
}
