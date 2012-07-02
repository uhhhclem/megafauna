package megafauna

import "math/rand"

// Biome types - Land, Water, or Orogeny
const BiomeTypeKeys = "LWO"

// Dietary DNA types - Browser, Grazer, Husker, Insectivore, or Physiology
const DietaryDNAKeys = "BGHIP"

// Instinct keys - Manual dexterity, Natural history, Social skills, or Language
const InstinctKeys = "MNSL"

// Latitude keys - Tropics, Horse latitudes, Jet stream, Arctic, or Orogeny
const LatitudeKeys = "THJAO"

// Roadrunner DNA keys - Aggression, Marine, Nocturnal, or Speedy
const RoadrunnerDNAKeys = "AMNS"

// Latitude is one of the four latitude bands on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Key     string
	Name    string
	Habitat []*Habitat
}

// Shuffle implements the Fisher-Yates-Knuth shuffling algorithm.  It takes a slice of
// strings (typically keys to maps of tiles or cards) and randomizes their order.  Note
// that the caller is responsible for seeding the random number generator.
func Shuffle(keys []string) {
	length := len(keys)
	for i, _ := range keys {
		j := i + rand.Intn(length-i)
		keys[i], keys[j] = keys[j], keys[i]
	}

}
