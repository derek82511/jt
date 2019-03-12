package web

import (
	"derek82511/jt/service/web/socket"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func SetupWebsocket(app *iris.Application) {
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	ws.OnConnection(func(c websocket.Connection) {
		c.On("register", func(jobID string) {
			socket.SetConnection(jobID, &c)

			c.Emit("onRegister", "ok")
		})
	})

	app.Get("/echo", ws.Handler())

	app.Any("/iris-ws.js", websocket.ClientHandler())
}
