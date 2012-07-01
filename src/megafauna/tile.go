package megafauna

type Tile struct {
	Key string
	Title              string
	Subtitle           string
	LatitudeKey        string
	IsMesozoic         bool
}

type BiomeTile struct {
	Tile
	IsLand			   bool
	IsWater            bool
	IsOrogeny          bool
	Requirements       *DNASpec
	Niche              *Niche
	RooterRequirements *DNASpec
	RedStar            bool
	BlueStar           bool
	ClimaxNumber       int
}

type ImmigrantTile struct {
	Tile
	IsHerbivore        bool
	DNA				   *DNASpec
}