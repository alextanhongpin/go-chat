package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/alextanhongpin/go-chat/database"
	"github.com/alextanhongpin/go-chat/server"
	"github.com/alextanhongpin/go-chat/ticket"
)

const (
	port         = ":4000"
	redisPort    = ":6379"
	redisChannel = "chat"
)

func main() {
	var (
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASS")
		dbName = os.Getenv("DB_NAME")
		// port   = os.Getenv("PORT")
	)
	db, err := database.New(dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ticketMachine := ticket.NewMachine([]byte("secret"), "go-chat", 5*time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	// mux.HandleFunc("/ws", cs.ServeWS())
	s := server.New()
	defer s.Close()

	mux.HandleFunc("/ws", s.ServeWS(ticketMachine, db))
	mux.HandleFunc("/auth", handleAuth(ticketMachine, db))
	// mux.HandleFunc("/chat-histories", handleHistory)
	// go checkGoroutine()

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

func handleAuth(machine ticket.Dispenser, db database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query user from the database based on the Authorization Bearer Token provided
		// Use the user_id obtained to create a new "Ticket" for the websocket

		if r.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: Validate if the user is a valid user by checking the database.

		user, err := db.GetUserByName(req.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Create new ticket.
		ticket := machine.New(strconv.Itoa(user.ID))

		// Sign ticket.
		token, err := machine.Sign(ticket)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Return as json response
		json.NewEncoder(w).Encode(authResponse{
			Token: token,
		})
	}
}

type authRequest struct {
	UserID string `json:"user_id"`
}

type authResponse struct {
	Token string `json:"token"`
}
