package repository

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/user"
	"github.com/chris910512/travel-chat/internal/domain/repository"
	"gorm.io/gorm"
	"time"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Create(user *user.User) error {
	user.LastActive = time.Now()
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) GetByID(id uint) (*user.User, error) {
	var u user.User
	err := r.db.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepositoryImpl) GetByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepositoryImpl) Update(user *user.User) error {
	return r.db.Save(user).Error
}

func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&user.User{}, id).Error
}

func (r *userRepositoryImpl) List(offset, limit int) ([]*user.User, error) {
	var users []*user.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) GetByDestination(country, city string) ([]*user.User, error) {
	var users []*user.User
	err := r.db.Where("country = ? AND city = ?", country, city).Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) GetActiveUsers() ([]*user.User, error) {
	var users []*user.User
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	err := r.db.Where("last_active > ?", tenMinutesAgo).Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) UpdateLastActive(userID uint) error {
	return r.db.Model(&user.User{}).Where("id = ?", userID).
		Update("last_active", time.Now()).Error
}

func (r *userRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Model(&user.User{}).Count(&count).Error
	return count, err
}
