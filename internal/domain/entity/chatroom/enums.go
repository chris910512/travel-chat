package chatroom

import "encoding/json"

// RoomType - 채팅방 타입
type RoomType int

const (
	RoomTypePublic  RoomType = iota // 0 - 전체 채팅방
	RoomTypePrivate                 // 1 - 1:1 채팅방
)

func (rt *RoomType) String() string {
	switch *rt {
	case RoomTypePublic:
		return "public"
	case RoomTypePrivate:
		return "private"
	default:
		return "unknown"
	}
}

func (rt *RoomType) IsValid() bool {
	return *rt >= RoomTypePublic && *rt <= RoomTypePrivate
}

func (rt *RoomType) MarshalJSON() ([]byte, error) {
	return json.Marshal(rt.String())
}

func (rt *RoomType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*rt = RoomTypeFromString(s)
	return nil
}

func RoomTypeFromString(s string) RoomType {
	switch s {
	case "public":
		return RoomTypePublic
	case "private":
		return RoomTypePrivate
	default:
		return RoomTypePublic
	}
}
