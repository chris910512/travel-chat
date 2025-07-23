package repository

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/chatroom"
)

type ChatRoomRepository interface {
	Create(chatRoom *chatroom.ChatRoom) error
	GetByID(id uint) (*chatroom.ChatRoom, error)
	GetByLocation(country, city string, roomType chatroom.RoomType) (*chatroom.ChatRoom, error)
	GetOrCreatePublicRoom(country, city string) (*chatroom.ChatRoom, error)
	CreatePrivateRoom(country, city, user1Name, user2Name string) (*chatroom.ChatRoom, error)
	Update(chatRoom *chatroom.ChatRoom) error
	Delete(id uint) error
}
