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

const port = ":8080"

var room *chat.Room

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/ws", handleWebSocket)

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	log.Println("you are in room:", query.Get("room"))
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if room == nil {
		room = chat.NewRoom()
		go room.Run()
	}

	client := chat.NewClient(ws, room)
	subscription := chat.NewSubscription(query.Get("room"), client)
	client.Subscribe(subscription)

	go subscription.Read()
	go subscription.Write()
}
