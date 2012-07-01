package megafauna

// Biome represents a biome on the board, i.e. a Habitat with a biome tile in it, along with
// (possible) animals in its predator, herbivore, and rooter slots.
type Biome struct {
	Key				   string
	Tile			   *BiomeTile
	Predator           []*Animal
	Herbivore          []*Animal
	Rooter             []*Animal
	Habitat            *Habitat
}
