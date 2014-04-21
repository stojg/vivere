package websocket

import (
	"encoding/json"
	"log"
)

type Messenger interface {
	Message() *Message
}

type Message struct {
	Event   string
	Message interface{}
}

func (m Message) JSON() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		log.Println("Error:", err)
		log.Printf("%v\n\n", m)
	}
	return []byte(json)
}
