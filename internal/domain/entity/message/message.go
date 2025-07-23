package message

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Content     string         `gorm:"not null;type:text" json:"content"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	ChatRoomID  uint           `gorm:"not null" json:"chat_room_id"`
	MessageType MessageType    `gorm:"default:0" json:"message_type"`
	ExpiresAt   *time.Time     `json:"expires_at"` // 메시지 만료 시간 (nullable)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// IsExpired - 메시지가 만료되었는지 확인
func (m *Message) IsExpired() bool {
	if m.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*m.ExpiresAt)
}

// SetExpiration - 메시지 만료 시간 설정
func (m *Message) SetExpiration(duration time.Duration) {
	expiresAt := m.CreatedAt.Add(duration)
	m.ExpiresAt = &expiresAt
}

// SetPublicChatExpiration - 전체 채팅 만료 시간 설정 (6시간)
func (m *Message) SetPublicChatExpiration() {
	m.SetExpiration(6 * time.Hour)
}

// SetPrivateChatExpiration - 1:1 채팅 만료 시간 설정 (24시간)
func (m *Message) SetPrivateChatExpiration() {
	m.SetExpiration(24 * time.Hour)
}
