package message

import "encoding/json"

// MessageType - 메시지 타입
type MessageType int

const (
	MessageTypeText MessageType = iota
	MessageTypeImage
	MessageTypeSystem // 시스템 메시지 (입장/퇴장 등)
)

func (mt *MessageType) String() string {
	switch *mt {
	case MessageTypeText:
		return "text"
	case MessageTypeImage:
		return "image"
	case MessageTypeSystem:
		return "system"
	default:
		return "text"
	}
}

func (mt *MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

func (mt *MessageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*mt = MessageTypeFromString(s)
	return nil
}

func MessageTypeFromString(s string) MessageType {
	switch s {
	case "text":
		return MessageTypeText
	case "image":
		return MessageTypeImage
	case "system":
		return MessageTypeSystem
	default:
		return MessageTypeText
	}
}
