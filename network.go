package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/stojg/vivere/lib/client"
	. "github.com/stojg/vivere/lib/components"
	. "github.com/stojg/vector"
	"golang.org/x/net/websocket"
	"net/http"
)

var clients []*client.Client

func init() {
	Println("Inititalising Network")

	ch := client.NewClientManager()
	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "static/index.html")
			return
		}
		http.ServeFile(w, r, "static/"+r.URL.Path[1:])
	})

	go func() {
		Println(http.ListenAndServe(":8080", nil))
	}()

	go func(client chan *client.Client) {
		for {
			select {
			case newClient := <-client:
				Println("New client connected")
				clients = append(clients, newClient)
				//client.Update(world.Serialize(true))
				//world.players = append(world.players, client)
			}
		}
	}(ch.NewClients())
}

func binaryStream(buf *bytes.Buffer, lit Literal, val interface{}) {
	binary.Write(buf, binary.LittleEndian, lit)
	switch val.(type) {
	case uint8:
		binary.Write(buf, binary.LittleEndian, byte(val.(uint8)))
	case uint16:
		binary.Write(buf, binary.LittleEndian, float32(val.(uint16)))
	case EntityType:
		binary.Write(buf, binary.LittleEndian, float32(val.(EntityType)))
	case float32:
		binary.Write(buf, binary.LittleEndian, float32(val.(float32)))
	case float64:
		binary.Write(buf, binary.LittleEndian, float32(val.(float64)))
	case Entity:
		binary.Write(buf, binary.LittleEndian, float32(val.(Entity)))
	case *Vector3:
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[0]))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[1]))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[2]))
	case *Quaternion:
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).R))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).I))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).J))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).K))
	default:
		panic(fmt.Errorf("Havent found out how to serialise literal %v with value of type '%T'", lit, val))
	}

}
