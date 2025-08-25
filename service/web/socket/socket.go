package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	connections = make(map[string]*websocket.Conn)
	mu          sync.Mutex
)

func SetConnection(id string, c *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	// Close any existing connection for this ID
	if conn, ok := connections[id]; ok {
		conn.Close()
	}
	connections[id] = c
}

func GetConnection(id string) *websocket.Conn {
	mu.Lock()
	defer mu.Unlock()
	return connections[id]
}

func RemoveConnection(id string) {
	mu.Lock()
	defer mu.Unlock()
	if conn, ok := connections[id]; ok {
		conn.Close()
		delete(connections, id)
	}
}
