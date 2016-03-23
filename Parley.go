package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		websocks(c.Writer, c.Request)
	})

	r.Run("localhost:3000")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocks(responseWriter http.ResponseWriter, request *http.Request) {

	connection, err := wsupgrader.Upgrade(responseWriter, request, nil)

	if err != nil {
		// TODO log this
		return
	}

	for {
		t, msg, err := connection.ReadMessage()
		if err != nil {
			break
		}
		connection.WriteMessage(t, msg)
	}
}
