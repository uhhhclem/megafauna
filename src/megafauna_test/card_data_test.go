package megafauna_test

import (
	"megafauna"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	reader := strings.NewReader(megafauna.MutationCardData)
	cards := make(megafauna.MutationCardMap)
	err := cards.Parse(reader)
	if err != nil {
		t.Errorf("cards.Parse returned %v", err.Error())
	}
	
	var cardPantingFound, cardLungsFound bool
	
	// loop through the cards, instead of looking them up, so that we don't have any dependency on the keys,
	// which are pretty arbitrarily assigned and could easily change.
	for _, card := range cards {
		if card.Title == "Panting" {
			cardPantingFound = true
			if card.MinSize != 1 { t.Error("Expected MinSize to be 1.")}
			if card.MaxSize != 5 { t.Error("Expected MaxSize to be 5.")}
			if card.Mutation.Spec != "PP" { t.Error("Expected Mutation.Spec to be PP.") }
			if card.InstinctKey != "" { t.Error("Expected InstinctKey to be empty.") }
			if !card.Event.IsMilankovich { t.Error("Expected a Milankovich event.") }
			keys := card.Event.MilankovichLatitudeKeys
			if len(keys) != 2 || keys[0] != "H" || keys[1] != "A" {
				t.Errorf("Expected Milankovich latitudes to be HA, but it's %v.", keys)
			}
		}
		if card.Title == "Flow-Through Lungs" {
			cardLungsFound = true
			if !card.Event.IsCatastrophe { t.Error("Expected event to be catastrophe.") }
			if card.Event.Description != "Asteroid impact global cooling" {
				t.Errorf("Event description shouldn't be %v", card.Event.Description)
			}
			if card.Event.CatastropheLevel != 5 { t.Error("Should be a level 5 catastrophe.") }
			if card.Event.CatastropheIsWarming { t.Error("Should be a global-cooling catastrophe.") }
		}
	}
	
	if !cardPantingFound { t.Error("Didn't find Panting card.") }
	if !cardLungsFound { t.Error("Didn't find Flow-Through Lungs card.") }
}