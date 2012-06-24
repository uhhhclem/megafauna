package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestMakeNiche(t *testing.T) {
	n, err := megafauna.MakeNiche("Size")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !n.Size {
		t.Errorf("Expected a Size niche and didn't get one.")
	} 
	
	n, err = megafauna.MakeNiche("3")
	if err != nil {
		t.Errorf(err.Error())
	}
	if n.Size {
		t.Errorf("Size should be false.")
	}
	if n.Dentition != 3 {
		t.Errorf("Expected a Dentition of 3 and got %v.", n.Dentition)
	} 

	n, err = megafauna.MakeNiche("I")
	if err != nil {
		t.Errorf(err.Error())
	}
	if n.Size {
		t.Errorf("Size should be false.")
	}
	if n.Dentition != 0 {
		t.Errorf("Dentition should be 0.")
	}
	if n.DNA != "I" {
		t.Errorf("Expected DNA of I and got %v.", n.Dentition)
	} 
	
	n, err = megafauna.MakeNiche("Q")
	if err == nil {
		t.Errorf("Should have gotten an error and didn't.")
	}
}
