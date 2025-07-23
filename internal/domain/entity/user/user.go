package user

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/shared"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"not null" json:"-"` // JSON에서 제외
	Name          string         `gorm:"not null" json:"name"`
	Age           int            `json:"age"`
	Gender        Gender         `gorm:"default:0" json:"gender"`
	ProfilePic    string         `json:"profile_pic"`
	Country       string         `json:"country"` // 여행 국가
	City          string         `json:"city"`    // 여행 도시
	TravelStart   time.Time      `json:"travel_start"`
	TravelEnd     time.Time      `json:"travel_end"`
	Bio           string         `gorm:"type:text" json:"bio"`
	TravelPurpose TravelPurpose  `gorm:"default:0" json:"travel_purpose"`
	TravelBudget  int            `json:"travel_budget"` // 여행 예산 (만원 단위)
	TravelStyle   TravelStyle    `gorm:"default:0" json:"travel_style"`
	LastActive    time.Time      `gorm:"autoUpdateTime" json:"last_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetDestination - Key for ChatRoom of public chatroom type
func (u *User) GetDestination() string {
	return shared.FormatDestination(u.Country, u.City)
}

// IsActiveRecently - 최근 활동 여부 확인 (10분 이내)
func (u *User) IsActiveRecently() bool {
	return time.Since(u.LastActive) < 10*time.Minute
}

// GetActivityStatus - 활동 상태 문자열 반환
func (u *User) GetActivityStatus() string {
	duration := time.Since(u.LastActive)

	if duration < 10*time.Minute {
		return "온라인"
	} else if duration < time.Hour {
		return "10분 전 활동"
	} else if duration < 24*time.Hour {
		return "1시간 전 활동"
	} else if duration < 7*24*time.Hour {
		return "1일 전 활동"
	} else {
		return "7일 이상 전 활동"
	}
}

func (u *User) UpdateLastActive() {
	u.LastActive = time.Now()
}
