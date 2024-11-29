package realtime

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NotesSyncHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	for {
		var message string
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		err = conn.WriteJSON(message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
