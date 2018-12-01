package service

import (
	"context"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type GetConversationsRequest struct {
	RoomID string
}

type GetConversationsResponse struct {
	Data []entity.Conversation `json:"data"`
	Room string                `json:"room"`
}
type GetConversations func(ctx context.Context, req GetConversationsRequest) (*GetConversationsResponse, error)

func NewGetConversationsService(repo repository.Conversation) GetConversations {
	return func(ctx context.Context, req GetConversationsRequest) (*GetConversationsResponse, error) {
		conversations, err := repo.GetConversations(req.RoomID)
		if err != nil {
			return nil, err
		}
		return &GetConversationsResponse{
			Data: conversations,
			Room: req.RoomID,
		}, nil
	}
}
