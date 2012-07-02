package megafauna_test

import (
	"math/rand"
	"megafauna"
	"strings"
	"testing"
)

func TestShuffle(t *testing.T) {
	rand.Seed(1)
	s := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	megafauna.Shuffle(s)
	expected := "2804691357"
	actual := strings.Join(s, "")
	if actual != expected {
		t.Errorf("Shuffle failed; expected %v, got %v", expected, actual)
	}
}
