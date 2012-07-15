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

// Board is a slice of Latitude rows.  The rules call it "the map", but since map's a reserved word in Go,
// we'll be calling it the board. 
type Board struct
{
	Latitudes [][]*Habitat
	LatitudeMap map[string] *Latitude
}

// Latitude is one of the four latitude rows on the board (Tropic, Horse Latitudes, etc.)
type Latitude struct {
	Key     string
	Name    string
	Habitats []*Habitat
}


// Habitat is a space on the board that contains a slot for the biome, the predator triangle,
// and the rooter triangle.
type Habitat struct {
	Latitude         *Latitude
	ClimaxNumber     int
	IsOrogeny        bool
	AdjacentHabitats []*Habitat
	Biome            *Biome
}

func NewBoard() *Board {

	board := new(Board)
	board.Latitudes = make([][]*Habitat, 4)
	board.LatitudeMap = make(map[string] *Latitude)
	
	populateLatitudes(board)
	setOrogenyHabitats(board)
	setAdjacentHabitats(board)
	setLatitudeMap(board)
	
	return board			
}
 
func populateLatitudes(board *Board) {
	climaxNumbers := make([]string, 4)
	climaxNumbers[0] = "236541"
	climaxNumbers[1] = "614532"
	climaxNumbers[2] = "243561"
	climaxNumbers[3] = "73856412"
	
	for i, key := range latitudeKeys {
		board.Latitudes[i] = makeLatitude(string(key), climaxNumbers[i])
	}
	return
}

// makeLatitude creates a slice of Habitats from a string of climax numbers.
func makeLatitude(key string, climaxNumbers string) []*Habitat {
	lat := make([]*Habitat, len(climaxNumbers))
	for i, ch := range climaxNumbers {
		h := new(Habitat)
		h.AdjacentHabitats = make([]*Habitat, 4)
		h.ClimaxNumber, _ = strconv.Atoi(string(ch))
		lat[i] = h
	}
	return lat
}

// setOrogenyHabitats marks the Habitats on the board that are orogeny.
func setOrogenyHabitats(board *Board) {
	lat := board.Latitudes
	lat[0][4].IsOrogeny = true
	lat[1][1].IsOrogeny = true
	lat[1][3].IsOrogeny = true
	lat[2][0].IsOrogeny = true
	lat[2][4].IsOrogeny = true
	lat[3][1].IsOrogeny = true
	return
}

// setAdjacentHabitats builds the adjacency array for each habitat on the board.
func setAdjacentHabitats(board *Board) {
	var h *Habitat
	lat := board.Latitudes
	
	// for the top 4 rows of habitats, the method is simple, since the slices of habitats
	// comprise a rectangular array of cells.
	for row := 0; row < 4; row++ {
		for col := 0; col < 5; col++ {
			h = lat[row][col]
			if row > 0 {
				h.AdjacentHabitats[MapDirectionN] = lat[row-1][col]
			}
			if row < 3 {
				h.AdjacentHabitats[MapDirectionS] = lat[row+1][col]
			}
			if col > 0 {
				h.AdjacentHabitats[MapDirectionE] = lat[row][col-1]
			}
			if col < 5 {
				h.AdjacentHabitats[MapDirectionW] = lat[row][col+1]
			}
		}
	}
	
	// the last two habitats in the Tropics are special, because they're actually in their own row on the board.
	h = lat[3][6]
	h.AdjacentHabitats[MapDirectionN] = lat[3][2]
	h.AdjacentHabitats[MapDirectionE] = lat[3][7]
	lat[3][2].AdjacentHabitats[MapDirectionS] = h
	
	h = lat[3][7]
	h.AdjacentHabitats[MapDirectionW] = lat[3][6]
	h.AdjacentHabitats[MapDirectionN] = lat[3][3]
	lat[3][3].AdjacentHabitats[MapDirectionS] = h
	return
}

// setLatitudeMap initializes the LatitudeMap part of the board.
func setLatitudeMap(board *Board) {
	m := board.LatitudeMap
	
	for row, letter:= range latitudeKeys {
		key := string(letter)
		m[key] = new(Latitude)
		m[key].Key = key
		
		lat := m[string(key)]
		cols := len(lat.Habitats)
		for col := 0; col < cols; col++ {
			lat.Habitats[col] = board.Latitudes[row][col]
		}
	}
	
	m["A"].Name = "Arctic"
	m["J"].Name = "Jet Stream"
	m["H"].Name = "Horse Latitude"
	m["T"].Name = "Tropics"
	
	return
}