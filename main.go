package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
		dbUser    = os.Getenv("DB_USER")
		dbPass    = os.Getenv("DB_PASS")
		dbName    = os.Getenv("DB_NAME")
		jwtSecret = os.Getenv("JWT_SECRET")
		jwtIssuer = os.Getenv("JWT_ISSUER")
	)

	db, err := database.New(dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ticketMachine := ticket.NewMachine([]byte(jwtSecret), jwtIssuer, 5*time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	// mux.HandleFunc("/ws", cs.ServeWS())
	s := server.New(db)
	defer s.Close()

	mux.HandleFunc("/ws", s.ServeWS(ticketMachine, db))
	mux.HandleFunc("/auth", handleAuth(ticketMachine, db))
	// mux.HandleFunc("/chat-histories", handleHistory)
	// go checkGoroutine()

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
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
