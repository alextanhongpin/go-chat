package repository

import "github.com/alextanhongpin/go-chat/entity"

type Conversation interface {
	GetConversations(roomID string) ([]entity.Conversation, error)
}
