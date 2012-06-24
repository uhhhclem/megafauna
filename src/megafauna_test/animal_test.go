package megafauna_test

import (
	"megafauna"
	"testing"
)

func testHerbivoreContest(t *testing.T) {
	// animal1 is suited, animal2 is not
	animal1 := megafauna.Animal { 3, 2, megafauna.MakeDNASpec("BBG")}
	animal2 := megafauna.Animal { 2, 2, megafauna.MakeDNASpec("IIM")}

	h := new(megafauna.HerbivoreContest)
	h.Animals = []*megafauna.Animal {&animal1, &animal2}
	h.Requirements = megafauna.MakeDNASpec("BB")
	h.Niche = megafauna.Niche { true, 0, "" }
	
	result := h.FindWinner()
	if result != &animal1 {
		t.Errorf("animal1 should have won on suitability.")
	}
	
	// now animal2 is suited, and it's bigger, and this is a size niche
	animal2 = megafauna.Animal { 1, 3, megafauna.MakeDNASpec("BG")}
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("animal2 should have won on size.")
	}

	// now they're both the same size, but animal2 has more teeth
	animal2 = megafauna.Animal { 2, 2, megafauna.MakeDNASpec("BG")}
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("animal2 should have won on dentition.")
	}

}
