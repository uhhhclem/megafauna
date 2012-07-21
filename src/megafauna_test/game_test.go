package megafauna_test

import (
	"megafauna"
	"sort"
	"testing"
)

func testSortPlayers(t *testing.T) {
	s := make(megafauna.SortablePlayerCollection, 4)
	dentitions := []int{3, 2, 5, 4}
	for i, d := range dentitions {
		p := new(megafauna.Player)
		p.Dentition = d
		s[i] = p
	}
	sort.Sort(s)
	for i, p := range s {
		if p.Dentition != i+2 {
			t.Errorf("SortedPlayerCollection is not in order by dentition.")
			return
		}
	}
}
