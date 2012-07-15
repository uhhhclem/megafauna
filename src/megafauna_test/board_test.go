package megafauna_test

import (
	"megafauna"
	"testing"
)

func testNewBoard_CheckAdjacency(t *testing.T) {
	b := megafauna.NewBoard()
	var n, e, s, w = megafauna.MapDirectionN, megafauna.MapDirectionE, megafauna.MapDirectionS, megafauna.MapDirectionW

	t6 := b.HabitatMap["T6"]
	t2 := b.HabitatMap["T2"]
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

func testNewBoard_CheckClimaxNumbers(t *testing.T) {
	b := megafauna.NewBoard()

	check := func(key string, expected int) {
		if b.HabitatMap[key].ClimaxNumber != expected {
			t.Error("Incorrect climax number for %v", key)
		}
	}

	check("J1", 1)
	check("A2", 3)
	check("H5", 1)
	check("T0", 7)
	check("T5", 4)
}

func testNewBoard_CheckLatitudeMap(t *testing.T) {
	b := megafauna.NewBoard()

	a0 := b.HabitatMap["A0"]
	a5 := b.HabitatMap["A5"]
	t0 := b.HabitatMap["T0"]
	t7 := b.HabitatMap["T7"]

	if b.LatitudeMap["A"].Habitats[0] != a0 || b.LatitudeMap["A"].Habitats[5] != a5 {
		t.Error("Latitude map isn't set up properly.")
	}
	if b.LatitudeMap["T"].Habitats[0] != t0 || b.LatitudeMap["T"].Habitats[7] != t7 {
		t.Error("Latitude map isn't set up properly.")
	}
}

func testNewBoard_CheckIsOrogeny(t *testing.T) {
	b := megafauna.NewBoard()

	for k, v := range b.HabitatMap {
		if should_be_orogeny(k) != v.IsOrogeny {
			t.Errorf("IsOrogeny is wrong for key %v", k)
		}
	}
}

func should_be_orogeny(key string) bool {
	var orogeny_keys = []string{"A5", "J1", "J3", "H0", "H4", "T1"}
	for _, k := range orogeny_keys {
		if key == k {
			return true
		}
	}
	return false
}
