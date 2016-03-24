package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Connection struct {
	uname  string
	socket *websocket.Conn
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections []Connection

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws/:uname", func(c *gin.Context) {
		uname := c.Param("uname")
		websocks(c.Writer, c.Request, uname)
	})

	r.Run("localhost:3000")
}

func websocks(responseWriter http.ResponseWriter, request *http.Request, uname string) {

	socket, err := wsupgrader.Upgrade(responseWriter, request, nil)
	connections = append(connections, Connection{uname, socket})

	if err != nil {
		log.Print(err)
		return
	}

	for {
		t, msg, err := socket.ReadMessage()
		if err != nil {
			log.Print(err)
			break
		}

		msg = append([]byte(uname+" says "), msg...)

		for _, connection := range connections {
			connection.socket.WriteMessage(t, msg)
		}
	}
}
