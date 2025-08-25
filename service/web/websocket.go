package web

import (
	"derek82511/jt/service/web/socket"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity.
		// In a production environment, you might want to check the origin.
		return true
	},
}

// Message is the structure for websocket communication
type Message struct {
	Event   string `json:"event"`
	Payload string `json:"payload"`
}

func serveWs(ctx iris.Context) {
	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	if err != nil {
		ctx.StopWithError(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	var jobID string

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if jobID != "" {
				socket.RemoveConnection(jobID)
			}
			break // Exit loop on error (e.g., client disconnects)
		}

		var msg Message
		if err := json.Unmarshal(p, &msg); err != nil {
			continue // Ignore malformed messages
		}

		if msg.Event == "register" {
			jobID = msg.Payload
			socket.SetConnection(jobID, conn)
			// Acknowledge registration
			response := Message{Event: "onRegister", Payload: "ok"}
			respBytes, _ := json.Marshal(response)
			conn.WriteMessage(websocket.TextMessage, respBytes)
		}
	}
}

func SetupWebsocket(app *iris.Application) {
	app.Get("/ws", serveWs)
}
