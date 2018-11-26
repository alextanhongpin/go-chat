package controller

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/alextanhongpin/go-chat/repository"
)

var pattern = regexp.MustCompile(`^\/conversations\/([\w+])\/?$`)

func getConversations(repo repository.Conversation, w http.ResponseWriter, r *http.Request) {
	submatches := pattern.FindStringSubmatch(r.URL.Path)
	if submatches == nil {
		http.Error(w, "room_id is required", http.StatusBadRequest)
		return
	}
	roomID := submatches[1]
	conversations, err := repo.GetConversations(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(M{
		"data": conversations,
		"room": roomID,
	})
}
