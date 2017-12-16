package chat

type Message struct {
	Handle string `json:"handle"`
	Text   string `json:"text"`
	Room   string `json:"room"`
}
