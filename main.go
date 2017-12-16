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

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
