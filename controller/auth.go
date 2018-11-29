package controller

import (
	"context"
	"errors"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/pkg/token"
	"github.com/alextanhongpin/go-chat/repository"
)

type postAuthRequest struct {
	// UserID string `json:"user_id"`
}

type postAuthResponse struct {
	// Token string `json:"token"`
	Name string `json:"name"`
}

type postAuthorizeService func(ctx context.Context, req postAuthRequest) (*postAuthResponse, error)

func MakePostAuthorizeService(repo repository.User) postAuthorizeService {
	return func(ctx context.Context, req postAuthRequest) (*postAuthResponse, error) {
		userID := ctx.Value(entity.ContextKeyUserID).(string)
		user, err := repo.GetUser(userID)
		if err != nil {
			return nil, err
		}
		return &postAuthResponse{
			Name: user.Name,
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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *loginRequest) Validate() error {
	if l.Email == "" {
		return errors.New("email is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type loginService func(ctx context.Context, req loginRequest) (*loginResponse, error)

func MakeLoginService(repo repository.User, signer token.Signer) loginService {
	return func(ctx context.Context, req loginRequest) (*loginResponse, error) {
		if err := req.Validate(); err != nil {
			return nil, err
		}
		user, err := repo.GetUserByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if err := user.ComparePassword(req.Password); err != nil {
			return nil, err
		}
		accessToken, err := signer.Sign(user.ID)
		if err != nil {
			return nil, err
		}
		return &loginResponse{
			AccessToken: accessToken,
			ExpiresIn:   signer.ExpiresIn(),
		}, nil
	}
}
