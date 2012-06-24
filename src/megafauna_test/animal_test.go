package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestHerbivoreContest(t *testing.T) {
	// animal1 is suited, animal2 is not
	animal1 := megafauna.Animal { 1, 2, megafauna.MakeDNASpec("BBG")}
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
	animal2 = megafauna.Animal { 2, 3, megafauna.MakeDNASpec("BG")}
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("animal2 should have won on size.")
	}

	// now animal2 is suited, and this is an "I" niche
	animal2 = megafauna.Animal { 2, 2, megafauna.MakeDNASpec("BGI")}
	h.Niche = megafauna.Niche { false, 0, "I" }
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("h.Scores = %v", h.Scores)
		t.Errorf("animal2 should have won because it's an I niche.")
	}

	// what if it's a player-color niche, i.e. a niche with animal1's dentition?
	h.Niche = megafauna.Niche { false, animal1.Dentition, "" }
	result = h.FindWinner()
	if result != &animal1 {
		t.Errorf("h.Scores = %v", h.Scores)
		t.Errorf("animal1 should have won its player-color niche.")
	}

	// now they're both the same size (and it's a Size niche), but animal2 has more teeth
	animal2 = megafauna.Animal { 2, 2, megafauna.MakeDNASpec("BG")}
	h.Niche = megafauna.Niche { true, 0, "" }
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("h.Scores = %v", h.Scores)
		t.Errorf("animal2 should have won on dentition.")
	}

}
