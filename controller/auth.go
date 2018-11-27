package controller

import (
	"context"
	"errors"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/pkg/token"
	"github.com/alextanhongpin/go-chat/repository"
)

type postAuthRequest struct {
	UserID string `json:"user_id"`
}

type postAuthResponse struct {
	Token string `json:"token"`
}

type postAuthorizeService func(ctx context.Context, req postAuthRequest) (*postAuthResponse, error)

func MakePostAuthorizeService(repo repository.User, signer token.Signer) postAuthorizeService {
	return func(ctx context.Context, req postAuthRequest) (*postAuthResponse, error) {
		user, err := repo.GetUserByName(req.UserID)
		if err != nil {
			return nil, err
		}

		// Sign the user.
		token, err := signer.Sign(user.ID)
		if err != nil {
			return nil, err
		}
		return &postAuthResponse{
			Token: token,
		}, nil
	}
}

type registerRequest struct {
	Email        string `json:"email"`
	ConfirmEmail string `json:"confirm_email"`
	Password     string `json:"password"`
	Name         string `json:"name"`
}

func (r *registerRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.ConfirmEmail == "" {
		return errors.New("confirm_email is required")
	}
	if r.Email != r.ConfirmEmail {
		return errors.New("email does not match")
	}
	// TODO: Validate email format.
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

type registerResponse struct {
	AccessToken string `json:"access_token"`
}
type registerService func(ctx context.Context, req registerRequest) (*registerResponse, error)

func MakeRegisterService(repo repository.User, signer token.Signer) registerService {
	return func(ctx context.Context, req registerRequest) (*registerResponse, error) {
		if err := req.Validate(); err != nil {
			return nil, err
		}
		_, err := repo.GetUserByEmail(req.Email)
		if err != entity.ErrUserNotFound {
			return nil, err
		}
		user := entity.NewUser(req.Name, req.Email)
		if err := user.SetPassword(req.Password); err != nil {
			return nil, err
		}
		if err = repo.CreateUser(user); err != nil {
			return nil, err
		}
		token, err := signer.Sign(user.ID)
		if err != nil {
			return nil, err
		}
		return &registerResponse{
			AccessToken: token,
		}, nil
	}
}
