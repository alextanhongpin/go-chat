package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/alextanhongpin/go-chat/chat"
	"github.com/alextanhongpin/go-chat/ticket"
)

const (
	port         = ":4000"
	redisPort    = ":6379"
	redisChannel = "chat"
)

func main() {

	cs := chat.NewServer(redisPort, redisChannel)

	go cs.Run()
	go cs.Subscribe()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/ws", cs.ServeWS())

	mux.HandleFunc("/auth", handleAuth)
	// mux.HandleFunc("/chat-histories", handleHistory)
	go checkGoroutine()

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func checkGoroutine() {
	var ch chan int
	if false {
		ch = make(chan int, 1)
		ch <- 1
	}
	go func(ch chan int) {
		<-ch
	}(ch)
	c := time.Tick(1 * time.Second)

	go func() {
		for range c {
			fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
		}
	}()
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	// Query user from the database based on the Authorization Bearer Token provided
	// Use the user_id obtained to create a new "Ticket" for the websocket
	userID := "abc123"

	// Create new ticket
	tic := ticket.New(userID, 1*time.Hour)

	// Sign ticket
	token, err := ticket.Sign(tic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return as json response
	fmt.Fprintf(w, `{"ticket": "%s"}`, token)
}
