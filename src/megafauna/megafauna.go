package megafauna

const BiomeTypeKeys = "LWO"
const DietaryDNAKeys = "BGHIP"
const LatitudeKeys = "THJAO"
const RoadrunnerDNAKeys = "AMNS"

// Latitude is one of the four latitude bands on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Key     string
	Name    string
	Habitat []*Habitat
}
