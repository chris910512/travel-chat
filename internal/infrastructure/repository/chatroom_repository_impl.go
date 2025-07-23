package repository

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/chatroom"
	"github.com/chris910512/travel-chat/internal/domain/repository"
	"gorm.io/gorm"
)

type chatRoomRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRoomRepository(db *gorm.DB) repository.ChatRoomRepository {
	return &chatRoomRepositoryImpl{
		db: db,
	}
}

func (r *chatRoomRepositoryImpl) Create(chatRoom *chatroom.ChatRoom) error {
	return r.db.Create(chatRoom).Error
}

func (r *chatRoomRepositoryImpl) GetByID(id uint) (*chatroom.ChatRoom, error) {
	var room chatroom.ChatRoom
	err := r.db.First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRoomRepositoryImpl) GetByLocation(country, city string, roomType chatroom.RoomType) (*chatroom.ChatRoom, error) {
	var room chatroom.ChatRoom
	err := r.db.Where("country = ? AND city = ? AND room_type = ?",
		country, city, roomType).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRoomRepositoryImpl) GetOrCreatePublicRoom(country, city string) (*chatroom.ChatRoom, error) {
	// 먼저 기존 방이 있는지 확인
	room, err := r.GetByLocation(country, city, chatroom.RoomTypePublic)
	if err == nil {
		return room, nil
	}

	// 없으면 새로 생성
	newRoom := &chatroom.ChatRoom{
		Country:  country,
		City:     city,
		RoomType: chatroom.RoomTypePublic,
	}
	newRoom.GeneratePublicRoomName()

	err = r.Create(newRoom)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (r *chatRoomRepositoryImpl) CreatePrivateRoom(country, city, user1Name, user2Name string) (*chatroom.ChatRoom, error) {
	newRoom := &chatroom.ChatRoom{
		Country:  country,
		City:     city,
		RoomType: chatroom.RoomTypePrivate,
	}
	newRoom.GeneratePrivateRoomName(user1Name, user2Name)

	err := r.Create(newRoom)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (r *chatRoomRepositoryImpl) Update(chatRoom *chatroom.ChatRoom) error {
	return r.db.Save(chatRoom).Error
}

func (r *chatRoomRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&chatroom.ChatRoom{}, id).Error
}
