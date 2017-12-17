package chat

// Message represents the data that is passed through the websocket
type Message struct {
	Handle string `json:"handle"`
	Text   string `json:"text"`
	Room   string `json:"room"`
}
