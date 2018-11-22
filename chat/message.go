package chat

type MessageType string

const (
	MessageTypeText     = MessageType("text")
	MessageTypeStatus   = MessageType("status")
	MessageTypePresence = MessageType("presence")
	MessageTypeOnline   = MessageType("online")
	MessageTypeOffline  = MessageType("offline")
)

func (m MessageType) String() string {
	return string(m)
}
