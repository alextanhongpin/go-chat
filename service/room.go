package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type GetRoomsRequest struct {
}

type GetRoomsResponse struct {
	Data []entity.UserRoom `json:"data"`
}

type GetRooms func(ctx context.Context, req GetRoomsRequest) (*GetRoomsResponse, error)

func NewGetRoomsService(repo repository.Room) GetRooms {
	return func(ctx context.Context, req GetRoomsRequest) (*GetRoomsResponse, error) {
		userID, _ := ctx.Value(entity.ContextKeyUserID).(string)
		if userID == "" {
			return nil, errors.New("user_id is required")
		}

		rooms, err := repo.GetRooms(userID)
		if err != nil {
			return nil, err
		}
		return &GetRoomsResponse{rooms}, nil
	}
}

type PostRoomsRequest struct {
	// The user creating the room.
	UserID string `json:"-"`
	// The other user to be added in the room.
	FriendID string `json:"friend_id"`
}

type PostRoomsResponse struct {
	Data entity.UserRoom `json:"data"`
}

type PostRooms func(ctx context.Context, req PostRoomsRequest) (*PostRoomsResponse, error)

func NewPostRoomsService(repo repository.Room) PostRooms {
	return func(ctx context.Context, req PostRoomsRequest) (*PostRoomsResponse, error) {
		roomID, err := repo.CreateRoom(
			req.UserID,
			req.FriendID)
		if err != nil {
			return nil, err
		}
		return &PostRoomsResponse{
			Data: entity.UserRoom{
				UserID: req.FriendID,
				RoomID: strconv.FormatInt(roomID, 10),
			},
		}, nil
	}
}
