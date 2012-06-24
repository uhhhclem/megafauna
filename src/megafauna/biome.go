package megafauna

// Biome represents a biome on the board, i.e. a Habitat with a biome tile in it, along with
// (possible) animals in its predator, herbivore, and rooter slots.
type Biome struct {
	Title string
	Subtitle string
	ClimaxNumber int
	Requirements DNASpec
	Niche Niche
	RooterRequirements DNASpec
	Predator []*Animal
	Herbivore []*Animal
	Rooter []*Animal
	Habitat *Habitat
}

