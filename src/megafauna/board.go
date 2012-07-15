package megafauna

import (
	"strconv"
)

const (
	MapDirectionN = iota
	MapDirectionE
	MapDirectionS
	MapDirectionW
)

const latitudeKeys = "AJHT"

// Board contains the Habitats on the board.
type Board struct {
	Habitats    [][]*Habitat         // the four outer slices are the latitude bands
	LatitudeMap map[string]*Latitude // map of latitude keys to the habitats in that latitude
	HabitatMap  map[string]*Habitat  // for looking up habitats by their key
}

// Latitude is one of the four latitude rows on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Key      string
	Name     string
	Habitats []*Habitat
}

// Habitat is a space on the board that contains a slot for the biome, the predator triangle,
// and the rooter triangle.
type Habitat struct {
	Key              string
	Latitude         *Latitude
	ClimaxNumber     int
	IsOrogeny        bool
	AdjacentHabitats []*Habitat
	Biome            *Biome
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

func populateHabitats(board *Board) {
	climaxNumbers := make([]string, 4)
	climaxNumbers[0] = "236541"
	climaxNumbers[1] = "614532"
	climaxNumbers[2] = "243561"
	climaxNumbers[3] = "73856412"

	for i, key := range latitudeKeys {
		board.Habitats[i] = makeHabitatsInLatitude(string(key), climaxNumbers[i])
	}
	return
}

// makeHabitatsInLatitude creates a slice of Habitats from a string of climax numbers.
func makeHabitatsInLatitude(key string, climaxNumbers string) []*Habitat {
	lat := make([]*Habitat, len(climaxNumbers))
	for i, ch := range climaxNumbers {
		h := new(Habitat)
		h.AdjacentHabitats = make([]*Habitat, 4)
		h.Key = key + string(i)
		h.ClimaxNumber, _ = strconv.Atoi(string(ch))
		lat[i] = h
	}
	return lat
}

// setOrogenyHabitats marks the Habitats on the board that are orogeny.
func setOrogenyHabitats(board *Board) {
	hab := board.Habitats
	hab[0][4].IsOrogeny = true
	hab[1][1].IsOrogeny = true
	hab[1][3].IsOrogeny = true
	hab[2][0].IsOrogeny = true
	hab[2][4].IsOrogeny = true
	hab[3][1].IsOrogeny = true
	return
}

// setAdjacentHabitats builds the adjacency array for each habitat on the board.
func setAdjacentHabitats(board *Board) {
	var h *Habitat
	hab := board.Habitats

	// for the top 4 rows of habitats, the method is simple, since the slices of habitats
	// comprise a rectangular array of cells.
	for row := 0; row < 4; row++ {
		for col := 0; col < 5; col++ {
			h = hab[row][col]
			if row > 0 {
				h.AdjacentHabitats[MapDirectionN] = hab[row-1][col]
			}
			if row < 3 {
				h.AdjacentHabitats[MapDirectionS] = hab[row+1][col]
			}
			if col > 0 {
				h.AdjacentHabitats[MapDirectionE] = hab[row][col-1]
			}
			if col < 5 {
				h.AdjacentHabitats[MapDirectionW] = hab[row][col+1]
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
		m[key] = new(Latitude)
		m[key].Key = key

		lat := m[string(key)]
		cols := len(lat.Habitats)
		for col := 0; col < cols; col++ {
			lat.Habitats[col] = board.Habitats[row][col]
		}
	}

	m["A"].Name = "Arctic"
	m["J"].Name = "Jet Stream"
	m["H"].Name = "Horse Latitude"
	m["T"].Name = "Tropics"

	return
}

// setHabitatMap puts all of the Habitats on the board into a map keyed on Habitat.Key, e.g. the leftmost Habitat
// in the horse latitudes is H0 and the rightmost arctic habitat is A5.
func setHabitatMap(board *Board) {
	for key, lat := range board.LatitudeMap {
		for i, h := range lat.Habitats {
			h.Key = key + string(i)
		}
	}
}
