package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	verbosity int = VERB_NORM
)

const (
	VERB_NORM  int = 0
	VERB_INFO  int = 1
	VERB_DEBUG int = 2
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	envVerbosity := os.Getenv("VERBOSITY")
	if envVerbosity == "" {
		verbosity = 0
	} else {
		verbosity, _ = strconv.Atoi(envVerbosity)
	}
}

func DebugFPS(framesPerSec float64) {
	warningFPS := (1 / framesPerSec) - 1

	go func() {
		timer := 1 * time.Second
		prev := Frame
		prevTime := time.Now()
		for {
			currentTime := <-time.After(timer)
			fps := float64(Frame-prev) / float64(currentTime.Sub(prevTime).Seconds())
			if fps < warningFPS {
				Printf("fps: %0.1f < %0.1f frame %d\n", fps, warningFPS, Frame)
			} else {
				dPrintf("fps: %0.1f frame %d\n", fps, Frame)
			}
			prev = Frame
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
