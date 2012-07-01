package megafauna

// Biome represents a biome on the board, i.e. a habitat with a biome tile in it, along with
// (possible) animals in its predator, herbivore, and rooter slots.
type Biome struct {
	Key       string
	BiomeTile *BiomeTile
	Predator  []*Animal
	Herbivore []*Animal
	Rooter    []*Animal
	Habitat   *Habitat
}
