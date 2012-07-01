package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestHerbivoreContestBasicTests(t *testing.T) {
	// animal1 is suited, animal2 is not
	animal1 := megafauna.Animal{2, 2, megafauna.MakeDNASpec("BBG")}
	animal2 := megafauna.Animal{3, 2, megafauna.MakeDNASpec("IIM")}

	h := new(megafauna.HerbivoreContest)
	h.Animals = []*megafauna.Animal{&animal1, &animal2}
	h.Requirements = megafauna.MakeDNASpec("BB")

	var err error
	h.Niche, err = megafauna.MakeNiche("Size")
	if err != nil {
		t.Error(err)
	}

	result := h.FindWinner()
	if result != &animal1 {
		t.Errorf("Scores were: %v", h.Scores)
		t.Errorf("animal1 should have won on suitability.")
	}

	// now animal2 is suited, and it's bigger, and this is a size niche
	animal2 = megafauna.Animal{3, 3, megafauna.MakeDNASpec("BB")}
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("Scores were: %v", h.Scores)
		t.Errorf("animal2 should have won on size.")
	}

	// now animal2 is suited, and this is an "I" niche
	animal2 = megafauna.Animal{3, 2, megafauna.MakeDNASpec("BBI")}
	h.Niche, err = megafauna.MakeNiche("I")
	if err != nil {
		t.Error(err)
	}
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("Scores were: %v", h.Scores)
		t.Errorf("animal2 should have won because it's an I niche.")
	}

	// what if it's a player-color niche, i.e. a niche with animal1's dentition?
	h.Niche, err = megafauna.MakeNiche("2")
	if err != nil {
		t.Error(err)
	}
	result = h.FindWinner()
	if result != &animal1 {
		t.Errorf("Scores were: %v", h.Scores)
		t.Errorf("animal1 should have won its player-color niche.")
	}

	// now they're both the same size (and it's a Size niche), but animal2 has more teeth
	animal2 = megafauna.Animal{3, 2, megafauna.MakeDNASpec("BB")}
	h.Niche, err = megafauna.MakeNiche("Size")
	if err != nil {
		t.Error(err)
	}
	result = h.FindWinner()
	if result != &animal2 {
		t.Errorf("Scores were: %v", h.Scores)
		t.Errorf("animal2 should have won on dentition.")
	}

}
