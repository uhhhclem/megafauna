package megafauna_test

import (
	"megafauna"
    "testing"
)

func TestParseBiome(t *testing.T) {
	biomes, err := megafauna.Parse("testdata/test_biomes.json")
	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
	if len(biomes) == 0 {
		t.Errorf("Expected some biomes")
	}
}

func TestEmitBiomes(t *testing.T) {
	biome1 := megafauna.Biome{
		Title:	"Test Biome",
		Subtitle: "This is the subtitle",
		ClimaxNumber: 1,
		Requirements: megafauna.DNASpec{
			Spec: "BBG",
			Breakdown: nil,
		}
	}
	biomes := []megafauna.Biome{biome1}
	err := megafauna.Emit(biomes)
	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
	
}