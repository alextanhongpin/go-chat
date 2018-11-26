package controller

import (
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/go-chat/repository"
	"github.com/alextanhongpin/go-chat/ticket"
)

type postAuthRequest struct {
	UserID string `json:"user_id"`
}

type postAuthResponse struct {
	Token string `json:"token"`
}

func postAuthorize(dispenser ticket.Dispenser, repo repository.User, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req postAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repo.GetUserByName(req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new ticket.
	ticket := dispenser.New(user.ID)

	// Sign ticket.
	token, err := dispenser.Sign(ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return as json response
	json.NewEncoder(w).Encode(postAuthResponse{
		Token: token,
	})
}
