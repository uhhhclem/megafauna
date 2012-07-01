package megafauna

// Tile represents a Biome or Immigrant tile.
type Tile struct {
	Key         string // unique identifier, defined in data
	Title       string
	Subtitle    string
	LatitudeKey string // defined in LatitudeKeys
	IsMesozoic  bool   // if false, tile is Cenozoic
}

type BiomeTile struct {
	Tile
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

type ImmigrantTile struct {
	Tile
	IsHerbivore bool // if false, immigrant is a 1-tooth predator
	DNA         *DNASpec
}
