package service

import (
	"context"
	"errors"

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
