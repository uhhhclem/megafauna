package megafauna_test

import (
	"megafauna"
	"sort"
	"testing"
)

func TestSortPlayers(t *testing.T) {
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

func TestNewGame_Players(t *testing.T) {
	var g *megafauna.Game

	g = megafauna.NewGame([]string{""})
	if g != nil {
		t.Error("Shouldn't create a game if there aren't enough names.")
		return
	}
	g = megafauna.NewGame([]string{"A", "B", "C", "D", "E"})
	if g != nil {
		t.Error("Shouldn't create a game if there aren't too many names.")
		return
	}

	g = megafauna.NewGame([]string{"Matthew", "Mark", "Luke", "John"})
	if len(g.Players) != 4 {
		t.Error("g.Players isn't the right length.")
		return
	}

	var ok bool
	found := make([]*megafauna.Player, 0)
	for i, p := range g.Players {
		switch p.Color {
		default:
			t.Errorf("Unexpected player color: %v", p.Color)
			return
		case "Red":
			found = append(found, p)
			ok = p.Dentition == 2 && p.IsDinosaur && i == 0 && p.Genes == 3
		case "Orange":
			found = append(found, p)
			ok = p.Dentition == 3 && !p.IsDinosaur && i == 1 && p.Genes == 4
		case "Green":
			found = append(found, p)
			ok = p.Dentition == 4 && p.IsDinosaur && i == 2 && p.Genes == 4
		case "White":
			ok = p.Dentition == 5 && !p.IsDinosaur && i == 3 && p.Genes == 4
			found = append(found, p)
		}
		if !ok {
			t.Errorf("There's something wrong with the %v player.", p.Color)
			return
		}
	}
	if len(found) != 4 {
		t.Error("Didn't find all 4 colors, somehow.")
	}

	g = megafauna.NewGame([]string{"Tinker", "Evers", "Chance"})
	var prevDentition int
	for _, p := range g.Players {
		if p.Dentition < prevDentition {
			t.Errorf("Players aren't sorted by dentition.")
			return
		}
		prevDentition = p.Dentition
	}
}

func TestGetCard(t *testing.T) {
	g := megafauna.NewGame([]string{"Dick", "Jane"})

	if len(g.CardKeys) < 10 {
		t.Error("There should be at least 10 cards.")
		return
	}
	for _, k := range g.CardKeys {
		mut, gen := g.GetCard(k)
		if mut == nil && gen == nil {
			t.Error("We should always get either a mutation or a genotype card.")
			return
		}
		if mut != nil && gen != nil {
			t.Error("Only one of the objects returned should have a value.")
			return
		}
	}
}

func TestGetPlayer(t *testing.T) {
	g := megafauna.NewGame([]string{"John", "Paul", "George", "Ringo"})
	for _, d := range []int{2, 3, 4, 5} {
		p := g.GetPlayer(d)
		if p == nil || p.Dentition != d {
			t.Errorf("Didn't get the player with dentition %v.", d)
			return
		}
	}
}
