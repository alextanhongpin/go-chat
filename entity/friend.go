package entity

type FriendRequestAction string

var (
	AcceptFriend = FriendRequestAction("accept")
	BlockFriend  = FriendRequestAction("block")
	RejectFriend = FriendRequestAction("reject")
)

type Friend struct {
	ID          string `json:"id"`
	Status      string `json:"status,omitempty"`
	Name        string `json:"name,omitempty"`
	IsRequested bool   `json:"is_requested"`
}

type FilterFriendOption string

var (
	FilterFriends   = FilterFriendOption("friends")
	FilterRequested = FilterFriendOption("requested")
	FilterPending   = FilterFriendOption("pending")
	FilterBlocked   = FilterFriendOption("blocked")
)
