package megafauna

// Habitat is a space on the board that contains a slot for the biome, the predator triangle,
// and the rooter triangle.
type Habitat struct {
	Latitude         *Latitude
	ClimaxNumber     int
	IsOrogeny        bool
	AdjacentHabitats []*Habitat
	Biome            *Biome
}
