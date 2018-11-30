package repository

import "github.com/alextanhongpin/go-chat/entity"

type Friendship interface {
	Add(userID, targetID, actorID int) error
	Handle(requestID int, req entity.FriendRequestAction) error
	Accept(requestID int) error // Accept.
	Reject(requestID int) error
	Block(requestID int) error
	GetRequested(id int) ([]entity.Friend, error)
	GetPending(id int) ([]entity.Friend, error)
	GetBlock(id int) ([]entity.Friend, error)
	GetFriends(id int) ([]entity.Friend, error)
}
