package megafauna

// Animal can be either an immigrant or a player's species.
type Animal struct {
	Dentition int
	Size      int
	Genome    DNASpec
}

// Immigrant represents an immigrant tile.
type Immigrant struct {
	Animal
	Title    string
	Subtitle string
	Latitude *Latitude
	IsLand   bool
	IsSea    bool
}

// Species represents one of the player's four species
type Species struct {
	Animal
}

// HerbivoreContest is not very well thought out.  The idea is that in order to figure out who wins in an herbivore
// contest, you'll build a slice of animals and sort them.  I think the actual solution is going to be more of a
// score-based mechanic, e.g. add 100 points if the animal matches the biome, add 10 if it matches the niche, and then
// add the Dentition, and everyone under 100 dies and the winner is the survivor with the highest score.
type HerbivoreContest struct {
	Animals      []Animal
	Requirements DNASpec
	Niche        string
}

// Len is part of the sort.Interface interface.
func (h HerbivoreContest) Len() int {
	return len(h.Animals)
}

// Swap is part of the sort.Interface interface.
func (h HerbivoreContest) Swap(i, j int) {
	h.Animals[i], h.Animals[j] = h.Animals[j], h.Animals[i]
}

// Less is part of the sort.Interface interface.
func (h HerbivoreContest) Less(i, j int) {
	// a1, a2 := h.Animals[i], h.Animals[j]
	// suitability test goes here
}
