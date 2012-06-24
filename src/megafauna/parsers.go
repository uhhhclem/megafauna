package megafauna

import (
	"encoding/json"
	"io/ioutil"
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
