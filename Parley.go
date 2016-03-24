package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Connection struct {
	username string
	socket   *websocket.Conn
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections []Connection

func main() {

	port := flag.Int("p", 3000, "server port 3000")
	flag.Parse()

	// check if port number is between 1 and 65535
	if *port <= 0 || *port > 1<<16-1 {
		log.Fatal(fmt.Sprintf("port %v number invalid", *port))
	}

	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws/:uname", func(c *gin.Context) {
		username := c.Param("uname")
		openws(c.Writer, c.Request, username)
	})

	r.Run(fmt.Sprintf("localhost:%v", *port))
}

func openws(responseWriter http.ResponseWriter, request *http.Request, username string) {

	socket, err := wsupgrader.Upgrade(responseWriter, request, nil)
	connections = append(connections, Connection{username, socket})

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

		msg = buildMsg(msg, username)

		for _, connection := range connections {
			connection.socket.WriteMessage(t, msg)
		}
	}
}

func buildMsg(msg []byte, username string) []byte {
	now := time.Now().Format("15:04:05")
	now =  fmt.Sprintf("%s", now)
	return append([]byte(now + " | " + username + ": "), msg...);
}
