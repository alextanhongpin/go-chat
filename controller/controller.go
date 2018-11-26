package controller

import "net/http"

type Controller struct {
}

func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request)         {}
func (c *Controller) GetConversations(w http.ResponseWriter, r *http.Request) {}
func (c *Controller) GetRooms(w http.ResponseWriter, r *http.Request)         {}
