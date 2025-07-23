package errors

import "errors"

// 사용자 관련 에러들
var (
	ErrUserNotFound       = errors.New("사용자를 찾을 수 없습니다")
	ErrEmailAlreadyExists = errors.New("이미 사용 중인 이메일입니다")
	ErrInvalidCredentials = errors.New("이메일 또는 비밀번호가 올바르지 않습니다")
	ErrWeakPassword       = errors.New("비밀번호는 최소 6자 이상이어야 합니다")
	ErrInvalidEmail       = errors.New("올바르지 않은 이메일 형식입니다")
	ErrInvalidTravelDates = errors.New("여행 시작일은 종료일보다 빨라야 합니다")
	ErrPastTravelDate     = errors.New("여행 시작일은 현재 날짜 이후여야 합니다")
	ErrUnauthorized       = errors.New("인증이 필요합니다")
	ErrForbidden          = errors.New("접근 권한이 없습니다")
)

// 에러 타입 체크 헬퍼 함수들
func IsUserNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}

func IsEmailAlreadyExists(err error) bool {
	return errors.Is(err, ErrEmailAlreadyExists)
}

func IsInvalidCredentials(err error) bool {
	return errors.Is(err, ErrInvalidCredentials)
}

func IsWeakPassword(err error) bool {
	return errors.Is(err, ErrWeakPassword)
}

func IsInvalidEmail(err error) bool {
	return errors.Is(err, ErrInvalidEmail)
}

func IsInvalidTravelDates(err error) bool {
	return errors.Is(err, ErrInvalidTravelDates)
}

func IsPastTravelDate(err error) bool {
	return errors.Is(err, ErrPastTravelDate)
}

func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}
