package megafauna

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Game is one discrete game of Bios Megafauna.
type Game struct {
	Players       SortablePlayerCollection // slice of Player objects, in player order
	MutationCards MutationCardMap          // master map of all mutation cards
	GenotypeCards GenotypeCardMap          // master map of all genotype cards
	CardKeys      []string                 // shuffled slice of keys to all cards
}

// SortablePlayerCollection is used to sort Players by Dentition.
type SortablePlayerCollection []*Player

// Len returns the length of a SortablePlayerCollection.
func (p SortablePlayerCollection) Len() int {
	return len(p)
}

// Less returns which Player in a SortablePlayerCollection has fewer teeth.
func (p SortablePlayerCollection) Less(i, j int) bool {
	return p[i].Dentition < p[j].Dentition
}

// Swap swaps Players in a SortablePlayerCollection. 
func (p SortablePlayerCollection) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// NewGame creates a new Game and initializes the Players.
func NewGame(names []string) *Game {
	g := new(Game)
	g.createPlayers(names)
	if g.Players == nil {
		return nil
	}
	return g
}

// createPlayers creates the Player objects in the game from a list of their names.  If it fails, Players will be nil.
func (g *Game) createPlayers(names []string) {
	if len(names) < 2 || len(names) > 4 {
		return
	}

	// create and name the players
	players := make(SortablePlayerCollection, len(names))
	for index, name := range names {
		p := new(Player)
		p.Name = name
		players[index] = p
	}

	// Randomly assign dentitions to each player.  Shuffle takes strings, so we'll have to convert
	// the dentitions to ints.  
	dentitions := []string{"2", "3", "4", "5"}
	colors := []string{"Red", "Orange", "Green", "White"}
	Shuffle(dentitions)
	for index, p := range players {
		p.Dentition, _ = strconv.Atoi(dentitions[index])
		p.Color = colors[p.Dentition-2]
		p.IsDinosaur = p.Dentition == 2 || p.Dentition == 4
		p.Species = make([]*Species, 4)
		p.AnimalTokens = []int{8, 8, 8, 8}
	}

	// sort the players by Dentition
	sort.Sort(players)

	// everyone gets 4 genes except the starting player
	for i, p := range players {
		p.Genes = 4
		if i == 0 {
			p.Genes -= 1
		}
	}

	g.Players = players

}

// createCards initializes the deck.
func (g *Game) createCards() {
	g.CardKeys = make([]string, 0)

	// get the mutation cards
	g.MutationCards = make(MutationCardMap)
	r := strings.NewReader(MutationCardSourceData)
	g.MutationCards.Parse(r)

	// get the genotype cards (this it TODO)
	g.GenotypeCards = make(GenotypeCardMap)

	// add the keys from both collections to CardKeys
	for k, _ := range g.MutationCards {
		g.CardKeys = append(g.CardKeys, k)
	}
	for k, _ := range g.GenotypeCards {
		g.CardKeys = append(g.CardKeys, k)
	}
	Shuffle(g.CardKeys)
}

// GetCard returns the MutationCard or GenotypeCard for a given key.
func (g *Game) GetCard(key string) (*MutationCard, *GenotypeCard) {
	mut := g.MutationCards[key]
	gen := g.GenotypeCards[key]
	return mut, gen
}

// GetPlayer returns the Player with the given dentition.
func (g *Game) GetPlayer(dentition int) *Player {
	for _, p := range g.Players {
		if p.Dentition == dentition {
			return p
		}
	}
	return nil
}

// Player defines a player and his resources.
type Player struct {
	Name         string     // the player's name
	Color        string     // the player's color
	Dentition    int        // how many teeth the player has
	IsDinosaur   bool       // indicates if the player species are dinosaurs or mammals
	Species      []*Species // the players' species
	Genes        int        // number of genes the player currently has
	AnimalTokens []int      // number of animal tokens (of silhouettes 0-3) are in the player's supply
}

// String formats a Player for display.
func (p *Player) String() string {
	return fmt.Sprintf("%v [%v/%v]", p.Name, p.Color, p.Dentition)
}
