package entity

type FriendRequestAction string

var AcceptFriend = FriendRequestAction("accept")
var RejectFriend = FriendRequestAction("reject")
var BlockFriend = FriendRequestAction("block")

type Friend struct {
	ID string
}
