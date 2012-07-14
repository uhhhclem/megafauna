package megafauna

import "math/rand"
import "time"

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

// ValidDNAKeys contains all valid DNA keys
const ValidDNAKeys = DietaryDNAKeys + RoadrunnerDNAKeys

// Latitude is one of the four latitude bands on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Key     string
	Name    string
	Habitat []*Habitat
}

// SeedRand seeds the RNG for calls to Shuffle.  If this is 0, Shuffle will use the system time as the seed.
func SeedRand(value int64) {
	if value == 0 {
		value = time.Now().Unix()
	}
	rand.Seed(value)
}

// Shuffle implements the Fisher-Yates-Knuth shuffling algorithm.  It takes a slice of
// strings (typically keys to maps of tiles or cards) and randomizes their order.  This
// automatically seeds the RNG unless SetSeed is called first.  Note that this resets the
// seed to 0; if you want Shuffle to behave deterministically, call SetSeed before every call.
func Shuffle(keys []string) {
	length := len(keys)
	for i, _ := range keys {
		j := i + rand.Intn(length-i)
		keys[i], keys[j] = keys[j], keys[i]
	}

}
