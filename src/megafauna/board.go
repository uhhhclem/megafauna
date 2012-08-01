package megafauna

import (
	"fmt"
	"strconv"
)

const (
	MapDirectionN = iota
	MapDirectionE
	MapDirectionS
	MapDirectionW
)

// latitudeKeys is local, because megafauna.LatitudeKeys includes the O for orogeny. 
const latitudeKeys = "AJHT"

// Board contains the Habitats on the board as well as various data structures to support lookup and searching.
// LatitudeMap is used to find all the habitats with a given latitude key, and HabitatMap is used to find a
// habitat given its key.
type Board struct {
	Habitats    [][]*Habitat         // the four outer slices are the latitude bands
	LatitudeMap map[string]*Latitude // map of latitude keys to the habitats in that latitude
	HabitatMap  map[string]*Habitat  // for looking up habitats by their key
}

// Latitude is one of the four latitude rows on the board (Tropic, Horse Latitudes, etc.).  Note that there will
// also be an orogeny "latitude" with a Key of O - this facilitates using FindLowestClimax to place orogeny biomes.
type Latitude struct {
	Key      string
	Name     string
	Habitats []*Habitat
}

// Habitat is a space on the board that contains a slot for the biome, the predator triangle,
// and the rooter triangle.
type Habitat struct {
	Key              string     // The latitude key plus a zero-based index, e.g. "T0" is the first habitat in the Tropics
	ClimaxNumber     int        // The printed climax number for the habitat
	IsOrogeny        bool       // True for orogeny habitats 
	AdjacentHabitats []*Habitat // Indexed by map direction; an entry is nil if there's no habitat in that direction
	Biome            *Biome     // The Biome, if any in this habitat
}

// NewBoard initializes the board.
func NewBoard() *Board {

	board := new(Board)
	board.Habitats = make([][]*Habitat, 4)
	board.LatitudeMap = make(map[string]*Latitude)

	populateHabitats(board)
	setOrogenyHabitats(board)
	setAdjacentHabitats(board)
	setLatitudeMap(board)
	setHabitatMap(board)

	return board
}

// populateHabitats creates all the Habitat objects and stores them in the board's Habitats field.
func populateHabitats(board *Board) {
	climaxNumbers := make([]string, 4)
	climaxNumbers[0] = "263541"
	climaxNumbers[1] = "614532"
	climaxNumbers[2] = "243561"
	climaxNumbers[3] = "73856412"

	for row, ch := range latitudeKeys {
		latitudeKey := string(ch)
		board.Habitats[row] = makeHabitatsInLatitude(latitudeKey, climaxNumbers[row])
	}
	return
}

// makeHabitatsInLatitude creates a slice of Habitats from a string of climax numbers.
func makeHabitatsInLatitude(latitudeKey string, climaxNumbers string) []*Habitat {
	lat := make([]*Habitat, len(climaxNumbers))
	for col, ch := range climaxNumbers {
		h := new(Habitat)
		h.AdjacentHabitats = make([]*Habitat, 4)
		h.ClimaxNumber, _ = strconv.Atoi(string(ch))
		h.Key = fmt.Sprintf("%v%v", latitudeKey, col)
		lat[col] = h
	}
	return lat
}

// setOrogenyHabitats marks the Habitats on the board that are orogeny.  It also creates the LatitudeMap for orogeny habitats.
func setOrogenyHabitats(board *Board) {
	hab := board.Habitats
	hab[0][4].IsOrogeny = true
	hab[1][1].IsOrogeny = true
	hab[1][3].IsOrogeny = true
	hab[2][0].IsOrogeny = true
	hab[2][4].IsOrogeny = true
	hab[3][1].IsOrogeny = true

	lat := new(Latitude)
	lat.Key = "O"
	lat.Name = "Orogeny"
	lat.Habitats = make([]*Habitat, 6)

	index := 0
	for row, _ := range board.Habitats {
		for col, _ := range board.Habitats[row] {
			h := board.Habitats[row][col]
			if h.IsOrogeny {
				lat.Habitats[index] = h
				index++
			}
		}
	}

	board.LatitudeMap["O"] = lat

	return
}

