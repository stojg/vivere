// Copyright 2013 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"log"
	"net/http"
	"time"
	"fmt"
	"encoding/json"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if r.URL.Path[1:] == "" {
		http.ServeFile(w, r, "static/index.html")
		return
	}
	http.ServeFile(w, r, "static/"+r.URL.Path[1:])
}

type Message struct {
	Name string
	Rotation float32
	Time int64
}

func main() {
	go h.run()
	m := Message{"bunny", 0, time.Now().UnixNano()}

	go func() {
		timer := time.Tick(16 * time.Millisecond)
		for now := range timer {
			m.Rotation = m.Rotation + 0.01
			m.Time = now.UnixNano()
			b, _ := json.Marshal(m)
			h.broadcast <- b
			fmt.Printf("%s\n",b)
		}
	}()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}



}
