package megafauna


// DNA represents a single DNA value, e.g. BB or AAA
type DNA struct {
	Letter string
	Value  int
}

// IsRoadrunner tells you whether or not this is roadrunner DNA.
func (d *DNA) IsRoadrunner() bool {
	return d.Letter == "A" || d.Letter == "M" || d.Letter == "N" || d.Letter == "S"
}

// IsDietary tells you whether or not this is dietary DNA.
func (d *DNA) IsDietary() bool {
	return !d.IsRoadrunner()
}

// DNASpec represents a DNA specification, i.e. a genome, or a biome's requirements.
type DNASpec struct {
	Spec      string         // e.g. BBG
	Breakdown map[string]*DNA // Spec broken down for ease of comparison
}

// MakeDNASpec makes a DNASpec given an input spec, e.g. "BGGGA"
func MakeDNASpec(spec string) *DNASpec {
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
