package service

import (
	"context"
	"errors"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type AddFriendRequest struct {
	UserID   int `json:"-"`
	TargetID int `json:"-"`
}

func (a *AddFriendRequest) Validate() error {
	if a.UserID == a.TargetID {
		return errors.New("id cannot be the same")
	}
	return nil
}

func (a *AddFriendRequest) Sort() (l, r int) {
	if a.UserID > a.TargetID {
		return a.TargetID, a.UserID
	}
	return a.UserID, a.TargetID
}

type AddFriendResponse struct {
	Status bool `json:"status"`
}

type AddFriend func(ctx context.Context, req AddFriendRequest) (*AddFriendResponse, error)

func NewAddFriendService(repo repository.Friendship) AddFriend {
	return func(ctx context.Context, req AddFriendRequest) (*AddFriendResponse, error) {
		if err := req.Validate(); err != nil {
			return nil, err
		}

		l, r := req.Sort()
		if err := repo.AddFriend(l, r, req.UserID); err != nil {
			return nil, err
		}

		return &AddFriendResponse{Status: true}, nil
	}
}

type HandleFriendRequest struct {
	// RequestID int
	UserID   int                        `json:"-"`
	TargetID int                        `json:"-"`
	Action   entity.FriendRequestAction `json:"action,omitempty"`
}

func (h *HandleFriendRequest) Sort() (int, int) {
	if h.UserID > h.TargetID {
		return h.TargetID, h.UserID
	}
	return h.UserID, h.TargetID
}

type HandleFriendResponse struct {
	Status bool `json:"status"`
}

type HandleFriend func(ctx context.Context, req HandleFriendRequest) (*HandleFriendResponse, error)

func NewHandleFriendService(repo repository.Friendship) HandleFriend {
	return func(ctx context.Context, req HandleFriendRequest) (*HandleFriendResponse, error) {
		var err error
		l, r := req.Sort()
		switch req.Action {
		case entity.AcceptFriend:
			err = repo.AcceptFriend(l, r)
		case entity.BlockFriend:
			err = repo.BlockFriend(l, r)
		case entity.RejectFriend:
			err = repo.RejectFriend(l, r)
		}
		return &HandleFriendResponse{Status: true}, err
	}
}

type ListFriendRequest struct {
	Filter entity.FilterFriendOption
	UserID int
}

type ListFriendResponse struct {
	Friends []entity.Friend
}

type ListFriend func(ctx context.Context, req ListFriendRequest) (*ListFriendResponse, error)

func NewListFriend(repo repository.Friendship) ListFriend {
	return func(ctx context.Context, req ListFriendRequest) (*ListFriendResponse, error) {
		var res []entity.Friend
		var err error
		switch req.Filter {
		case entity.FilterFriends:
			res, err = repo.GetMutualFriends(req.UserID)
		case entity.FilterRequested:
			res, err = repo.GetRequestedFriends(req.UserID)
		case entity.FilterPending:
			res, err = repo.GetPendingFriends(req.UserID)
		case entity.FilterBlocked:
			res, err = repo.GetBlockedFriends(req.UserID)
		}
		return &ListFriendResponse{
			Friends: res,
		}, err
	}
}
