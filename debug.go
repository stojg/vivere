package main

import (
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
}

func PrintFPS(world *World) {
	go func() {
		timer := 1 * time.Second
		prev := world.Frame
		prevTime := time.Now()
		for {
			currentTime := <-time.After(timer)
			fps := float64(world.Frame-prev) / float64(currentTime.Sub(prevTime).Seconds())
			if fps < 30 {
				Printf("fps: %0.1f frame %d\n", fps, world.Frame)
			} else {
				dPrintf("fps: %0.1f frame %d\n", fps, world.Frame)
			}
			prev = world.Frame
			prevTime = currentTime
		}
	}()
}

func Printf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Println(a ...interface{}) {
	log.Println(a...)
}

func dPrintf(format string, a ...interface{}) {
	if verbosity < VERB_DEBUG {
		return
	}
	Printf(format, a...)
}

func dPrintln(a ...interface{}) {
	if verbosity < VERB_DEBUG {
		return
	}
	Println(a...)
}

func iPrintf(format string, a ...interface{}) {
	if verbosity < VERB_INFO {
		return
	}
	Printf(format, a...)
}

func iPrintln(a ...interface{}) {
	if verbosity < VERB_INFO {
		return
	}
	Println(a...)
}
