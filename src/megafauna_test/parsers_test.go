package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestParseBiome(t *testing.T) {
	biomes, err := megafauna.Parse("testdata/test_biomes.csv")
	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
	if len(biomes) == 0 {
		t.Errorf("Expected some biomes")
	}
}
