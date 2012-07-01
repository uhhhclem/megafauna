package megafauna

// Card represents a card, either a Mutation card or a Genotype card.
type Card struct {
	Key string
}

type MutationCard struct {
	Card
	DNA      *DNASpec
	Title    string // the large text on the card
	Subtitle string // the small text on the card - note this can be quite long
	Reminder string // a few mutation cards have reminders, e.g. Caudal Fin
	MinSize  int
	MaxSize  int
}

type GenotypeCard struct {
	MammalData   *GenotypeCardData
	DinosaurData *GenotypeCardData
}

type GenotypeCardData struct {
	SilhouetteIndex int
	MinSize         int
	MaxSize         int
	DNAGenome       *DNASpec
	Family          string
	Title           string
	Subtitle        string
}
