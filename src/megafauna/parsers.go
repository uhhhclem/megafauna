package megafauna

import (
	"encoding/csv"
	"io"
)

// Parseable types implement a Parse method.  I'm not really sure that we need a Parseable interface per se; I
// can't really think of a useful function that would take a Parseable as an argument.  Maybe there could be 
// an Init(slices []Parseable, filenames []string) that parsed all of the data at startup. 
type Parseable interface {
	Parse(io.Reader) error
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
