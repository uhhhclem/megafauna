package megafauna

// Tile represents a Biome or Immigrant tile.
type Tile struct {
	Key           string // unique identifier, defined in data
	Title         string
	Subtitle      string
	LatitudeKey   string             // defined in LatitudeKeys
	IsMesozoic    bool               // if false, tile is Cenozoic
	BiomeData     *BiomeTileData     // if present, this is a biome tile
	ImmigrantData *ImmigrantTileData // if present, this is an immigrant tile
}

// BiomeTileData contains the data for biome Tiles.
type BiomeTileData struct {
	IsLand             bool
	IsWater            bool
	IsOrogeny          bool
	Niche              *Niche
	Requirements       *DNASpec
	RooterRequirements *DNASpec
	RedStar            bool
	BlueStar           bool
	ClimaxNumber       int
}

// ImmigrantTileData contains the data for immigrant Tiles.
type ImmigrantTileData struct {
	IsHerbivore bool // if false, immigrant is a 1-tooth predator
	DNA         *DNASpec
}

// IsBiomeTile indicates if the Tile is a biome.
func (t *Tile) IsBiomeTile() bool {
	return t.BiomeData != nil
}

// IsImmigrantTile indicates if the Tile is an immigrant.
func (t *Tile) IsImmigrantTile() bool {
	return t.ImmigrantData != nil
}
