package repository

import "github.com/alextanhongpin/go-chat/entity"

type Post interface {
	GetPosts() ([]entity.Post, error)
	GetPost(id string) (entity.Post, error)
	CreatePost(post entity.Post) (int64, error)
	UpdatePost(post entity.Post) error
	DeletePost(id string) error
}
