package megafauna

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Indices into the CSV data for the mutation cards
const (
	mutationCardKeyField = iota
	mutationCardMinSizeField
	mutationCardMaxSizeField
	mutationCardMutationField
	mutationCardInstinctField
	mutationCardSupertitleField
	mutationCardTitleField
	mutationCardSubtitleField
	mutationCardEventTypeField
	mutationCardEventDescriptionField
	mutationCardCatastropheLevelField
	mutationCardCatastropheIsWarmingField
	mutationCardMilankovichLatitudesField
)

var (
	ErrInvalidEventType        = errors.New("Invalid event type.")
	ErrInvalidLatitudeKey      = errors.New("Invalid latitude key.")
	ErrInvalidDNASpec          = errors.New("Invalid DNA spec.")
	ErrInvalidInstinctKey      = errors.New("Invalid instinct key.")
	ErrInvalidCatastropheLevel = errors.New("Invalid catastrophe level.")
	ErrInvalidMinSize          = errors.New("Invalid minimum size.")
	ErrInvalidMaxSize          = errors.New("Invalid maximum size.")
)

// MutationCardMap is a map of string keys to MutationCard objects.
type MutationCardMap map[string]*MutationCard

// Parse parses CSV data from the Reader passed in into the receiver MutationCardMap.
func (cards MutationCardMap) Parse(r io.Reader) error {
	csvReader := csv.NewReader(r)
	csvReader.TrailingComma = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		m := new(MutationCard)
		m.Key = record[mutationCardKeyField]
		m.MinSize, err = strconv.Atoi(record[mutationCardMinSizeField])
		if err != nil {
			return ErrInvalidMinSize
		}
		m.MaxSize, err = strconv.Atoi(record[mutationCardMaxSizeField])
		if err != nil {
			return ErrInvalidMaxSize
		}
		instinct := record[mutationCardInstinctField]
		if instinct != "" {
			if len(instinct) != 1 || !strings.Contains(InstinctKeys, instinct) {
				return ErrInvalidInstinctKey
			}
			m.InstinctKey = instinct
		}
		m.Mutation = MakeDNASpec(record[mutationCardMutationField])
		if m.Mutation == nil {
			return ErrInvalidDNASpec
		}

		m.Supertitle = record[mutationCardSupertitleField]
		m.Title = record[mutationCardTitleField]
		m.Subtitle = record[mutationCardSubtitleField]

		var catLevel int
		if record[mutationCardCatastropheLevelField] != "" {
			catLevel, err = strconv.Atoi(record[mutationCardCatastropheLevelField])
			if err != nil {
				return ErrInvalidCatastropheLevel
			}
		}
		catWarming := (record[mutationCardCatastropheIsWarmingField] == "T")
		m.Event, err = MakeEvent(record[mutationCardEventTypeField], record[mutationCardMilankovichLatitudesField], catLevel, catWarming)
		if err != nil {
			return err
		}
		m.Event.Description = record[mutationCardEventDescriptionField]
		cards[m.Key] = m
	}
	return nil
}

// MakeEvent creates an event from the appropriate fields in the card data.  Returns error if the data is invalid
func MakeEvent(eventType string, milankovichLatitude string, catastropheLevel int, catastropheIsWarming bool) (*Event, error) {

	var e Event
	eventKey := eventType[0]

	switch {
	case eventKey == 'T':
		e.IsDrawTwo = true
		return &e, nil
	case eventKey == 'C':
		e.IsCatastrophe = true
		e.CatastropheLevel = catastropheLevel
		e.CatastropheIsWarming = catastropheIsWarming
		return &e, nil
	case eventKey == 'M':
		e.IsMilankovich = true
		e.MilankovichLatitudeKeys = make([]string, len(milankovichLatitude))
		for index, key := range milankovichLatitude {
			if !strings.Contains(LatitudeKeys, string(key)) {
				return nil, ErrInvalidLatitudeKey
			}
			e.MilankovichLatitudeKeys[index] = string(key)
		}
		return &e, nil
	}
	return nil, ErrInvalidEventType
}

// CSV data for the mutation cards.
const MutationCardData = `1,1,6,S,,Breathing while running,Carrier's Constant Diaphragm,,T,,,,
2,1,4,SS,,Unidirectional respiration,Flow-Through Lungs,Shown is the bird system.,C,Asteroid impact global cooling,5,FALSE,
3,1,4,PP,,Homoiotherm,Feathers,"A better insulator than fur, but prone to parasites and matting.",T,,,,
4,1,5,B,M,Pubic bone shift,Biped Stance,,T,,,,
5,1,4,S,,Achilles tendon,Flexure Heel,This tendon provides elastic energy storage in hopping and running,T,,,,
6,1,6,P,S,K strategy,Parental protection,,MP,,,,HA
7,1,3,S,,Tail corset,Digitigrade Hopping,A digitigrade is an animal that walks on its toes,T,,,,
8,1,4,N,N,,Infrared Pit Sensor,The ability to sense thermal radiation helps to detect warm-blooded predators or prey.,T,,,,
9,1,6,P,,Precocious,Placental Reproduction,,ME,,,,T
10,1,6,B,M,,Tripod Stance,Walking plantigrade (with soles of feet flat on the ground) sacrifices speed for stability and weight-bearing,T,,,,
11,4,6,BB,,Four-chambered heart,Long Neck,,T,,,,
12,1,1,H,,Internal husking ridges,Seed-cracking Bill,,C,"Volcanic acid rain, global warming",6,TRUE,
13,1,5,PP,,Thermoregulation,Panting,Many mammals and birds use this form of evaporative cooling,MP,,,,HA
14,1,6,M,,,Sculling Tail,,T,,,,
15,1,6,P,,Seasonal,Hibernation,,T,,,,
16,1,4,PP,,Homoiothermic,Fur,Homoiotherms maintain a constant body temperature despite ambient temperatures.,T,,,,
17,1,2,II,,,Anteater Tongue,,C,Solar Flare global cooling,7,FALSE,
18,2,3,AA,,,Saber Tooth,These animals lunged from ambush to eviscerate the belly of their prey.,T,,,,
19,1,5,A,,Heterodont,Scimitar Incisors Carnassial Molars,,T,,,,
20,1,1,GI,,Dilambodont Cheek Teeth,"W-shaped teeth used by shrews, moles, and bats.",,T,,,,`
