package megafauna

// Latitude is one of the four latitude bands on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Key string
	Name    string
	Habitat []*Habitat
}
