package controller

import (
	"context"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type getConversationsRequest struct {
	RoomID string
}

type getConversationsResponse struct {
	Data []entity.Conversation `json:"data"`
	Room string                `json:"room"`
}
type getConversationsService func(ctx context.Context, req getConversationsRequest) (*getConversationsResponse, error)

func MakeGetConversationsService(repo repository.Conversation) getConversationsService {
	return func(ctx context.Context, req getConversationsRequest) (*getConversationsResponse, error) {
		conversations, err := repo.GetConversations(req.RoomID)
		if err != nil {
			return nil, err
		}
		return &getConversationsResponse{
			Data: conversations,
			Room: req.RoomID,
		}, nil
	}
}
