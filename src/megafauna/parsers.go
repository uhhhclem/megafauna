package megafauna

import (
	"encoding/json"
	"io/ioutil"
)

func parse(fileName string) (biomes []Biome, err error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var message []Biome
	err2 := json.Unmarshal(b, &message)
	if err2 != nil {
		return nil, err
	}
	return message, nil
}
