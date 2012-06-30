package megafauna_test

import (
	"megafauna"
	"strings"
	"testing"
)

func TestBiomeSliceParse(t *testing.T) {
	data := "Variscan orogeny,African Podocarp High Forest,67,B,,SIZE,false,false\nCypressales,Dawn Redwood Forest,40,BBH,,A,false,true"

	biomes := make(megafauna.BiomeSlice, 300)

	reader := strings.NewReader(data)
	err := biomes.Parse(reader)
	if err != nil {
		t.Errorf("An error occurred: %v", err)
	}
	if len(biomes) == 0 {
		t.Errorf("Expected some biomes.")
	}
	
	b := biomes[1]
	expectedTitle := "Cypressales"
	expectedSubtitle := "Dawn Redwood Forest"
	expectedClimaxNumber := 40
	expectedRequirementsSpec := "BBH"
	expectedRooterRequirementsSpec := ""
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