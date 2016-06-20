package main

import (
	"math/rand"
	"time"
)

const (
	SEC_PER_UPDATE float64 = 0.016
)

type Updatable interface {
	Update(elapsed float64)
}

var (
	Frame uint
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	level := NewLevel()

	var previous time.Time = time.Now()
	var lag float64 = 0

	Println("Starting the game loop")
	DebugFPS(SEC_PER_UPDATE)

	for {
		Frame += 1
		now := time.Now()
		elapsed := now.Sub(previous).Seconds()
		previous = now
		lag += elapsed

		level.Update(elapsed)

		buf := level.Draw()
		if buf.Len() > 0 {
			for _, client := range clients {
				go client.Update(buf)
			}
		}
		lag -= SEC_PER_UPDATE
		// save some CPU cycles by sleeping for a while
		time.Sleep(time.Duration((SEC_PER_UPDATE-lag)*1000) * time.Millisecond)
	}
}
