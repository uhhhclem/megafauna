package megafauna

import (
	"bufio"
	"encoding/csv"
	"os"
)

func Parse(fileName string) (biomes []Biome, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	biomeArray := make([]Biome, len(records))
	for key, record := range records {
		biomeArray[key] = Biome{
			Title: record[0],
		}
	}
	return biomeArray, nil
}
