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

// BiomeSlice is a Parseable of *Biomes.
type BiomeSlice []*Biome

// Parse takes a Reader containing Biome data in CSV format, parses the data into Biomes, and populates
// the (pre-allocated) slice with the biomes.
func (biomes BiomeSlice) Parse(r io.Reader) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for i, record := range records {
		b := new(Biome)
		b.Title = record[0]
		b.Subtitle = record[1]
		b.ClimaxNumber, err = strconv.Atoi(record[2])
		if err != nil {
			return err
		}
		b.Requirements = MakeDNASpec(record[3])
		b.RooterRequirements = MakeDNASpec(record[4])
		b.Niche, err = MakeNiche(record[5])
		if err != nil {
			return err
		}
		b.RedStar, err = strconv.ParseBool(record[6])
		if err != nil {
			return nil
		}
		b.BlueStar, err = strconv.ParseBool(record[7])
		if err != nil {
			return err
		}

		biomes[i] = b
	}
	return nil
}

