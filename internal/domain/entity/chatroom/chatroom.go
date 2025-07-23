package chatroom

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/shared"
	"gorm.io/gorm"
	"time"
)

type ChatRoom struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Country   string         `gorm:"not null;size:100" json:"country"` // 국가
	City      string         `gorm:"not null;size:100" json:"city"`    // 도시
	RoomType  RoomType       `gorm:"not null;default:0" json:"room_type"`
	Name      string         `gorm:"size:200" json:"name"` // 채팅방 이름 (1:1의 경우 자동 생성)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetRoomKey - 채팅방 키 생성 (국가-도시 조합)
func (c *ChatRoom) GetRoomKey() string {
	return shared.FormatDestination(c.Country, c.City)
}

// IsPublic - 전체 채팅방 여부
func (c *ChatRoom) IsPublic() bool {
	return c.RoomType == RoomTypePublic
}

// IsPrivate - 1:1 채팅방 여부
func (c *ChatRoom) IsPrivate() bool {
	return c.RoomType == RoomTypePrivate
}

// GeneratePublicRoomName - 전체 채팅방 이름 생성
func (c *ChatRoom) GeneratePublicRoomName() {
	c.Name = c.Country + " " + c.City + " 여행자 채팅"
}

// GeneratePrivateRoomName - 1:1 채팅방 이름 생성
func (c *ChatRoom) GeneratePrivateRoomName(user1Name, user2Name string) {
	c.Name = user1Name + " & " + user2Name
}
