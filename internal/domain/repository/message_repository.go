package repository

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/message"
	"time"
)

type MessageRepository interface {
	Create(message *message.Message) error
	GetByID(id uint) (*message.Message, error)
	GetByChatRoom(chatRoomID uint, limit int) ([]*message.Message, error)
	DeleteExpired() error
	DeleteExpiredBefore(before time.Time) error
	Count() (int64, error)
}
