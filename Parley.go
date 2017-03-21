package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/code-inflation/parley/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
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

var dblog *core.Dblogger = core.NewDblogger()

func main() {
	port := parsePort()
	setUpRouter(port)
}

func parsePort() int {
	port := flag.Int("p", 3000, "server port 3000")
	flag.Parse()

	// check if port number is between 1 and 65535
	if *port <= 0 || *port > 1<<16-1 {
		log.Fatal(fmt.Sprintf("port %v number invalid", *port))
	}

	return *port
}

func setUpRouter(port int) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"msgs": dblog.FindAllMsg(),
		})
	})

	r.GET("/ws/:uname", func(c *gin.Context) {
		username := c.Param("uname")
		openws(c.Writer, c.Request, username)
	})

	r.Run(fmt.Sprintf("localhost:%v", port))
}

func openws(responseWriter http.ResponseWriter, request *http.Request, username string) {

	socket, err := wsupgrader.Upgrade(responseWriter, request, nil)
	connections = append(connections, Connection{username, socket})

	if err != nil {
		log.Print(err)
		return
	}

	for {
		t, inputBytes, err := socket.ReadMessage()
		if err != nil {
			log.Print(err)
			break
		}

		msg := core.Message{Username: username, Text: string(inputBytes), Time: time.Now()}
		jsonBytes := msg.BuildJSON()
		go dblog.SaveMsg(msg)

		for _, connection := range connections {
			go connection.socket.WriteMessage(t, jsonBytes)
		}
	}
}
