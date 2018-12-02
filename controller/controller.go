package controller

import (
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/go-chat/service"
	"github.com/julienschmidt/httprouter"
)

// M represents a generic map.
type M map[string]interface{}

// Controller represents the transport layer for the service.
type Controller struct {
	// endpoints
}

// New returns a pointer to a new Controller.
func New() *Controller {
	return new(Controller)
}

// PostAuthorize handles authorization for the user.
func (c *Controller) PostAuthorize(svc service.Authorize) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req service.AuthorizeRequest
		res, err := svc(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Return as json response
		json.NewEncoder(w).Encode(res)
	}
}

// GetConversations returns a list of conversations.
func (c *Controller) GetConversations(svc service.GetConversations) httprouter.Handle {
	// var pattern = regexp.MustCompile(`^\/conversations\/([\w+])\/?$`)
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// submatches := pattern.FindStringSubmatch(r.URL.Path)
		// if submatches == nil {
		//         http.Error(w, "room_id is required", http.StatusBadRequest)
		//         return
		// }
		// roomID := submatches[1]
		res, err := svc(r.Context(), service.GetConversationsRequest{RoomID: ps.ByName("id")})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

// PostLogin handles the user authentication request.
func (c *Controller) PostLogin(svc service.Login) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var request service.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := svc(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

// PostRegister handles the user registration request.
func (c *Controller) PostRegister(svc service.Register) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var request service.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := svc(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}
