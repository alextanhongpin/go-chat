package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/service"
	"github.com/julienschmidt/httprouter"
)

// PostFriendship creates a new friend request.
func (c *Controller) PostFriendship(svc service.AddFriend) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		id, _ := ctx.Value(entity.ContextKeyUserID).(string)
		userID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id = ps.ByName("id")
		targetID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req := service.AddFriendRequest{
			UserID:   userID,
			TargetID: targetID,
		}
		res, err := svc(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

func (c *Controller) PatchFriendship(svc service.HandleFriend) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		id, _ := ctx.Value(entity.ContextKeyUserID).(string)
		userID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id = ps.ByName("id")
		targetID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var req service.HandleFriendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.TargetID = targetID
		req.UserID = userID
		res, err := svc(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

func (c *Controller) GetContacts(svc service.GetContacts) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		id, _ := ctx.Value(entity.ContextKeyUserID).(string)
		userID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := svc(ctx, service.GetContactsRequest{UserID: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}
