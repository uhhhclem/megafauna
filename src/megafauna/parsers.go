package megafauna

import (
	"encoding/csv"
	"io"
	"strconv"
)

// Parseable types implement a Parse method.  I'm not really sure that we need a Parseable interface per se; I
// can't really think of a useful function that would take a Parseable as an argument.  Maybe there could be 
// an Init(slices []Parseable, filenames []string) that parsed all of the data at startup. 
type Parseable interface {
	Parse(io.Reader) error
}

// BiomeMap is a Parseable of *Biomes.
type BiomeMap map[string]*Biome

// Indices into the CSV file for the Biome fields.
const (
	BiomeKeyField = iota
	BiomeTitleField
	BiomeSubtitleField
	BiomeClimaxNumberField
	BiomeRequirementsField
	BiomeRooterRequirementsField
	BiomeNicheField
	BiomeRedStarField
	BiomeBlueStarField
)

// Parse takes a Reader containing Biome data in CSV format, parses the data into Biomes, and populates
// the (pre-made) map with the biomes.
func (biomes BiomeMap) Parse(r io.Reader) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		b := new(Biome)
		b.Key = record[BiomeKeyField]
		b.Title = record[BiomeTitleField]
		b.Subtitle = record[BiomeSubtitleField]
		b.ClimaxNumber, err = strconv.Atoi(record[BiomeClimaxNumberField])
		if err != nil {
			return err
		}
		b.Requirements = MakeDNASpec(record[BiomeRequirementsField])
		b.RooterRequirements = MakeDNASpec(record[BiomeRooterRequirementsField])
		b.Niche, err = MakeNiche(record[BiomeNicheField])
		if err != nil {
			return err
		}
		b.RedStar, err = strconv.ParseBool(record[BiomeRedStarField])
		if err != nil {
			return nil
		}
		b.BlueStar, err = strconv.ParseBool(record[BiomeBlueStarField])
		if err != nil {
			return err
		}
		biomes[b.Key] = b
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