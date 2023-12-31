package models

type MessageType int

const (
	CreatedMessage MessageType = 0
	RegularMessage MessageType = 1

	JoinedMessage MessageType = 2
	LeftMessage   MessageType = 3
	KickMessage   MessageType = 4
)

type Message struct {
	UserId    int64       `json:"user_id"`
	Type      MessageType `json:"message_type"`
	Contents  string      `json:"contents"`
	Timestamp string      `json:"timestamp"`
}
