package service

import (
	"context"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/repository"
)

type GetPostsRequest struct{}
type GetPostsResponse struct {
	Data []entity.Post `json:"data"`
}

type GetPosts func(context.Context, GetPostsRequest) (*GetPostsResponse, error)

func NewGetPostsService(repo repository.Post) GetPosts {
	return func(ctx context.Context, req GetPostsRequest) (*GetPostsResponse, error) {
		res, err := repo.GetPosts()
		return &GetPostsResponse{
			Data: res,
		}, err
	}
}

type GetPostRequest struct {
	ID string
}
type GetPostResponse struct {
	Data entity.Post `json:"data"`
}

type GetPost func(context.Context, GetPostRequest) (*GetPostResponse, error)

func NewGetPostService(repo repository.Post) GetPost {
	return func(ctx context.Context, req GetPostRequest) (*GetPostResponse, error) {
		res, err := repo.GetPost(req.ID)
		return &GetPostResponse{
			Data: res,
		}, err
	}
}

type CreatePostRequest struct {
	Text   string
	UserID string
}

type CreatePostResponse struct {
	Success bool
}

type CreatePost func(context.Context, CreatePostRequest) (*CreatePostResponse, error)

func NewCreatePostService(repo repository.Post) CreatePost {
	return func(ctx context.Context, req CreatePostRequest) (*CreatePostResponse, error) {
		_, err := repo.CreatePost(entity.Post{
			UserID: req.UserID,
			Text:   req.Text,
		})
		return &CreatePostResponse{
			Success: true,
		}, err
	}
}

type UpdatePostRequest struct {
	Text string
}

type UpdatePostResponse struct {
	Success bool
}

type UpdatePost func(context.Context, UpdatePostRequest) (*UpdatePostResponse, error)

func NewUpdatePostService(repo repository.Post) UpdatePost {
	return func(ctx context.Context, req UpdatePostRequest) (*UpdatePostResponse, error) {
		err := repo.UpdatePost(entity.Post{
			Text: req.Text,
		})
		return &UpdatePostResponse{
			Success: true,
		}, err
	}
}

type DeletePostRequest struct {
	ID string
}
type DeletePostResponse struct {
	Success bool
}
type DeletePost func(context.Context, DeletePostRequest) (*DeletePostResponse, error)

func NewDeletePostService(repo repository.Post) DeletePost {
	return func(ctx context.Context, req DeletePostRequest) (*DeletePostResponse, error) {
		err := repo.DeletePost(req.ID)
		return &DeletePostResponse{
			Success: true,
		}, err
	}
}
