package megafauna

// Niche is used as the first tiebreaker in herbivore contests.  A Niche will be *one* of these: Size, 
// a player color (Dentition), or a DNA letter.
type Niche struct {
	Size      bool
	Dentition int
	DNA       string
}

// Latitude is one of the four latitude bands on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Name    string
	Habitat []*Habitat
}

// Habitat is a space on the board that contains a slot for the biome, the predator triangle,
// and the rooter triangle.
type Habitat struct {
	Latitude         *Latitude
	ClimaxNumber     int
	AdjacentHabitats []*Habitat
	Biome            *Biome
}
