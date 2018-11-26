package repository

import "github.com/alextanhongpin/go-chat/entity"

// Conversation represents the DAO for conversation.
type Conversation interface {
	GetConversations(roomID string) ([]entity.Conversation, error)
}
