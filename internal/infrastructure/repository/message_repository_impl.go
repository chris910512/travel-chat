package repository

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/message"
	"github.com/chris910512/travel-chat/internal/domain/repository"
	"gorm.io/gorm"
	"time"
)

type messageRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepositoryImpl{
		db: db,
	}
}

func (r *messageRepositoryImpl) Create(message *message.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepositoryImpl) GetByID(id uint) (*message.Message, error) {
	var msg message.Message
	err := r.db.First(&msg, id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepositoryImpl) GetByChatRoom(chatRoomID uint, limit int) ([]*message.Message, error) {
	var messages []*message.Message
	err := r.db.Where("chat_room_id = ?", chatRoomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *messageRepositoryImpl) DeleteExpired() error {
	now := time.Now()
	return r.db.Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Delete(&message.Message{}).Error
}

func (r *messageRepositoryImpl) DeleteExpiredBefore(before time.Time) error {
	return r.db.Where("expires_at IS NOT NULL AND expires_at < ?", before).
		Delete(&message.Message{}).Error
}

func (r *messageRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Model(&message.Message{}).Count(&count).Error
	return count, err
}