// setAdjacentHabitats builds the adjacency array for each habitat on the board.
func setAdjacentHabitats(board *Board) {
	var h *Habitat
	hab := board.Habitats

	// for the top 4 rows of habitats, the method is simple, since the slices of habitats
	// comprise a rectangular array of cells.
	for row := 0; row < len(hab); row++ {
		for col := 0; col < 6; col++ {
			h = hab[row][col]
			if row > 0 {
				h.AdjacentHabitats[MapDirectionN] = hab[row-1][col]
			}
			if row < 3 {
				h.AdjacentHabitats[MapDirectionS] = hab[row+1][col]
			}
			if col > 0 {
				h.AdjacentHabitats[MapDirectionW] = hab[row][col-1]
			}
			if col < 5 {
				h.AdjacentHabitats[MapDirectionE] = hab[row][col+1]
			}
		}
	}

	// the last two habitats in the Tropics are special, because they're actually in their own row on the board.
	h = hab[3][6]
	h.AdjacentHabitats[MapDirectionN] = hab[3][2]
	h.AdjacentHabitats[MapDirectionE] = hab[3][7]
	hab[3][2].AdjacentHabitats[MapDirectionS] = h

	h = hab[3][7]
	h.AdjacentHabitats[MapDirectionW] = hab[3][6]
	h.AdjacentHabitats[MapDirectionN] = hab[3][3]
	hab[3][3].AdjacentHabitats[MapDirectionS] = h
	return
}

// setLatitudeMap initializes the LatitudeMap part of the board.
func setLatitudeMap(board *Board) {
	m := board.LatitudeMap
	for row, letter := range latitudeKeys {
		key := string(letter)
		lat := new(Latitude)
		lat.Key = key
		cols := 6
		if key == "T" {
			cols = 8
		}
		lat.Habitats = make([]*Habitat, cols)
		for col := 0; col < cols; col++ {
			lat.Habitats[col] = board.Habitats[row][col]
		}
		m[key] = lat
	}

	m["A"].Name = "Arctic"
	m["J"].Name = "Jet Stream"
	m["H"].Name = "Horse Latitude"
	m["T"].Name = "Tropics"

	return
}

// setHabitatMap puts all of the Habitats on the board into a map keyed on Habitat.Key.
func setHabitatMap(board *Board) {
	board.HabitatMap = make(map[string]*Habitat)
	for row, _ := range board.Habitats {
		for col, _ := range board.Habitats[row] {
			h := board.Habitats[row][col]
			board.HabitatMap[h.Key] = h
		}
	}
}

// FindLowestClimax returns the habitat with the lowest climax number in the requested latitude.
func (b *Board) FindLowestClimax(latitudeKey string) (*Habitat, error) {
	var result *Habitat
	min := 1000 // all climax numbers are less than 1000
	lat := b.LatitudeMap[latitudeKey]
	if lat == nil {
		return nil, ErrInvalidLatitudeKey
	}
	for _, h := range b.LatitudeMap[latitudeKey].Habitats {
		climaxNumber := h.ClimaxNumber
		if h.Biome != nil {
			climaxNumber = h.Biome.GetClimaxNumber()
		}
		if climaxNumber < min {
			min = climaxNumber
			result = h
		}
	}
	return result, nil
}

// PlaceTileOnBoard places a tile on the board, and can return a displaced tile.
// PlaceTileOnBoard looks at the tile's Latitude, and then finds the Habitats for that Latitude, 
// and among those chooses the one that has the lowest climax number.
func (b *Board) PlaceTileOnBoard(t *Tile) (*Tile, error) {
	climax, err := b.FindLowestClimax(t.LatitudeKey)
	if err == nil {
		return nil, ErrInvalidLatitudeKey
	}
	if t.IsBiomeTile() {
		displacedTile := climax.Biome.Tile
		// swap the existing biome for the new one
		return displacedTile, nil
	}
	return nil, nil
}
