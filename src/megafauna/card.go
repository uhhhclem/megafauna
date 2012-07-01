package megafauna

// Card represents a card, either a DNA card or a Genotype card.
type Card struct {
	Key				   string
}

type DNACard struct {
	Card
	DNA *DNASpec
	Title string
	MinSize int
	MaxSize int
}