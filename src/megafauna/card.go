package megafauna

type Card struct {
	Key      string
	Event    *Event
	Mutation *MutationCard
	Genotype *GenotypeCard
}

// MutationCard contains mutation-card-specific fields.
type MutationCard struct {
	MinSize     int
	MaxSize     int
	Mutation    *DNASpec
	InstinctKey string
	Supertitle  string // text appearing above the title
	Title       string // the large text on the card
	Subtitle    string // the small text on the card - note this can be quite long
	Reminder    string // a few mutation cards have reminders, e.g. Caudal Fin
}

// GenotypeCard contains genotype-card-specific fields.
type GenotypeCard struct {
	MammalData   *GenotypeCardData
	DinosaurData *GenotypeCardData
	Event        *Event
}

// GenotypeCardData represents the data in one half of a genotype Card.
type GenotypeCardData struct {
	SilhouetteIndex int
	MinSize         int
	MaxSize         int
	DNASpec         *DNASpec
	Family          string
	Title           string
	Subtitle        string
}

// Event represents the event portion of a Card.
type Event struct {
	IsDrawTwo               bool
	IsMilankovich           bool
	IsWarming               bool
	IsCooling               bool
	MilankovichLatitudeKeys []string
	IsCatastrophe           bool
	CatastropheLevel        int
	Description             string
}
