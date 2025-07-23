package shared

import (
	"fmt"
	"strings"
)

// FormatDestination - 국가-도시를 표준 형식으로 포맷팅
func FormatDestination(country, city string) string {
	if country == "" || city == "" {
		return ""
	}
	country = strings.TrimSpace(country)
	city = strings.TrimSpace(city)
	return fmt.Sprintf("%s-%s", country, city)
}

// ParseDestination - "국가-도시" 문자열을 파싱
func ParseDestination(destination string) (country, city string) {
	if destination == "" {
		return "", ""
	}

	parts := strings.Split(destination, "-")
	if len(parts) != 2 {
		return "", ""
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

// ValidateDestination - 목적지 형식 검증
func ValidateDestination(country, city string) bool {
	return strings.TrimSpace(country) != "" && strings.TrimSpace(city) != ""
}

// NormalizeDestination - 목적지 정규화 (대소문자, 공백 처리)
func NormalizeDestination(country, city string) (string, string) {
	return strings.TrimSpace(country), strings.TrimSpace(city)
}
