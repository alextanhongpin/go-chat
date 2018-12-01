package repository

import "github.com/alextanhongpin/go-chat/entity"

type Friendship interface {
	AddFriend(userID, targetID, actorID int) error
	AcceptFriend(requestID int) error // Accept.
	RejectFriend(requestID int) error
	BlockFriend(requestID int) error
	GetRequestedFriends(id int) ([]entity.Friend, error)
	GetPendingFriends(id int) ([]entity.Friend, error)
	GetBlockedFriends(id int) ([]entity.Friend, error)
	GetMutualFriends(id int) ([]entity.Friend, error)
}
