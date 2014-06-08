package creator

import (
	"testing"
)

func TestInit(t *testing.T) {

	c := &Creator{}
	c.Init(32, 200, 200)

	if c.tileSize != 32 {
		t.Error("tile size wasnt set.")
	}
}
