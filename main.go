package main

import (
	"log"
	"net/http"

	"github.com/alextanhongpin/go-chat/chat"
)

const port = ":3000"
const redisPort = ":6379"

func main() {
	// Create a new pool and pass it in as an option

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/ws", chat.Server(redisPort))

	// mux.HandleFunc("/auth", handleAuth)

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)

	log.Fatal(http.ListenAndServe(port, mux))
}

// func handleAuth(w http.ResponseWriter, r *http.Request) {
// 	// Create a jwt `ticket` that contains the user scope and id to allow them to connect to the websocket
// }
