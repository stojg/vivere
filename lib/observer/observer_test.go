package observer

import (
	"testing"
	"time"
)

func TestMessageRouting(t *testing.T) {

	expected := "Foo data"
	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Second)
			Publish("foo", expected)
		}
	}()

	gotMessage := false
	eventCh1 := make(chan interface{})
	Subscribe("foo", eventCh1)
	go func() {
		for {
			actual := <-eventCh1
			if actual != expected {
				t.Error("Expected %v, but got", expected, actual)
			}
			gotMessage = true
		}
	}()

	<-time.After(2 * time.Second)
	if !gotMessage {
		t.Error("Never got the message")
	}
}
