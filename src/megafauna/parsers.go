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
		niche, err := MakeNiche(record[5])
		if err != nil {
			return nil, err
		}
		RedStar, err := strconv.ParseBool(record[6])
		if err != nil {
			return nil, err
		}
		BlueStar, err := strconv.ParseBool(record[7])
		if err != nil {
			return nil, err
		}

		biomeArray[key] = Biome{
			Title:              record[0],
			Subtitle:           record[1],
			ClimaxNumber:       climaxNumber,
			Requirements:       MakeDNASpec(record[3]),
			Niche:              niche,
			RooterRequirements: MakeDNASpec(record[4]),
			RedStar:            RedStar,
			BlueStar:           BlueStar,
		}
	}
	return biomeArray, nil
}
