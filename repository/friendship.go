package repository

import "github.com/alextanhongpin/go-chat/entity"

type Friendship interface {
	AddFriend(userID, targetID, actorID int) error
	AcceptFriend(a, b int) error // Accept.
	RejectFriend(a, b int) error
	BlockFriend(a, b int) error
	GetRequestedFriends(id int) ([]entity.Friend, error)
	GetPendingFriends(id int) ([]entity.Friend, error)
	GetBlockedFriends(id int) ([]entity.Friend, error)
	GetMutualFriends(id int) ([]entity.Friend, error)
}
