package megafauna

import (
	"fmt"
	"strconv"
	"strings"
)

// Niche is used as the first tiebreaker in herbivore contests.  A Niche will be *one* of these: Size, 
// a player color (Dentition), or a DNA letter.
type Niche struct {
	Size      bool
	Dentition int
	DNA       string
}

type MakeNicheError struct {
	spec string
}

func (e *MakeNicheError) Error() string {
	return fmt.Sprintf("MakeNiche: %v is not a valid niche spec.", e.spec)
}

// MakeNiche makes a new Niche, given an input string that is one of "Size", "2" through "5" (for player dentition),
// or a DNA letter.
func MakeNiche(spec string) (n *Niche, err error) {
	spec = strings.TrimSpace(spec)
	n = new(Niche)
	if strings.ToUpper(spec) == "SIZE" {
		n.Size = true
		return
	}
	dentition, converr := strconv.ParseInt(spec, 0, 0)
	if converr == nil {
		n.Dentition = int(dentition)
		return
	} 
	if len(spec) == 1 && strings.Contains("BGHIANMS", spec) {
		n.DNA = spec
	}
	err = MakeNicheError { spec }
}
