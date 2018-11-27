package controller

import (
	"context"
	"errors"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type getRoomsRequest struct {
}

type getRoomsResponse struct {
	Data []entity.UserRoom `json:"data"`
}

type getRoomsService func(ctx context.Context, req getRoomsRequest) (*getRoomsResponse, error)

func MakeGetRoomsService(repo repository.Room) getRoomsService {
	return func(ctx context.Context, req getRoomsRequest) (*getRoomsResponse, error) {
		userID, _ := ctx.Value(entity.ContextKeyUserID).(string)
		if userID == "" {
			return nil, errors.New("user_id is required")
		}

		rooms, err := repo.GetRooms(userID)
		if err != nil {
			return nil, err
		}
		return &getRoomsResponse{rooms}, nil
	}
}
