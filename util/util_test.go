package util

import (
	"testing"
	"time"
)

func TestIsOpen(t *testing.T) {
	t.Run("isopen", func(t *testing.T) {
		tim:=time.Now()
		IsOpen("159.138.30.154:8082")
		t.Log(time.Now().Sub(tim))
	})
}
