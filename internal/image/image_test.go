package image

import (
	"testing"
)

func TestCreateImage(t *testing.T) {
	err := LastListenedTo("LA really long track woooooooooooooooooooooooooooooooooow", "artist")
	if err != nil {
		t.Fatal(err)
	}
}
