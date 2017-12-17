package main

import (
	"log"
	"net/http"

	"github.com/alextanhongpin/go-chat/chat"
)

const port = ":3000"

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/ws", chat.Server())

	mux.HandleFunc("/auth", handleAuth)

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	// Create a jwt `ticket` that contains the user scope and id to allow them to connect to the websocket
}
