package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestDNA(t *testing.T) {
	var dna megafauna.DNA
	dna.Letter = "A"
	if !dna.IsRoadrunner() {
		t.Errorf("Expected A to be roadrunner.")
	}
	if dna.IsDietary() {
		t.Errorf("Expected A to not be dietary.")
	}
	dna.Letter = "P"
	if dna.IsRoadrunner() {
		t.Errorf("Expected P to not be roadrunner.")
	}
	if !dna.IsDietary() {
		t.Errorf("Expected P to be dietary.")
	}

}

func TestMakeDNASpec(t *testing.T) {
	spec := megafauna.MakeDNASpec("BBGAAA")
	if spec.Spec != "BBGAAA" {
		t.Errorf("Spec set incorrectly - shouldn't be %v", spec.Spec)
	}

	if spec.Breakdown["B"].Value != 2 {
		t.Errorf("Breakdown of B shouldn't be %v.", spec.Breakdown["B"].Value)
	}
	if spec.Breakdown["G"].Value != 1 {
		t.Errorf("Breakdown of G shouldn't be %v.", spec.Breakdown["G"].Value)
	}
	if spec.Breakdown["A"].Value != 3 {
		t.Errorf("Breakdown of A shouldn't be %v.", spec.Breakdown["A"].Value)
	}
}

func TestCanPreyOn(t *testing.T) {
	predators := []string{"BGA", "M", "MM", "INS"}

	prey := []string{"IG", "MGG", "MP", "HHS"}
	for i, _ := range predators {
		predatorSpec := megafauna.MakeDNASpec(predators[i])
		preySpec := megafauna.MakeDNASpec(prey[i])
		if !predatorSpec.CanPreyOn(preySpec) {
			t.Errorf("%v should be able to prey on %v", predatorSpec.Spec, preySpec.Spec)
		}
	}
	prey = []string{"IGAA", "MSGG", "MN", "A"}
	for i, _ := range predators {
		predatorSpec := megafauna.MakeDNASpec(predators[i])
		preySpec := megafauna.MakeDNASpec(prey[i])
		if predatorSpec.CanPreyOn(preySpec) {
			t.Errorf("%v shouldn't be able to prey on %v", predatorSpec.Spec, preySpec.Spec)
		}
	}

}

func TestCanFeedOn(t *testing.T) {
	eaters := []string{"BGA", "M", "MMP", "INS"}

	food := []string{"BG", "", "MP", "IAA"}
	for i, _ := range eaters {
		eaterSpec := megafauna.MakeDNASpec(eaters[i])
		foodSpec := megafauna.MakeDNASpec(food[i])
		if !eaterSpec.CanFeedOn(foodSpec) {
			t.Errorf("%v should be able to feed on %v", eaterSpec.Spec, foodSpec.Spec)
		}
	}

	food = []string{"BGG", "H", "MPP", "IIN"}
	for i, _ := range eaters {
		eaterSpec := megafauna.MakeDNASpec(eaters[i])
		foodSpec := megafauna.MakeDNASpec(food[i])
		if eaterSpec.CanFeedOn(foodSpec) {
			t.Errorf("%v shouldn't be able to feed on %v", eaterSpec.Spec, foodSpec.Spec)
		}
	}

}
