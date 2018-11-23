package chat

type MessageType string

const (
	MessageTypeText     = MessageType("text")
	MessageTypeStatus   = MessageType("status")
	MessageTypePresence = MessageType("presence")
	MessageTypeOnline   = MessageType("online")
	MessageTypeOffline  = MessageType("offline")
	MessageTypeAuth     = "auth"
	MessageTypeMessage  = "message"
)

func (m MessageType) String() string {
	return string(m)
}

type Message struct {
	// The text content of the message.
	Text string `json:"data"`

	// Unique ID for the message created.
	ID string `json:"id"`

	// The type for client to handle polymorphism.
	Type string `json:"type"`

	Timestamp int64 `json:"timestamp"`
	// User might not persist timestamp that long (when messages are loaded from db).
	// k = HMAC(timestamp,id,msg, sk)
	// timestamp|id|(msg|sender)k|HMAC(timestamp|id|sender|msg, sk))
	// Hash string

	// The sender of the message. Usually a user id.
	Sender string `json:"sender"`

	// The receiver in this case is the room id (chat group), rather than
	// the user id. This provides better flexibility, as it allows us to
	// send conversations to more than a person (users in the same group,
	// rather than just direct message to a single user).
	Receiver string `json:"room"`
}
