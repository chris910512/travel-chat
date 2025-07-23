package repository

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/user"
)

type UserRepository interface {
	Create(user *user.User) error
	GetByID(id uint) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Update(user *user.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*user.User, error)

	GetByDestination(country, city string) ([]*user.User, error)
	GetActiveUsers() ([]*user.User, error)
	UpdateLastActive(userID uint) error
}
