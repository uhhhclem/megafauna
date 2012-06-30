package megafauna

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
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
		climaxNumber, err := strconv.Atoi(record[2])
		if err != nil {
			return err
		}
		niche, err := MakeNiche(record[4])
		if err != nil {
			return err
		}
		b := new(Biome)
		b.Title = record[0]
		b.Subtitle = record[1]
		b.ClimaxNumber = climaxNumber
		b.Requirements = MakeDNASpec(record[3])
		b.Niche = niche
		biomes[i] = b
	}
	return nil
}

