package megafauna

// Biome represents a biome on the board, i.e. a habitat with a biome tile in it, along with
// (possible) animals in its predator, herbivore, and rooter slots.  
//
// Note that since a biome
// by definition has a biome tile in it, methods like GetClimaxNumber() don't check to see
// if Tile is nil or Tile.BiomeTileData is nil.
type Biome struct {
	Key       string
	Tile      *Tile // this will be either a biome or an immigrant
	Predator  []*Animal
	Herbivore []*Animal
	Rooter    []*Animal
	Habitat   *Habitat
}

// GetClimaxNumber returns a Biome's ClimaxNumber.
func (b *Biome) GetClimaxNumber() int {
	return b.Tile.BiomeData.ClimaxNumber
}
