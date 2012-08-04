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

// Game is one discrete game of Bios Megafauna.
type Game struct {
	Players SortablePlayerCollection // slice of Player objects, in player order
	Board   *Board                   // the game's board
	//
	// card-related fields
	//
	Cards                map[string]*Card
	CardKeys             []string // shuffled slice of keys to all cards
	TriassicCardKeys     []string
	JurassicCardKeys     []string
	CretaceousCardKeys   []string
	TertiaryCardKeys     []string
	UpperDisplayCardKeys []string
	LowerDisplayCardKeys []string
	LowerDisplayGenes    []int
	//
	// tile-related fields
	//
	Tiles            map[string]*Tile
	MesozoicTileKeys []string // shuffled slice of keys to the Mesozoic tiles.
	CenozoicTileKeys []string // shuffled slice of keys to the Cenozoic tiles.
	TarpitTileKeys   []string // keys of the Tiles in the Tarpit.
}

// NewGame creates a new Game and initializes the Players.
func NewGame(names []string) (*Game, error) {
	var err error

	g := new(Game)
	g.Board = NewBoard()
	err = g.createCards()
	if err != nil {
		return nil, err
	}
	err = g.createTiles()
	if err != nil {
		return nil, err
	}
	g.createPlayers(names)
	if g.Players == nil {
		return nil, ErrInvalidPlayers
	}
	return g, nil
}

// createPlayers creates the Player objects in the game from a list of their names.  If it fails, Players will be nil.
func (g *Game) createPlayers(names []string) {
	if len(names) < 2 || len(names) > 4 {
		return
	}

	// Randomly order the dentitions, so that we can randomly assign one to each player.  
	// Shuffle takes strings, so we'll have to convert the dentitions to ints.  
	dentitions := []string{"2", "3", "4", "5"}
	Shuffle(dentitions)

	// create and name the players
	players := make(SortablePlayerCollection, len(names))
	for index, name := range names {
		dentition, _ := strconv.Atoi(dentitions[index])
		p := NewPlayer(name, dentition)
		players[index] = p
		g.Tiles[p.HomelandTile.Key] = p.HomelandTile
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

	// deal out the initial stacks
	g.TriassicCardKeys = g.dealCards(len(g.Players) * 3)
	g.JurassicCardKeys = g.dealCards(5)
	g.CretaceousCardKeys = g.dealCards(8)
	g.TertiaryCardKeys = g.dealCards(7)
	g.UpperDisplayCardKeys = g.dealCards(5)
	g.LowerDisplayCardKeys = g.dealCards(5)
	g.LowerDisplayGenes = make([]int, 5)
	return nil
}

// dealCards "deals" keys from CardKeys to the various stacks.
func (g *Game) dealCards(amount int) []string {
	stack := make([]string, 0)
	for i := 0; i < amount; i++ {
		key := g.CardKeys[0]
		stack = append(stack, key)
		// was that the last card in the deck?
		if len(g.CardKeys) == 1 {
			g.CardKeys = make([]string, 0)
			break
		}
		g.CardKeys = g.CardKeys[1:]
	}
	return stack
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
	Name             string             // the player's name
	Color            string             // the player's color
	Dentition        int                // how many teeth the player has
	IsDinosaur       bool               // indicates if the player species are dinosaurs or mammals
	Species          [][]*Animal        // the player's animals, indexed by silhouette (0-3)
	Genes            int                // number of genes the player currently has
	AnimalTokens     []int              // number of animal tokens (of silhouettes 0-3) are in the player's supply
	HomelandTile     *Tile              // the player's homeland tile
	InheritanceTiles []*InheritanceTile // the player's supply of unused inheritance tiles
}

// NewPlayer creates a new player given the player's name and dentition.
func NewPlayer(name string, dentition int) *Player {

	p := new(Player)

	p.Name = name
	p.Dentition = dentition

	colors := []string{"Red", "Orange", "Green", "White"}
	p.Color = colors[p.Dentition-2]
	p.IsDinosaur = p.Dentition == 2 || p.Dentition == 4
	p.Species = make([][]*Animal, 4)
	for i := 0; i < 3; i++ {
		p.Species[i] = make([]*Animal, 8)
	}
	p.AnimalTokens = []int{8, 8, 8, 8}
	p.InheritanceTiles = GetInheritanceTiles()

	t := new(Tile)
	b := new(BiomeTileData)
	t.BiomeData = b
	p.HomelandTile = t

	t.Key = p.Color
	t.HomelandPlayer = p
	t.IsLand = true

	// a homeland tile has no Requirements, just a Niche	
	niche, err := MakeNiche(strconv.Itoa(dentition))
	if err != nil {
		panic(err)
	}
	b.Niche = niche
	b.RedStar = true

	switch dentition {
	case 2:
		t.Title = "Giant Myriapods"
		t.LatitudeKey = "T"
		b.ClimaxNumber = 58
	case 3:
		t.Supertitle = "Alleghenian orogeny"
		t.Title = "Appalachian Cloud Forest"
		t.LatitudeKey = "O"
		b.IsOrogeny = true
		b.ClimaxNumber = 60
	case 4:
		t.Supertitle = "Caytoniales"
		t.Title = "Tree Fern-Cordaites Rainforest"
		t.LatitudeKey = "O"
		b.IsOrogeny = true
		b.ClimaxNumber = 55
	case 5:
		t.Supertitle = "Pteridophytes"
		t.Title = "Fern Understory"
		t.LatitudeKey = "A"
		b.ClimaxNumber = 75
	}

	return p
}

// String formats a Player for display.
func (p *Player) String() string {
	return fmt.Sprintf("%v [%v/%v]", p.Name, p.Color, p.Dentition)
}
