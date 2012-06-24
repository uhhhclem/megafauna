package megafauna

import (
	"sort" 
)


// Animal can be either an immigrant or a player's species.
type Animal struct {
	Dentition int
	Size      int
	Genome    *DNASpec
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
	Animals      []*Animal
	Scores	     []int
	Requirements *DNASpec
	Niche        Niche
}

// Len is part of the sort.Interface interface.
func (h HerbivoreContest) Len() int {
	return len(h.Animals)
}

// Swap is part of the sort.Interface interface.  Here we swap the animals and their scores.
func (h HerbivoreContest) Swap(i, j int) {
	h.Animals[i], h.Animals[j] = h.Animals[j], h.Animals[i]
	h.Scores[i], h.Scores[j] = h.Scores[j], h.Scores[i]
}

// Less is part of the sort.Interface interface.
func (h HerbivoreContest) Less(i, j int) bool {
	return h.Scores[i] < h.Scores[j]
}

// Reverse embeds a sort.Interface value and implements a reverse sort over
// that value.
type Reverse struct {
    // This embedded Interface permits Reverse to use the methods of
    // another Interface implementation.
    sort.Interface
}

// Less returns the opposite of the embedded implementation's Less method.
func (r Reverse) Less(i, j int) bool {
    return r.Interface.Less(j, i)
}

// FindWinner assigns scores to the animals in the contest, then sorts in reverse order to
// find the winner.  It returns nil if there are no animals, or if none is suitable to
// the Requirements.
func (h HerbivoreContest) FindWinner() *Animal {
	if len(h.Animals) == 0 {
		return nil
	}
	h.Scores = make([]int, len(h.Animals))
	
	// calculate the score for each animal
	for i, animal := range h.Animals {
		// suitability score is in the hundreds
		if animal.Genome.CanFeedOn(h.Requirements) {
			h.Scores[i] += 100
		}
		// niche score is in the tens
		if h.Niche.Size {
			h.Scores[i] += 10 * animal.Size
		} else if h.Niche.Dentition != 0 {
			if h.Niche.Dentition == animal.Dentition {
				h.Scores[i] += 10
			}
		} else if h.Niche.DNA != "" {
			_, ok := animal.Genome.Breakdown[h.Niche.DNA]
			if ok {
				h.Scores[i] += 10
			}
		} 
		// and the ones is the animal's dentition
		h.Scores[i] += animal.Dentition
	}
	sort.Sort(Reverse{h})
	if h.Scores[0] >= 100 {
		return h.Animals[0]
	}
	return nil
}