package fsm

import (
	"testing"
)

func TestNewColors(t *testing.T) {
	colors := NewColors()
	if len(colors.values) == 0 {
		t.Fatal("no colors")
	}
}

func TestColors_Pick(t *testing.T) {
	colors := NewColors()
	var picked = make(map[string]int)
	for i := 0; i < len(availableColors)*2; i++ {
		v := colors.Pick()
		if v == "" {
			t.Fatal("value not picked")
		}
		picked[v]++
	}

	for k, v := range picked {
		if v != 2 {
			t.Fatalf("excepted %s picked 2 times, picked %d", k, v)
		}
	}
}
