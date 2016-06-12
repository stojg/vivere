package creator

import (
	"testing"
)

func TestInit(t *testing.T) {

	c := NewCreator(1, 32, 200, 200)

	if c.tileSize != 32 {
		t.Error("tile size wasnt set.")
	}
}
