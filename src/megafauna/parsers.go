package megafauna

import (
	"encoding/csv"
	"errors"
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

// TileMap is a Parseable of Tiles
type TileMap map[string]*Tile

// Indices into the CSV file for the biome tile fields.
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

var (
	errInvalidType        = errors.New("Type must be one of " + BiomeTypeKeys + ".")
	errInvalidLatitudeKey = errors.New("LatitudeKey must be one of " + LatitudeKeys + ".")
)

// Parse takes a Reader containing Tile data in CSV format, parses the data into Tiles, and populates
// the (pre-made) map with the Tiles. 
func (tiles TileMap) Parse(r io.Reader) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		t := new(Tile)

		// first get fields common to all Tiles.
		t.Key = record[BiomeTileKeyField]
		t.LatitudeKey = record[BiomeTileLatitudeKeyField]

		if len(t.LatitudeKey) != 1 || !strings.Contains(LatitudeKeys, t.LatitudeKey) {
			return errInvalidLatitudeKey
		}

		t.Title = record[BiomeTileTitleField]
		t.Subtitle = record[BiomeTileSubtitleField]

		// now parse out the BiomeTileData.
		t.BiomeData = new(BiomeTileData)
		b := t.BiomeData

		switch record[BiomeTileTypeField] {
		case "L":
			t.IsLand = true
		case "S":
			t.IsSea = true
		case "O":
			t.IsLand = true
			b.IsOrogeny = true
		case "B":
			t.IsLand = true
			t.IsSea = true
		default:
			return errInvalidType
		}

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
		tiles[t.Key] = t
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
