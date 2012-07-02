package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestMakeEvent(t *testing.T) {
	var e *megafauna.Event
	e = megafauna.MakeEvent("T", "", 0, false)
	if !e.IsDrawTwo {
		t.Error("Expected event to be a draw-two event, and it isn't.")
	}

	e = megafauna.MakeEvent("C", "", 5, true)
	if !e.IsCatastrophe {
		t.Error("Expected event to be a catastrophe, and it isn't.")
	}
	if e.CatastropheLevel != 5 {
		t.Errorf("Expected catastrophe level to be 5, and it's %v.", e.CatastropheLevel)
	}
	if !e.CatastropheIsWarming {
		t.Error("Expected catastrophe to be warming, and it's cooling.")
	}

	e = megafauna.MakeEvent("MP", "HA", 0, false)
	if !e.IsMilankovich {
		t.Error("Expected event to be a Milankovich event, and it isn't.")
	}
	keys := e.MilankovichLatitudeKeys
	if len(keys) != 2 || keys[0] != "H" || keys[1] != "A" {
		t.Errorf("Milankovich latitudes didn't parse; they're %v.", e.MilankovichLatitudeKeys)
	}

}
