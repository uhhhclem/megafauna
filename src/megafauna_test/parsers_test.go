package megafauna_test

import (
	"megafauna"
	"strings"
	"testing"
)

func TestLatitudeMapParse(t *testing.T) {
	data := "1,Tropics\n2,Horse Latitudes\n3,Jet Stream\n4,Arctic"

	latitudes := make(megafauna.LatitudeMap)

	reader := strings.NewReader(data)
	err := latitudes.Parse(reader)
	if err != nil {
		t.Errorf("An error occurred: %v", err)
	}
	if len(latitudes) == 0 {
		t.Errorf("Expected some latitudes.")
	}

	key := "2"
	expected := "Horse Latitudes"
	if latitudes[key].Name != expected {
		t.Errorf("Expected Latitude with a key of %v to be %v, but it's %v.", key, expected, latitudes[key].Name)
	}
}
