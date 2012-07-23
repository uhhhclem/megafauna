package megafauna_test

import (
	"megafauna"
	"strings"
	"testing"
)

func TestMakeEvent(t *testing.T) {
	var e *megafauna.Event
	var err error

	e, err = megafauna.MakeEvent("Q", "", 0, false)
	if err != megafauna.ErrInvalidEventType {
		t.Error("Expected ErrInvalidEventType and didn't get it.")
	}

	e, err = megafauna.MakeEvent("T", "", 0, false)
	if !e.IsDrawTwo {
		t.Error("Expected event to be a draw-two event, and it isn't.")
	}

	e, err = megafauna.MakeEvent("C", "", 5, true)
	if !e.IsCatastrophe {
		t.Error("Expected event to be a catastrophe, and it isn't.")
	}
	if e.CatastropheLevel != 5 {
		t.Errorf("Expected catastrophe level to be 5, and it's %v.", e.CatastropheLevel)
	}
	if !e.CatastropheIsWarming {
		t.Error("Expected catastrophe to be warming, and it's cooling.")
	}

	e, err = megafauna.MakeEvent("MP", "HA", 0, false)
	if !e.IsMilankovich {
		t.Error("Expected event to be a Milankovich event, and it isn't.")
	}
	keys := e.MilankovichLatitudeKeys
	if len(keys) != 2 || keys[0] != "H" || keys[1] != "A" {
		t.Errorf("Milankovich latitudes didn't parse; they're %v.", e.MilankovichLatitudeKeys)
	}

}

func TestParseMutationCards(t *testing.T) {
	reader := strings.NewReader(megafauna.MutationCardSourceData)
	cards := make(map[string]interface{})
	err := megafauna.ParseMutationCards(reader, cards)
	if err != nil {
		t.Errorf("ParseMutationCards returned %v", err.Error())
	}

	var card *megafauna.MutationCard
	var cardPantingFound, cardLungsFound bool

	// loop through the cards, instead of looking them up, so that we don't have any dependency on the keys,
	// which are pretty arbitrarily assigned and could easily change.
	for _, v := range cards {
		switch v.(type) {
		default:
			continue
		case *megafauna.MutationCard:
			card = v.(*megafauna.MutationCard)
		}
		if card.Title == "Panting" {
			cardPantingFound = true
			if card.MinSize != 1 {
				t.Error("Expected MinSize to be 1.")
			}
			if card.MaxSize != 5 {
				t.Error("Expected MaxSize to be 5.")
			}
			if card.Mutation.Spec != "PP" {
				t.Error("Expected Mutation.Spec to be PP.")
			}
			if card.InstinctKey != "" {
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
		if card.Title == "Flow-Through Lungs" {
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
			if card.Event.CatastropheIsWarming {
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

func TestParseGenotypeCards(t *testing.T) {
	reader := strings.NewReader(megafauna.GenotypeCardSourceData)
	cards := make(map[string]interface{})
	err := megafauna.ParseGenotypeCards(reader, cards)
	if err != nil {
		t.Errorf("ParseGenotypeCards returned %v", err.Error())
	}

	// G3,rhino,Artodactyl ungulate,Swine,"pigs, hippos",1,4,GP,dino,Ornithischian ornithopod,Duckbills,"lambeosaurines, iguanodonts, hadrosaurs",2,5,GG,T,,,,

	card := cards["G3"].(*megafauna.GenotypeCard)
	m := card.MammalData
	if m.SilhouetteIndex != 1 ||
		m.Family != "Artodactyl ungulate" ||
		m.Title != "Swine" ||
		m.Subtitle != "pigs, hippos" ||
		m.MinSize != 1 ||
		m.MaxSize != 4 ||
		m.DNASpec.Spec != "GP" {
		t.Error("Didn't parse mammal data correctly.")
	}
	d := card.DinosaurData
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

func TestParseErrors(t *testing.T) {
	const invalidDNA = "8,1,4,NQ,N,,Infrared Pit Sensor,The ability to sense thermal radiation helps to detect warm-blooded predators or prey.,T,,,,\n"
	var err error

	r := strings.NewReader(invalidDNA)
	cards := make(map[string]interface{})
	err = megafauna.ParseMutationCards(r, cards)
	if err != megafauna.ErrInvalidDNASpec {
		t.Error("Expected ErrInvalidDNASpec and didn't get it.")
	}

	const invalidEventData = "8,1,4,N,N,,Infrared Pit Sensor,The ability to sense thermal radiation helps to detect warm-blooded predators or prey.,Q,,,,\n"
	r = strings.NewReader(invalidEventData)
	err = megafauna.ParseMutationCards(r, cards)
	if err != megafauna.ErrInvalidEventType {
		t.Error("Expected ErrInvalidEventType and didn't get it.")
	}

}
