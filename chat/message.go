package chat

// Message represents the data that is passed through the websocket
type Message struct {
	Type string `json:"type"`
	Room string `json:"room"`
	Data string `json:"data"`
}
