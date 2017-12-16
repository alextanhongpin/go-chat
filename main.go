package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/alextanhongpin/go-chat/chat"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("listening to port *:8080. press ctrl + c to cancel.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// defer ws.Close()

	room := chat.NewRoom()
	go room.Run()

	client := &chat.Client{
		Conn: ws,
		Room: room,
		Send: make(chan chat.Message),
	}
	subscription := &chat.Subscription{
		Client: client,
		Room:   "123",
	}
	client.Room.Register <- subscription
	go subscription.Read()
	go subscription.Write()

	// for {
	// 	var msg chat.Message
	// 	if err := ws.ReadJSON(&msg); err != nil {
	// 		log.Println("websocket closed")
	// 		break
	// 	}
	// 	log.Println("got message:", msg)
	// }

	// ws.WriteMessage(websocket.CloseMessage, []byte{})
}
