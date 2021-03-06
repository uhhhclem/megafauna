package megafauna

import "strings"

// DNA represents a single DNA value, e.g. BB or AAA
type DNA struct {
	Letter string
	Value  int
}

// IsRoadrunner tells you whether or not this is roadrunner DNA.
func (d *DNA) IsRoadrunner() bool {
	return len(d.Letter) == 1 && strings.Contains(RoadrunnerDNAKeys, d.Letter)
}

// IsDietary tells you whether or not this is dietary DNA.
func (d *DNA) IsDietary() bool {
	return len(d.Letter) == 1 && strings.Contains(DietaryDNAKeys, d.Letter)
}

// DNASpec represents a DNA specification, i.e. a genome, or a biome's requirements.
type DNASpec struct {
	Spec      string          // e.g. BBG
	Breakdown map[string]*DNA // Spec broken down for ease of comparison
}

// MakeDNASpec makes a DNASpec given an input spec, e.g. "BGGGA"  Returns nil if any
// of the letters in spec are not valid DNA.
func MakeDNASpec(spec string) *DNASpec {
	for _, char := range spec {
		if !strings.Contains(ValidDNAKeys, string(char)) {
			return nil
		}
	}
	d := new(DNASpec)
	d.Spec = spec
	d.Breakdown = make(map[string]*DNA)
	for _, char := range spec {
		letter := string(char)
		dna := d.Breakdown[letter]
		if dna == nil {
			dna = new(DNA)
			dna.Letter = letter
			dna.Value = 1
			d.Breakdown[letter] = dna
		} else {
			dna.Value++
		}
	}
	return d
}

// CanPreyOn determines whether or not an animal with this  DNASpec can prey on another.  
// It returns false if the prey has any roadrunner DNA that this one doesn't, or whose Value 
// exceeds this one's.
func (d *DNASpec) CanPreyOn(other *DNASpec) bool {
	for letter, otherDna := range other.Breakdown {
		if otherDna.IsRoadrunner() {
			thisDna, ok := d.Breakdown[letter]
			if !ok {
				return false
			}
			if thisDna.Value < otherDna.Value {
				return false
			}
		}
	}
	return true
}

// GetDNAValue gets the number of DNA of a given kind that are in the DNASpec
func (d *DNASpec) GetDNAValue(kind string) int {
	if d.Breakdown[kind] == nil {
		return 0
	}
	return d.Breakdown[kind].Value
}

// CanFeedOn determines whether or not an animal with this DNASpec can feed on something - i.e.
// if an herbivore can feed in a biome.  It returns false if the other DNASpec has any dietary DNA 
// that this one doesn't, or whose Value exceeds this one's.
func (d *DNASpec) CanFeedOn(other *DNASpec) bool {
	for letter, otherDna := range other.Breakdown {
		if otherDna.IsDietary() {
			thisDna, ok := d.Breakdown[letter]
			if !ok {
				return false
			}
			if thisDna.Value < otherDna.Value {
				return false
			}
		}
	}
	return true
}
