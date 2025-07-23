package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Age         int            `json:"age"`
	Gender      string         `json:"gender"`
	ProfilePic  string         `json:"profile_pic"`
	Destination string         `json:"destination"`
	TravelStart time.Time      `json:"travel_start"`
	TravelEnd   time.Time      `json:"travel_end"`
	Bio         string         `json:"bio"`
	LastActive  time.Time      `json:"last_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
