package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestCards(t *testing.T) {
	g, err := megafauna.NewGame([]string{"Tinker", "Evers", "Chance"})
	if err != nil {
		t.Error(err)
		return
	}
	checkMutationCards(t, g.Cards)
	checkGenotypeCards(t, g.Cards)
}

func checkMutationCards(t *testing.T, cards map[string]*megafauna.Card) {

	var cardPantingFound, cardLungsFound bool

	// loop through the cards, instead of looking them up, so that we don't have any dependency on the keys,
	// which are pretty arbitrarily assigned and could easily change.
	for _, card := range cards {
		if card.Mutation == nil {
			continue
		}
		m := card.Mutation
		if m.Title == "Panting" {
			cardPantingFound = true
			if m.MinSize != 1 {
				t.Error("Expected MinSize to be 1.")
			}
			if m.MaxSize != 5 {
				t.Error("Expected MaxSize to be 5.")
			}
			if m.Mutation.Spec != "PP" {
				t.Error("Expected Mutation.Spec to be PP.")
			}
			if m.InstinctKey != "" {
				t.Error("Expected InstinctKey to be empty.")
			}
			if !card.Event.IsMilankovich {
				t.Error("Expected a Milankovich event.")
			}
			keys := card.Event.MilankovichLatitudeKeys
			if len(keys) != 2 || keys[0] != "H" || keys[1] != "A" {
				t.Errorf("Expected Milankovich latitudes to be HA, but it's %v.", keys)
			}
		}
		if m.Title == "Flow-Through Lungs" {
			cardLungsFound = true
			if !card.Event.IsCatastrophe {
				t.Error("Expected event to be catastrophe.")
			}
			if card.Event.Description != "Asteroid impact global cooling" {
				t.Errorf("Event description shouldn't be %v", card.Event.Description)
			}
			if card.Event.CatastropheLevel != 5 {
				t.Error("Should be a level 5 catastrophe.")
			}
			if card.Event.IsWarming {
				t.Error("Should be a global-cooling catastrophe.")
			}
		}
	}

	if !cardPantingFound {
		t.Error("Didn't find Panting card.")
	}
	if !cardLungsFound {
		t.Error("Didn't find Flow-Through Lungs card.")
	}
}

func checkGenotypeCards(t *testing.T, cards map[string]*megafauna.Card) {

	// G3,rhino,Artodactyl ungulate,Swine,"pigs, hippos",1,4,GP,dino,Ornithischian ornithopod,Duckbills,"lambeosaurines, iguanodonts, hadrosaurs",2,5,GG,T,,,,

	card := cards["G3"]
	if card.Genotype == nil {
		t.Error("Card G3 should be a genotype card.")
		return
	}
	m := card.Genotype.MammalData
	if m.SilhouetteIndex != 1 ||
		m.Family != "Artodactyl ungulate" ||
		m.Title != "Swine" ||
		m.Subtitle != "pigs, hippos" ||
		m.MinSize != 1 ||
		m.MaxSize != 4 ||
		m.DNASpec.Spec != "GP" {
		t.Error("Didn't parse mammal data correctly.")
	}
	d := card.Genotype.DinosaurData
	if d.SilhouetteIndex != 0 ||
		d.Family != "Ornithischian ornithopod" ||
		d.Title != "Duckbills" ||
		d.Subtitle != "lambeosaurines, iguanodonts, hadrosaurs" ||
		d.DNASpec.Spec != "GG" {
		t.Error("Didn't parse dinosaur data correctly.")
	}
	if !card.Event.IsDrawTwo {
		t.Error("Didn't parse event data correctly.")
	}
}
