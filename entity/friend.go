package entity

type FriendRequestAction string

var (
	AcceptFriend = FriendRequestAction("accept")
	BlockFriend  = FriendRequestAction("block")
	RejectFriend = FriendRequestAction("reject")
)

type Friend struct {
	ID string
}

type FilterFriendOption string

var (
	FilterFriends   = FilterFriendOption("friends")
	FilterRequested = FilterFriendOption("requested")
	FilterPending   = FilterFriendOption("pending")
	FilterBlocked   = FilterFriendOption("blocked")
)
