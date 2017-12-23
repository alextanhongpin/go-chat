package chat

// Message represents the data that is passed through the websocket
type Message struct {
	Type  string `json:"type,omitempty"`
	Room  string `json:"room,omitempty"`
	Data  string `json:"data,omitempty"`
	Token string `json:"token,omitempty"`
}
