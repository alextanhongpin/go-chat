package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type message struct {
	Handle string `json:"handle"`
	Text   string `json:"text"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("listening to port *:8080. press ctrl + c to cancel.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func validateMessage(data []byte) (message, error) {
	var msg message
	if err := json.Unmarshal(data, &msg); err != nil {
		return msg, err
	}
	if msg.Handle == "" && msg.Text == "" {
		return msg, errors.New("Message has no handle or text")
	}
	return msg, nil
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for {
		mt, data, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				log.Println("Websocket closed")
				break
			}
			log.Println("Error reading websocket message")
		}

		switch mt {
		case websocket.TextMessage:
			msg, err := validateMessage(data)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("got message:", msg)

			// switch msg.Type {
			// case  "authenticate":
			// 	// jwt authenticate
			// }
		default:
			log.Println("Unknown message")
		}
	}

	ws.WriteMessage(websocket.CloseMessage, []byte{})
}
