package megafauna_test

import (
	"megafauna"
	"strings"
	"testing"
)

func TestBiomeMapParse(t *testing.T) {
	data := "1,Variscan orogeny,African Podocarp High Forest,67,B,,SIZE,false,false\n2,Cypressales,Dawn Redwood Forest,40,BB,H,A,false,true"

	biomes := make(megafauna.BiomeMap)

	reader := strings.NewReader(data)
	err := biomes.Parse(reader)
	if err != nil {
		t.Errorf("An error occurred: %v", err)
	}
	if len(biomes) == 0 {
		t.Errorf("Expected some biomes.")
	}
	
	b := biomes["2"]
	expectedTitle := "Cypressales"
	expectedSubtitle := "Dawn Redwood Forest"
	expectedClimaxNumber := 40
	expectedRequirementsSpec := "BB"
	expectedRooterRequirementsSpec := "H"
	expectedNicheDNA := "A"
	
	if b.Title != expectedTitle {
		t.Errorf("Expected Title to be %v, but it's %v", expectedTitle, b.Title)
	}
	if b.Subtitle != expectedSubtitle {
		t.Errorf("Expected Subtitle to be %v, but it's %v", expectedSubtitle, b.Subtitle)
	}
	if b.ClimaxNumber != expectedClimaxNumber {
		t.Errorf("Expected ClimaxNumber to be %v, but it's %v", expectedClimaxNumber, b.ClimaxNumber)
	}
	if b.Requirements.Spec != expectedRequirementsSpec {
		t.Errorf("Expected Requirements.Spec to be %v, but it's %v", expectedRequirementsSpec, b.Requirements.Spec)
	}
	if b.RooterRequirements.Spec != expectedRooterRequirementsSpec {
		t.Errorf("Expected RooterRequirements.Spec to be %v, but it's %v", expectedRooterRequirementsSpec, b.RooterRequirements.Spec)
	}
	if b.Niche.DNA != expectedNicheDNA {
		t.Errorf("Expected Niche.DNA to be %v, but it's %v", expectedNicheDNA, b.Niche.DNA)
	}
	if b.RedStar {
		t.Errorf("Expected RedStar to be false, but it's true.")
	}
	if !b.BlueStar {
		t.Errorf("Expected BlueStar to be true, but it's false.")
	}
	if b.Niche.DNA != expectedNicheDNA {
		t.Errorf("Expected Niche.DNA to be %v, but it's %v", expectedNicheDNA, b.Niche.DNA)
	}
}

func TestLatitudeMapParse(t *testing.T) {
	data := "1,Tropics\n2,Horse Latitudes\n3,Jet Stream\n4,Arctic"

	latitudes := make(megafauna.LatitudeMap)

	reader := strings.NewReader(data)
	err := latitudes.Parse(reader)
	if err != nil {
		t.Errorf("An error occurred: %v", err)
	}
	if len(latitudes) == 0 {
		t.Errorf("Expected some latitudes.")
	}
	
	key := "2"
	expected := "Horse Latitudes"
	if latitudes[key].Name != expected {
		t.Errorf("Expected Latitude with a key of %v to be %v, but it's %v.", key, expected, latitudes[key].Name) 
	}
}