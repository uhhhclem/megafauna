package megafauna

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Parse(fileName string) (biomes []Biome, err error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var message []Biome
	err = json.Unmarshal(b, &message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func Emit(biomes []Biome) (err error) {
	b, err := json.MarshalIndent(biomes, "", "  ")
	if err != nil {
		return err
	}
	os.Stdout.Write(b)
	return nil
}