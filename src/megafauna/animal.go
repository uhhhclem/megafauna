package megafauna

import (
	"fmt"
	"sort"
)

// Reverse embeds a sort.Interface value and implements a reverse sort over
// that value.
type Reverse struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
	sort.Interface
}

// Animal can be either an immigrant or a player's species.
type Animal struct {
	Dentition     int
	Size          int
	Genome        *DNASpec
	ImmigrantTile *Tile // if nil, this is a player species (also, dentition must be 2-5)
	Silhouette    int   // if a player species, the silhouette (0-3) chosen for this species
}

// HerbivoreContest is used to determine the winner of herbivore contests during the cull.  Set Animals, Requirements, 
// and Niche from a Biome containing herbivores, and then call FindWinner to find out who (if anyone) survived.
type HerbivoreContest struct {
	Animals      []*Animal
	Scores       []int
	Requirements *DNASpec
	Niche        *Niche
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

// Less returns the opposite of the embedded implementation's Less method.
func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

// FindWinner assigns scores to the animals in the contest, then sorts in reverse order to
// find the winner.  It returns nil if there are no animals, or if none is suitable to
// the Requirements.
func (h *HerbivoreContest) FindWinner() *Animal {
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
	fmt.Printf("Scores are %v\n", h.Scores)
	// an animal has to have a score of at least 100 to survive; among those, we pick
	// the one with the highest score.
	if h.Scores[0] >= 100 {
		return h.Animals[0]
	}

	// nobody has a score of over 100; valar morghulis.
	return nil
}

// CarnivoreContest is used to determine the winner of carnivore contests during the cull.  Set Carnivores,
// and Prey (which will have 0, 1, or 2 members), and call FindWinners to find out which
// (if any) carnivores survive.
type CarnivoreContest struct {
	Carnivores []*Animal
	Prey       []*Animal
}

// FindWinners finds the winner(s) (if any) of a carnivore contest. 
func (contest *CarnivoreContest) FindWinners() []*Animal {

	winners := make([]*Animal, 0)

	for _, p := range contest.Prey {
		feeders = make([]*Animal)
		for _, c := range contest.Carnivores {
			if c.canFeedOn(p) {
				feeders = append(feeders, c)
			}
		}
		// if anyone can feed on this prey, find out which one actually wins
		if len(feeders) > 0 {
			winners = winners.append(result, findWinningCarnivore(feeders)
		}
	}

	return winners
}

// canFeedOn indicates whether or not prey is suitable for the carnivore
func (carnivore *Animal) canFeedOn(prey *Animal) bool {

	// no cannibalism!
	if carnivore.Dentition > 1 && carnivore.Dentition < 6 && carnivore.Dentition == prey.Dentition {
		if carnivore.Silhouette == prey.Silhouette {
			return false
		}
	}

	// too big or too small!
	if carnivore.Size < prey.Size-1 || carnivore.Size > prey.Size+1 {
		return false
	}

	// can you catch me?
	return carnivore.Genome.CanFeedOn(prey.Genome)
}

// findWinningCarnivore finds the one highest-scoring carnivore among those that can eat a specific prey.
func findWinningCarnivore(carnivores *[]Animal) {
	var winner *Animal
	winningP := 0
	winningDentition := 10
	for _, c := range carnivores {
	
		// who has the most P DNA?
		p := c.Genome.GetValue("P")
		if p < winningP {
			continue
		}
		winningP = p
		
		// among those, who has the fewest teeth?
		if c.Dentition > winningDentition {
			continue
		}
		winningDentition = c.Dentition
		winner = c
	}
	return winner
}