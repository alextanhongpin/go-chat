package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/alextanhongpin/go-chat/chat"
)

const port = ":3000"
const redisPort = ":6379"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := chat.NewServer(":6379")
	go cs.Run(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/ws", cs.ServeWS())

	// mux.HandleFunc("/auth", handleAuth)
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

// func handleAuth(w http.ResponseWriter, r *http.Request) {
// 	// Create a jwt `ticket` that contains the user scope and id to allow them to connect to the websocket
// }
