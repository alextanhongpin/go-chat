package service

import (
	"context"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type GetUsersRequest struct {
	UserID int
}

type GetUsersResponse struct {
	Data []entity.Friend `json:"data,omitempty"`
}

type GetUsers func(context.Context, GetUsersRequest) (*GetUsersResponse, error)

func NewGetUsersService(repo repository.User) GetUsers {
	return func(ctx context.Context, req GetUsersRequest) (*GetUsersResponse, error) {
		users, err := repo.GetUsers(req.UserID)
		if err != nil {
			return nil, err
		}

		return &GetUsersResponse{users}, nil
	}
}
