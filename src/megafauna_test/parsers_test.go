package megafauna_test

import (
	"megafauna"
	"strings"
	"testing"
)

func TestBiomeTileMapParse(te *testing.T) {
	data := `1,O,O,Variscan orogeny,African Podocarp High Forest,67,B,,SIZE,false,false
2,T,L,Cypressales,Dawn Redwood Forest,40,BB,H,A,false,true`

	biomes := make(megafauna.TileMap)

	reader := strings.NewReader(data)
	err := biomes.Parse(reader)
	if err != nil {
		te.Errorf("An error occurred: %v", err)
	}
	if len(biomes) == 0 {
		te.Errorf("Expected some biomes.")
	}

	var t *megafauna.Tile
	var b *megafauna.BiomeTileData

	t = biomes["1"]
	b = t.BiomeData

	if t.LatitudeKey != "O" {
		te.Errorf("Expected tile 1 to be in the O latitude.")
	}
	if !b.IsOrogeny || t.IsWater || !t.IsLand {
		te.Errorf("Expected tile 1 to be an orogeny tile.")
	}

	t = biomes["2"]
	b = t.BiomeData

	if t.LatitudeKey != "T" {
		te.Errorf("Expected tile 2 to be in the T latitude.")
	}
	if b.IsOrogeny || t.IsWater || !t.IsLand {
		te.Errorf("Expected tile 2 to be a land tile.")
	}

	expectedTitle := "Cypressales"
	expectedSubtitle := "Dawn Redwood Forest"
	expectedClimaxNumber := 40
	expectedRequirementsSpec := "BB"
	expectedRooterRequirementsSpec := "H"
	expectedNicheDNA := "A"

	if t.Title != expectedTitle {
		te.Errorf("Expected Title to be %v, but it's %v", expectedTitle, t.Title)
	}
	if t.Subtitle != expectedSubtitle {
		te.Errorf("Expected Subtitle to be %v, but it's %v", expectedSubtitle, t.Subtitle)
	}
	if b.ClimaxNumber != expectedClimaxNumber {
		te.Errorf("Expected ClimaxNumber to be %v, but it's %v", expectedClimaxNumber, b.ClimaxNumber)
	}
	if b.Requirements.Spec != expectedRequirementsSpec {
		te.Errorf("Expected Requirements.Spec to be %v, but it's %v", expectedRequirementsSpec, b.Requirements.Spec)
	}
	if b.RooterRequirements.Spec != expectedRooterRequirementsSpec {
		te.Errorf("Expected RooterRequirements.Spec to be %v, but it's %v", expectedRooterRequirementsSpec, b.RooterRequirements.Spec)
	}
	if b.Niche.DNA != expectedNicheDNA {
		te.Errorf("Expected Niche.DNA to be %v, but it's %v", expectedNicheDNA, b.Niche.DNA)
	}
	if b.RedStar {
		te.Errorf("Expected RedStar to be false, but it's true.")
	}
	if !b.BlueStar {
		te.Errorf("Expected BlueStar to be true, but it's false.")
	}
	if b.Niche.DNA != expectedNicheDNA {
		te.Errorf("Expected Niche.DNA to be %v, but it's %v", expectedNicheDNA, b.Niche.DNA)
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
