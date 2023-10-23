package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Chat room to manage connected clients
type ChatRoom struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
}

func (room *ChatRoom) run() {
	for {
		select {
		case message := <-room.broadcast:
			// Send the received message to all connected clients
			for client := range room.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error: %v", err)
					client.Close()
					delete(room.clients, client)
				}
			}
		}
	}
}

func handleChat(w http.ResponseWriter, r *http.Request, room *ChatRoom) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	room.clients[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: %v", err)
			delete(room.clients, conn)
			break
		}

		room.broadcast <- message
	}
}

func main() {
	room := &ChatRoom{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
	go room.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleChat(w, r, room)
	})

	port := 8080
	fmt.Printf("Chat server started on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
