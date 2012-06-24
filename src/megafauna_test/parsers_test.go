package megafauna_test

import (
	"megafauna"
    "testing"
    "os"
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
			Breakdown: map[string] megafauna.DNA {
				"B": megafauna.DNA {
					Letter: "B",
					Value: 0,},
			},
		}, 
		Niche: megafauna.Niche {
			DNA: "P",
		},
		/**
		RooterRequirements: nil,
		Predator: nil,
		Herbivore: nil,
		Rooter: nil,
		Habitat: nil,
		**/
	}
	biomes := []megafauna.Biome{biome1}
	jsonBiomes, err := megafauna.Emit(biomes)
	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
	os.Stdout.Write(jsonBiomes)
	
}