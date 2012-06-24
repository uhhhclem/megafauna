package megafauna

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
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
		climaxNumber, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
		biomeArray[key] = Biome{
			Title:        record[0],
			Subtitle:     record[1],
			ClimaxNumber: climaxNumber,
			Requirements: MakeDNASpec(record[3]),
		}
	}
	return biomeArray, nil
}
