package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestGetTiles(t *testing.T) {
	tiles, err := megafauna.GetTiles()
	if err != nil {
		t.Errorf("GetTiles returned %v", err.Error())
		return
	}
	checkBiomeTiles(t, tiles)
	checkImmigrantTiles(t, tiles)
}

// checkBiomeTiles checks the tiles map for known biome tiles.
func checkBiomeTiles(t *testing.T, tiles map[string]*megafauna.Tile) {

	o := tiles["MH43"]
	if o == nil {
		t.Error("Couldn't find tile MH43.")
		return
	}
	if !o.IsMesozoic || o.LatitudeKey != "H" || o.IsLand || !o.IsSea {
		t.Error("Structured tile data for MH43 is incorrect.")
		t.Errorf("IsMesozoic=%v, LatitudeKey=%v, IsLand=%v, IsSea=%v", o.IsMesozoic, o.LatitudeKey, o.IsLand, o.IsSea)
	}
	if o.Supertitle != "Cephalopod" || o.Title != "Mediterranean Ammonites" {
		t.Error("Title data for MH43 is incorrect.")
	}

	b := o.BiomeData
	if b == nil {
		t.Error("MH43 didn't have BiomeData.")
		return
	}

	if b.ClimaxNumber != 43 || b.Requirements.Spec != "AM" || b.RooterRequirements != nil || !b.Niche.Size || b.RedStar || b.BlueStar || b.IsWarming || b.IsCooling {
		t.Error("MH43 BiomeData is incorrect.")
		if b.RooterRequirements != nil {
			t.Error("Tile has rooter requirements.")
		}
		t.Errorf("ClimaxNumber=%v, Requirements.Spec=%v, Niche.Size=%v, RedStar=%v, BlueStar=%v, IsWarming=%v, IsCooling=%v", b.ClimaxNumber, b.Requirements.Spec, b.Niche.Size, b.RedStar, b.BlueStar, b.IsWarming, b.IsCooling)
		return
	}

	o = tiles["MJ17"]
	if o == nil {
		t.Error("Couldn't find tile MJ17.")
	}
	if !o.IsMesozoic || o.LatitudeKey != "J" || !o.IsLand || o.IsSea {
		t.Error("Structured tile data for MJ17 is incorrect.")
		return
	}
	b = o.BiomeData
	if b.RooterRequirements == nil {
		t.Error("MJ17 should have rooter requirements.")
		return
	}
	if b.RooterRequirements.Spec != "H" {
		t.Error("MJ17 rooter requirements is wrong.")
	}

	o = tiles["MO94"]
	b = o.BiomeData
	if !b.IsOrogeny || !b.IsWarming {
		t.Error("Tile MO94 should be an orogeny tile that causes warming.")
	}

	for key, tile := range tiles {
		if tile.BiomeData == nil {
			continue
		}
		if tile.IsLand && tile.IsSea {
			t.Errorf("Tile %v is land and sea - that shouldn't be possible for biome tiles.", key)
		}
	}
}

// checkImmigrantTiles checks the tiles map for known immigrant tiles.
func checkImmigrantTiles(t *testing.T, tiles map[string]*megafauna.Tile) {

	o := tiles["M1"]
	if o == nil {
		t.Error("Couldn't find tile M1.")
		return
	}
	if !o.IsMesozoic || o.LatitudeKey != "A" || !o.IsLand || o.IsSea {
		t.Error("M1 basic tile data is wrong.")
	}
	if o.Supertitle != "Primitive mammal" || o.Title != "Multituberculates" {
		t.Error("M1 title data is wrong.")
	}
	if o.ImmigrantData == nil {
		t.Error("M1 doesn't have ImmigrantData.")
		return
	}
	i := o.ImmigrantData
	if i.Size != 1 || !i.IsHerbivore || i.DNA.Spec != "HIIN" {
		t.Error("ImmigrantData is wrong for tile M1.")
	}

	o = tiles["M16"]
	if o.ImmigrantData.IsHerbivore || o.ImmigrantData.Size != 0 || o.ImmigrantData.DNA.Spec != "AAM" {
		t.Error("Tile M16 is wrong.")
	}
}

func TestGetInheritanceTiles(t *testing.T) {

	tiles := megafauna.GetInheritanceTiles()
	if len(tiles) != 5 {
		t.Error("We should have gotten five inheritance tiles.")
		return
	}
	tile := tiles[4]
	if tile.Obverse.MinSize != 1 || tile.Obverse.MaxSize != 6 || tile.Obverse.DNA.Spec != "P" {
		t.Error("Obverse data is wrong for inheritance tile 5.")
	}
	if tile.Reverse.MinSize != 1 || tile.Reverse.MaxSize != 6 || tile.Reverse.DNA.Spec != "PP" {
		t.Error("Reverse data is wrong for inheritance tile 5.")
	}
	
}
