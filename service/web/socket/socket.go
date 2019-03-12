package socket

import "github.com/kataras/iris/websocket"

var wsJobConnectionMap *map[string]*websocket.Connection

func init() {
	wsJobConnectionMap = &map[string]*websocket.Connection{}
}

func SetConnection(jobID string, conn *websocket.Connection) {
	ReleaseConnection(jobID)
	(*wsJobConnectionMap)[jobID] = conn
}

func GetConnection(jobID string) *websocket.Connection {
	return (*wsJobConnectionMap)[jobID]
}

func ReleaseConnection(jobID string) {
	conn := GetConnection(jobID)
	if conn != nil {
		(*conn).Disconnect()
		delete((*wsJobConnectionMap), jobID)
	}
}
