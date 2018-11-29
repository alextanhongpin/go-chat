package controller

import (
	"encoding/json"
	"net/http"
	"regexp"
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
func (c *Controller) PostAuthorize(svc postAuthorizeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		var req postAuthRequest
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
func (c *Controller) GetConversations(svc getConversationsService) http.HandlerFunc {
	var pattern = regexp.MustCompile(`^\/conversations\/([\w+])\/?$`)
	return func(w http.ResponseWriter, r *http.Request) {
		submatches := pattern.FindStringSubmatch(r.URL.Path)
		if submatches == nil {
			http.Error(w, "room_id is required", http.StatusBadRequest)
			return
		}
		roomID := submatches[1]
		res, err := svc(r.Context(), getConversationsRequest{roomID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

// GetRooms endpoint returns a list of rooms.
func (c *Controller) GetRooms(svc getRoomsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := svc(r.Context(), getRoomsRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

// PostLogin handles the user authentication request.
func (c *Controller) PostLogin(svc loginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}
		var request loginRequest
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
func (c *Controller) PostRegister(svc registerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}
		var request registerRequest
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
