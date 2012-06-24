package megafauna_test

import (
	"megafauna"
	"testing"
)

func TestMakeNiche(t *Testing.T) {
	n, err := megafauna.MakeNiche("Size")
	if err {
		t.Errorf(err.Error())
	}
	if !n.Size {
		t.Errorf("Expected a Size niche and didn't get one.")
	} 
}
