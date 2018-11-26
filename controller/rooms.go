package controller

import (
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

func getRooms(repo repository.Room, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(entity.ContextKeyUserID).(string)
	if userID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	rooms, err := repo.GetRooms(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(M{
		"data": rooms,
	})
}
