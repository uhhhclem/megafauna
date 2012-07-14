package megafauna

// Card represents a card, either a Mutation card or a Genotype card.
type Card struct {
	Key   string
	Event *Event
}

// MutationCard represents a mutation card
type MutationCard struct {
	Card
	MinSize     int
	MaxSize     int
	Mutation    *DNASpec
	InstinctKey string
	Supertitle  string // text appearing above the title
	Title       string // the large text on the card
	Subtitle    string // the small text on the card - note this can be quite long
	Reminder    string // a few mutation cards have reminders, e.g. Caudal Fin
}

// GenotypeCard represents a genotype card.
type GenotypeCard struct {
	MammalData   *GenotypeCardData
	DinosaurData *GenotypeCardData
}

// GenotypeCardData represents the data in one half of a genotype card.
type GenotypeCardData struct {
	SilhouetteIndex int
	MinSize         int
	MaxSize         int
	DNAGenome       *DNASpec
	Family          string
	Title           string
	Subtitle        string
}

// Event represents the event portion of a card.
type Event struct {
	IsDrawTwo               bool
	IsMilankovich           bool
	MilankovichLatitudeKeys []string
	IsCatastrophe           bool
	CatastropheLevel        int
	CatastropheIsWarming    bool
	Description             string
}
