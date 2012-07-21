package megafauna

import (
	"fmt"
	"strconv"
)

type Game struct {
	Players SortablePlayerCollection // slice of Player objects, in player order
}

// SortablePlayerCollection is used to sort players by Dentition
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

// NewGame creates a new Game and initializes the Players
func NewGame(players SortablePlayerCollection) *Game {
	if len(players) < 2 || len(players) > 4 {
		return nil
	}
	g := new(Game)
	// Randomly assign dentitions to each player.  Shuffle takes strings, so we'll have to convert
	// the dentitions to ints.  
	dentitions := []string{"2", "3", "4", "5"}
	colors := []string{"Red", "Orange", "Green", "White"}
	Shuffle(dentitions)
	minDentition := 99
	for index, p := range players {
		p.Dentition, _ = strconv.Atoi(dentitions[index])
		if p.Dentition < minDentition {
			minDentition = p.Dentition
		}
		p.Color = colors[p.Dentition-2]
		p.IsDinosaur = p.Dentition == 2 || p.Dentition == 4
		p.Species = make([]*Species, 4)
	}

	// Everyone gets 4 genes except the starting player
	for _, p := range players {
		p.Genes = 4
		if p.Dentition == minDentition {
			p.Genes -= 1
		}
	}

	return g
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

type Player struct {
	Name       string     // the player's name
	Color      string     // the player's color
	Dentition  int        // how many teeth the player has
	IsDinosaur bool       // indicates if the player species are dinosaurs or mammals
	Species    []*Species // the players' species
	Genes      int        // how many genes the player currently has
}

func (p *Player) String() string {
	return fmt.Sprintf("%v [%v/%v]", p.Name, p.Color, p.Dentition)
}
