package controller

import (
	"context"

	"github.com/alextanhongpin/go-chat/repository"
	"github.com/alextanhongpin/go-chat/ticket"
)

type postAuthRequest struct {
	UserID string `json:"user_id"`
}

type postAuthResponse struct {
	Token string `json:"token"`
}

type postAuthorizeService func(ctx context.Context, req postAuthRequest) (*postAuthResponse, error)

func MakePostAuthorizeService(repo repository.User, signer ticket.Dispenser) postAuthorizeService {
	return func(ctx context.Context, req postAuthRequest) (*postAuthResponse, error) {
		user, err := repo.GetUserByName(req.UserID)
		if err != nil {
			return nil, err
		}

		// Create new ticket.
		ticket := signer.New(user.ID)

		// Sign ticket.
		token, err := signer.Sign(ticket)
		if err != nil {
			return nil, err
		}
		return &postAuthResponse{
			Token: token,
		}, nil
	}
}
