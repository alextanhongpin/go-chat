package service

import (
	"context"
	"errors"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/pkg/token"
	"github.com/alextanhongpin/go-chat/repository"
)

type AuthorizeRequest struct {
	// UserID string `json:"user_id"`
}

type AuthorizeResponse struct {
	// Token string `json:"token"`
	Name string `json:"name"`
}

type Authorize func(ctx context.Context, req AuthorizeRequest) (*AuthorizeResponse, error)

func NewAuthorizeService(repo repository.User) Authorize {
	return func(ctx context.Context, req AuthorizeRequest) (*AuthorizeResponse, error) {
		userID := ctx.Value(entity.ContextKeyUserID).(string)
		user, err := repo.GetUser(userID)
		if err != nil {
			return nil, err
		}
		return &AuthorizeResponse{
			Name: user.Name,
		}, nil
	}
}

type RegisterRequest struct {
	Email        string `json:"email"`
	ConfirmEmail string `json:"confirm_email"`
	Password     string `json:"password"`
	Name         string `json:"name"`
}

func (r *RegisterRequest) Validate() error {
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

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type Register func(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)

func NewRegisterService(repo repository.User, signer token.Signer) Register {
	return func(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
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
		return &RegisterResponse{
			AccessToken: token,
		}, nil
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Validate() error {
	if l.Email == "" {
		return errors.New("email is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type Login func(ctx context.Context, req LoginRequest) (*LoginResponse, error)

func NewLoginService(repo repository.User, signer token.Signer) Login {
	return func(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
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
		return &LoginResponse{
			AccessToken: accessToken,
			ExpiresIn:   signer.ExpiresIn(),
		}, nil
	}
}
