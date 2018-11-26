package controller

import (
	"net/http"

	"github.com/alextanhongpin/go-chat/repository"
	"github.com/alextanhongpin/go-chat/ticket"
)

// M represents a generic map.
type M map[string]interface{}

// Controller represents the transport layer for the service.
type Controller struct{}

// New returns a pointer to a new Controller.
func New() *Controller {
	return new(Controller)
}

// PostAuthorize handles authorization for the user.
func (c *Controller) PostAuthorize(signer ticket.Dispenser, repo repository.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postAuthorize(signer, repo, w, r)
	}
}

// GetConversations returns a list of conversations.
func (c *Controller) GetConversations(repo repository.Conversation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getConversations(repo, w, r)
	}
}

// GetRooms endpoint returns a list of rooms.
func (c *Controller) GetRooms(repo repository.Room) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getRooms(repo, w, r)
	}
}
