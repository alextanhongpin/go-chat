package service

import (
	"context"
	"errors"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type AddFriendRequest struct {
	UserID   int
	TargetID int
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
	Status bool
}

type AddFriendService func(ctx context.Context, req AddFriendRequest) (*AddFriendResponse, error)

func MakeAddFriendService(repo repository.Friendship) AddFriendService {
	return func(ctx context.Context, req AddFriendRequest) (*AddFriendResponse, error) {
		if err := req.Validate(); err != nil {
			return nil, err
		}

		l, r := req.Sort()
		if err := repo.Add(l, r, l); err != nil {
			return nil, err
		}

		return &AddFriendResponse{Status: true}, nil
	}
}

type HandleFriendRequest struct {
	RequestID int
	UserID    int
	Action    entity.FriendRequestAction
}

type HandleFriendResponse struct {
	Status bool
}

type HandleFriendService func(ctx context.Context, req HandleFriendRequest) (*HandleFriendResponse, error)

func MakeHandleFriendService(repo repository.Friendship) HandleFriendService {
	return func(ctx context.Context, req HandleFriendRequest) (*HandleFriendResponse, error) {
		var err error
		switch req.Action {
		case entity.AcceptFriend:
			err = repo.Accept(req.RequestID)
		case entity.BlockFriend:
			err = repo.Block(req.RequestID)
		case entity.RejectFriend:
			err = repo.Reject(req.RequestID)
		}
		return nil, err
	}
}

//
// type ListFriendRequest struct {
//         Filter entity.FilterFriendOption
//         UserID string
// }
//
// type ListFriendResponse struct {
//         Friends []entity.Friend
// }
//
// type ListFriendService func(ctx context.Context, req ListFriendRequest) (*ListFriendResponse, error)
//
// func MakeListFriendService(repo repository.Friendship) ListFriendService {
//         return func(ctx context.Context, req ListFriendRequest) (*ListFriendResponse, error) {
//
//                 switch req.Filter {
//                 case entity.FilterFriends:
//                 case entity.FilterRequested:
//                 case entity.FilterPending:
//                 case entity.FilterBlocked:
//                 }
//                 return nil, nil
//         }
// }
