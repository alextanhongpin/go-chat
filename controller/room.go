package controller

import (
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/service"

	"github.com/julienschmidt/httprouter"
)

// GetRooms endpoint returns a list of rooms.
func (c *Controller) GetRooms(svc service.GetRooms) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res, err := svc(r.Context(), service.GetRoomsRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

func (c *Controller) PostRooms(svc service.PostRooms) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		userID, _ := ctx.Value(entity.ContextKeyUserID).(string)
		var req service.PostRoomsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.UserID = userID
		res, err := svc(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}
