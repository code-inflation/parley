package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Message struct {
	Username string
	Text     string
	Time     time.Time
}

func (msg Message) BuildJson() []byte {

	timeString := fmt.Sprintf("%s", msg.Time.Format("15:04:05"))
	msgMap := map[string]string{"username": msg.Username, "text": msg.Text, "time": timeString}
	msgBytes, err := json.Marshal(msgMap)
	if err != nil {
		log.Print(err)
	}
	return msgBytes
}
