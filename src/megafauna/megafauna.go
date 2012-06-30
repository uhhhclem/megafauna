package megafauna

// Latitude is one of the four latitude bands on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Name    string
	Habitat []*Habitat
}
