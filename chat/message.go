package chat

// Message represents the data that is passed through the websocket
type Message struct {
	Type    string   `json:"type,omitempty"`
	Room    string   `json:"room,omitempty"`
	Data    string   `json:"data,omitempty"`
	Token   string   `json:"token,omitempty"`
	History []string `json:"history,omitempty"`
	// IP
	// Device (mobile, web etc)
	// UserAgent
	// ID
	// Profile Photo
	// Name
	// Location (Lat, Lng)
	// Country
}

// type History struct {
// 	Message   string    `json:"message"`
// 	CreatedAt time.Time `json:"created_at"`
// }
