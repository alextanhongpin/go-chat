package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alextanhongpin/go-chat/chat"
	"github.com/alextanhongpin/go-chat/database"
	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
	"github.com/alextanhongpin/go-chat/ticket"
)

func main() {
	var (
		dbUser    = os.Getenv("DB_USER")
		dbPass    = os.Getenv("DB_PASS")
		dbName    = os.Getenv("DB_NAME")
		jwtSecret = os.Getenv("JWT_SECRET")
		jwtIssuer = os.Getenv("JWT_ISSUER")
		port      = ":4000"
	)

	db, err := database.New(dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ticketDispenser := ticket.NewDispenser([]byte(jwtSecret), jwtIssuer, 5*time.Minute)

	s := chat.New(db)
	defer s.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	mux.HandleFunc("/ws", s.ServeWS(ticketDispenser, db))
	mux.HandleFunc("/auth", handleAuth(ticketDispenser, db))
	mux.HandleFunc("/rooms", authMiddleware(ticketDispenser)(handleGetRooms(db)))
	mux.HandleFunc("/conversations", authMiddleware(ticketDispenser)(handleGetConversations(db)))

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

type middleware func(http.HandlerFunc) http.HandlerFunc

func authMiddleware(dispenser ticket.Dispenser) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			auth := r.Header.Get("Authorization")
			if values := strings.Split(auth, " "); len(values) != 2 {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			} else {
				bearer, token := values[0], values[1]
				if bearer != "Bearer" {
					http.Error(w, "invalid bearer type", http.StatusUnauthorized)
					return
				}
				userID, err := dispenser.Verify(token)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx := r.Context()
				ctx = context.WithValue(ctx, entity.ContextKey("user_id"), userID)
				r = r.WithContext(ctx)

			}

			next.ServeHTTP(w, r)
		}
	}
}

func handleAuth(dispenser ticket.Dispenser, db database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("got user by name: %#v", user)

		// Create new ticket.
		ticket := dispenser.New(user.ID)

		// Sign ticket.
		token, err := dispenser.Sign(ticket)
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

func handleGetRooms(db repository.Room) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIDContext := ctx.Value(entity.ContextKey("user_id"))
		if userIDContext == nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		userID := userIDContext.(string)
		log.Println("getting rooms", userID)
		if userID == "" {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		rooms, err := db.GetRooms(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(getRoomsResponse{
			Data: rooms,
		})
	}
}

type getRoomsResponse struct {
	Data []entity.UserRoom `json:"data"`
}

func handleGetConversations(db repository.Conversation) http.HandlerFunc {
	pattern := regexp.MustCompile(`^\/conversations\/([\w+])\/?$`)
	// res := r.FindStringSubmatch("/users/1")
	return func(w http.ResponseWriter, r *http.Request) {
		submatches := pattern.FindStringSubmatch(r.URL.Path)
		if submatches == nil {
			http.Error(w, "room_id is required", http.StatusBadRequest)
			return
		}
		roomID := submatches[1]
		conversations, err := db.GetConversations(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(getConversationsResponse{
			Data: conversations,
			Room: roomID,
		})
	}
}

type getConversationsResponse struct {
	Data []entity.Conversation `json:"data"`
	Room string                `json:"room"`
}
