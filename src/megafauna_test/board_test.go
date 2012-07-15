package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestNewBoard_CheckHabitats(t *testing.T) {
	b := megafauna.NewBoard()
	if len(b.Habitats) != 4 {
		t.Error("Expected Habitats to have 4 elements.")
		return
	}
	ha := b.Habitats[0]
	if len(ha) != 6 {
		t.Error("Expected the Arctic to have 6 habitats.")
	}
	if ha[0] == nil {
		t.Error("No habitat in Arctic column 0.")
		return
	}
	if ha[0].Key != "A0" {
		t.Errorf("Expected habitat A0's Key to be A0, but it's %v", ha[0].Key)
	}
	ht := b.Habitats[3]
	if len(ht) != 8 {
		t.Error("Expected the Tropics to have 8 habitats.")
	}

	for row, _ := range b.Habitats {
		for col, _ := range b.Habitats[row] {
			if b.Habitats[row][col] == nil {
				t.Error("b.Habitats[%v][%v] didn't get initialized.", row, col)
			}
		}
	}

}

func TestNewBoard_CheckLatitudeMap(t *testing.T) {
	b := megafauna.NewBoard()

	lat := b.LatitudeMap["A"]
	if lat == nil {
		t.Error("LatitudeMap doesn't contain the arctic.")
		return
	}
	if len(lat.Habitats) != len(b.Habitats[0]) {
		t.Error("LatitudeMap.Habitats isn't the right length.")
		return
	}
	for col, h := range lat.Habitats {
		if h != b.Habitats[0][col] {
			t.Error("LatitudeMap.Habitats doesn't contain expected habitats.")
		}
	}

}

func testNewBoard_CheckHabitatMap(t *testing.T) {
	b := megafauna.NewBoard()

	keys := []string{"A0", "J3", "H5", "T7"}
	for _, key := range keys {
		h := b.HabitatMap[key]
		if h == nil {
			t.Errorf("HabitatMap doesn't contain %v.", key)
			return
		}
		if h.Key != key {
			t.Errorf("HabitatMap[%v] has a key of %v.", key, h.Key)
		}
	}
}

func TestNewBoard_CheckAdjacency(t *testing.T) {
	b := megafauna.NewBoard()
	var n, e, s, w = megafauna.MapDirectionN, megafauna.MapDirectionE, megafauna.MapDirectionS, megafauna.MapDirectionW

	t6 := b.HabitatMap["T6"]
	t2 := b.HabitatMap["T2"]
	if t6 == nil || t2 == nil {
		t.Error("HabitatMap lookup failed.")
		return
	}
	if t6.AdjacentHabitats[n] != t2 {
		t.Error("T2 is not north of T6")
	}
	if t2.AdjacentHabitats[s] != t6 {
		t.Error("T6 is not south of T2.")
	}

	a0 := b.HabitatMap["A0"]
	j0 := b.HabitatMap["J0"]
	a1 := b.HabitatMap["A1"]

	if a0.AdjacentHabitats[w] != nil || a0.AdjacentHabitats[n] != nil ||
		a0.AdjacentHabitats[e] != a1 || a0.AdjacentHabitats[s] != j0 {
		t.Error("A0 doesn't have the right neighbors.")
	}

	j4 := b.HabitatMap["J4"]

	if j4.AdjacentHabitats[w] == nil || j4.AdjacentHabitats[n] == nil ||
		j4.AdjacentHabitats[e] == nil || j4.AdjacentHabitats[s] == nil {
		t.Error("J4 doesn't have four neighbors assigned.")
	}

}

func TestNewBoard_CheckClimaxNumbers(t *testing.T) {
	b := megafauna.NewBoard()

	check := func(key string, expected int) {
		h := b.HabitatMap[key]
		if h == nil {
			t.Error("HabitatMap lookup failed.")
			return
		}
		if b.HabitatMap[key].ClimaxNumber != expected {
			t.Errorf("Incorrect climax number for %v", key)
		}
	}

	check("J1", 1)
	check("A2", 3)
	check("H5", 1)
	check("T0", 7)
	check("T5", 4)
}

func TestNewBoard_CheckOrogeny(t *testing.T) {
	b := megafauna.NewBoard()

	for k, v := range b.HabitatMap {
		if should_be_orogeny(k) != v.IsOrogeny {
			t.Errorf("IsOrogeny is wrong for key %v", k)
		}
	}

	if len(b.LatitudeMap["O"].Habitats) != 6 {
		t.Error("LatitudeMap for orogeny isn't initialized.")
		return
	}
	for _, hab := range b.LatitudeMap["O"].Habitats {
		if !should_be_orogeny(hab.Key) {
			t.Error("We have a non-orogeny habitat in the LatitudeMap, somehow.")
		}
	}
}

func should_be_orogeny(key string) bool {
	var orogeny_keys = []string{"A4", "J1", "J3", "H0", "H4", "T1"}
	for _, k := range orogeny_keys {
		if key == k {
			return true
		}
	}
	return false
}

func TestFindLowestClimax(t *testing.T) {
	b := megafauna.NewBoard()
	h, err := b.FindLowestClimax("bogus")
	if err != megafauna.ErrInvalidLatitudeKey {
		t.Error("Expected an ErrInvalidLatitudeKey here.")
		return
	}
	h, err = b.FindLowestClimax("T")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if h == nil {
		t.Error("FindLowestClimax returned nil.")
		return
	}
	if h.Key != "T6" {
		t.Errorf("Expected to find T6, but found %v.", h.Key)
	}
}
