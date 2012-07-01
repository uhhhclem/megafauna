package megafauna

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

// Parseable types implement a Parse method.  I'm not really sure that we need a Parseable interface per se; I
// can't really think of a useful function that would take a Parseable as an argument.  Maybe there could be 
// an Init(slices []Parseable, filenames []string) that parsed all of the data at startup. 
type Parseable interface {
	Parse(io.Reader) error
}

// TileMap is a Parseable of *BiomeTiles.
type BiomeTileMap map[string]*BiomeTile

// Indices into the CSV file for the BiomeTile fields.
const (
	BiomeTileKeyField = iota
	BiomeTileLatitudeKeyField
	BiomeTileTypeField
	BiomeTileTitleField
	BiomeTileSubtitleField
	BiomeTileClimaxNumberField
	BiomeTileRequirementsField
	BiomeTileRooterRequirementsField
	BiomeTileNicheField
	BiomeTileRedStarField
	BiomeTileBlueStarField
)

// BiomeTileParseError is the error returned if one of the key fields in the biome data is invalid.
type BiomeTileParseError struct {
	InvalidTypeKey     string
	InvalidLatitudeKey string
}

func (e *BiomeTileParseError) Error() string {
	if e.InvalidTypeKey != "" {
		return "Type must be one of " + BiomeTypeKeys + " but is " + e.InvalidTypeKey
	}
	return "LatitudeKey must be one of " + LatitudeKeys + " but is " + e.InvalidLatitudeKey

}

// Parse takes a Reader containing BiomeTile data in CSV format, parses the data into BiomeTiles, and populates
// the (pre-made) map with the BiomeTiles.
func (tiles BiomeTileMap) Parse(r io.Reader) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		b := new(BiomeTile)
		b.Key = record[BiomeTileKeyField]
		b.LatitudeKey = record[BiomeTileLatitudeKeyField]

		if len(b.LatitudeKey) != 1 || !strings.Contains(LatitudeKeys, b.LatitudeKey) {
			return &BiomeTileParseError{"", b.LatitudeKey}
		}

		switch record[BiomeTileTypeField] {
		case "L":
			b.IsLand = true
		case "W":
			b.IsWater = true
		case "O":
			{
				b.IsLand = true
				b.IsOrogeny = true
			}
		default:
			{
				return &BiomeTileParseError{record[BiomeTileTypeField], ""}
			}
		}

		b.Title = record[BiomeTileTitleField]
		b.Subtitle = record[BiomeTileSubtitleField]
		b.ClimaxNumber, err = strconv.Atoi(record[BiomeTileClimaxNumberField])
		if err != nil {
			return err
		}
		b.Requirements = MakeDNASpec(record[BiomeTileRequirementsField])
		b.RooterRequirements = MakeDNASpec(record[BiomeTileRooterRequirementsField])
		b.Niche, err = MakeNiche(record[BiomeTileNicheField])
		if err != nil {
			return err
		}
		b.RedStar, err = strconv.ParseBool(record[BiomeTileRedStarField])
		if err != nil {
			return nil
		}
		b.BlueStar, err = strconv.ParseBool(record[BiomeTileBlueStarField])
		if err != nil {
			return err
		}
		tiles[b.Key] = b
	}
	return nil
}

// LatitudeMap is a Parseable of Latitudes.
type LatitudeMap map[string]*Latitude

const (
	LatitudeKeyField = iota
	LatitudeNameField
)

// Parse parses the latitude data in r into a map of Latitude objects.
func (latitudes LatitudeMap) Parse(r io.Reader) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		l := new(Latitude)
		l.Key = record[LatitudeKeyField]
		l.Name = record[LatitudeNameField]
		latitudes[l.Key] = l
	}
	return nil
}
