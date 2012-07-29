package megafauna

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

var (
	ErrCardNotFound   = errors.New("Card not found.")
	ErrInvalidPlayers = errors.New("Invalid player list.")
)

// Game is one discrete game of Bios Megafauna.
type Game struct {
	Players          SortablePlayerCollection // slice of Player objects, in player order
	Cards            map[string]*Card
	CardKeys         []string // shuffled slice of keys to all cards
	Tiles            map[string]*Tile
	MesozoicTileKeys []string // shuffled slice of keys to the Mesozoic tiles.
	CenozoicTileKeys []string // shuffled slice of keys to the Cenozoic tiles.
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
func NewGame(names []string) (*Game, error) {
	var err error

	g := new(Game)
	g.createPlayers(names)
	if g.Players == nil {
		return nil, ErrInvalidPlayers
	}
	err = g.createCards()
	if err != nil {
		return nil, err
	}
	err = g.createTiles()
	if err != nil {
		return nil, err
	}
	return g, nil
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
func (g *Game) createCards() error {
	var err error

	// get the mutation cards
	g.Cards, err = GetCards()
	if err != nil {
		return err
	}

	// get their keys and shuffle them
	g.CardKeys = make([]string, 0)
	for k, _ := range g.Cards {
		g.CardKeys = append(g.CardKeys, k)
	}
	Shuffle(g.CardKeys)
	return nil
}

// createTiles initializes the tile stacks.
func (g *Game) createTiles() error {
	var err error

	// get the tiles
	g.Tiles, err = GetTiles()
	if err != nil {
		return err
	}

	// get the keys for the two stacks of tiles and shuffle them
	g.MesozoicTileKeys = make([]string, 0)
	g.CenozoicTileKeys = make([]string, 0)
	for k, t := range g.Tiles {
		if t.IsMesozoic {
			g.MesozoicTileKeys = append(g.MesozoicTileKeys, k)
		} else {
			g.CenozoicTileKeys = append(g.CenozoicTileKeys, k)
		}
	}
	Shuffle(g.MesozoicTileKeys)
	Shuffle(g.CenozoicTileKeys)
	return nil
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
