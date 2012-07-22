package megafauna

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Nicknames of the animal silhouettes; SilhouetteIndex is derived from this.
var DinosaurSilhouettes = []string{"dino", "fin", "bird", "croc"}
var MammalSilhouettes = []string{"cat", "rhino", "bat", "dolphin"}

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
	ErrInvalidSilhouetteAbbrev = errors.New("Invalid silhouette abbreviation.")
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

// Indices into the CSV data for the mutation cards
const (
	genotypeCardKeyField = iota
	genotypeCardMSilhouetteField
	genotypeCardMFamilyField
	genotypeCardMTitleField
	genotypeCardMSubtitleField
	genotypeCardMMinSizeField
	genotypeCardMMaxSizeField
	genotypeCardMDNASpecField
	genotypeCardDSilhouetteField
	genotypeCardDFamilyField
	genotypeCardDTitleField
	genotypeCardDSubtitleField
	genotypeCardDMinSizeField
	genotypeCardDMaxSizeField
	genotypeCardDDNASpecField
	genotypeCardEventTypeField
	genotypeCardEventDescriptionField
	genotypeCardCatastropheLevelField
	genotypeCardCatastropheIsWarmingField
	genotypeCardMilankovichLatitudesField
)

type GenotypeCardMap map[string]*GenotypeCard

// Parse parses CSV data from the Reader passed in into the receiver GenotypeCardMap.
func (cards GenotypeCardMap) Parse(r io.Reader) error {
	csvReader := csv.NewReader(r)
	csvReader.TrailingComma = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		var err error
		var data *GenotypeCardData
		g := new(GenotypeCard)
		g.Key = record[genotypeCardKeyField]
		data, err = parseGenotypeData(record, genotypeCardMSilhouetteField)
		if err != nil {
			return err
		}
		g.MammalData = data
		data, err = parseGenotypeData(record, genotypeCardDSilhouetteField)
		if err != nil {
			return err
		}
		g.DinosaurData = data

		var catLevel int
		if record[genotypeCardCatastropheLevelField] != "" {
			catLevel, err = strconv.Atoi(record[genotypeCardCatastropheLevelField])
			if err != nil {
				return ErrInvalidCatastropheLevel
			}
		}
		catWarming := (record[genotypeCardCatastropheIsWarmingField] == "T")
		g.Event, err = MakeEvent(record[genotypeCardEventTypeField], record[genotypeCardMilankovichLatitudesField], catLevel, catWarming)
		if err != nil {
			return err
		}
		g.Event.Description = record[genotypeCardEventDescriptionField]
		cards[g.Key] = g
	}
	return nil
}

// ConvertSilhouette returns the silhouette index for an abbreviation.
func ConvertSilhouette(abbrev string, isDinosaur bool) (int, error) {
	silhouettes := MammalSilhouettes
	if isDinosaur {
		silhouettes = DinosaurSilhouettes
	}
	for i := 0; i < 4; i++ {
		if abbrev == silhouettes[i] {
			return i, nil
		}
	}
	return -1, ErrInvalidSilhouetteAbbrev
}

// parseGenotypeData parses the genotype data starting in field startField of the record
func parseGenotypeData(record []string, startField int) (*GenotypeCardData, error) {
	var err error
	var silhouette int
	g := new(GenotypeCardData)

	isDinosaur := !(startField == genotypeCardMSilhouetteField)
	silhouette, err = ConvertSilhouette(record[startField], isDinosaur)
	startField++
	if err != nil {
		return nil, err
	}
	g.SilhouetteIndex = silhouette

	g.Family = record[startField]
	startField++
	g.Title = record[startField]
	startField++
	g.Subtitle = record[startField]
	startField++
	g.MinSize, err = strconv.Atoi(record[startField])
	if err != nil {
		return nil, ErrInvalidMinSize
	}
	startField++
	g.MaxSize, err = strconv.Atoi(record[startField])
	if err != nil {
		return nil, ErrInvalidMaxSize
	}
	startField++
	g.DNASpec = MakeDNASpec(record[startField])
	if g.DNASpec == nil {
		return nil, ErrInvalidDNASpec
	}
	startField++
	return g, nil
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
const MutationCardSourceData = `M1,1,6,S,,Breathing while running,Carrier's Constant Diaphragm,,T,,,,
M2,1,4,SS,,Unidirectional respiration,Flow-Through Lungs,Shown is the bird system.,C,Asteroid impact global cooling,5,FALSE,
M3,1,4,PP,,Homoiotherm,Feathers,"A better insulator than fur, but prone to parasites and matting.",T,,,,
M4,1,5,B,M,Pubic bone shift,Biped Stance,,T,,,,
M5,1,4,S,,Achilles tendon,Flexure Heel,This tendon provides elastic energy storage in hopping and running,T,,,,
M6,1,6,P,S,K strategy,Parental protection,,MP,,,,HA
M7,1,3,S,,Tail corset,Digitigrade Hopping,A digitigrade is an animal that walks on its toes,T,,,,
M8,1,4,N,N,,Infrared Pit Sensor,The ability to sense thermal radiation helps to detect warm-blooded predators or prey.,T,,,,
M9,1,6,P,,Precocious,Placental Reproduction,,ME,,,,T
M10,1,6,B,M,,Tripod Stance,Walking plantigrade (with soles of feet flat on the ground) sacrifices speed for stability and weight-bearing,T,,,,
M11,4,6,BB,,Four-chambered heart,Long Neck,,T,,,,
M12,1,1,H,,Internal husking ridges,Seed-cracking Bill,,C,"Volcanic acid rain, global warming",6,TRUE,
M13,1,5,PP,,Thermoregulation,Panting,Many mammals and birds use this form of evaporative cooling,MP,,,,HA
M14,1,6,M,,,Sculling Tail,,T,,,,
M15,1,6,P,,Seasonal,Hibernation,,T,,,,
M16,1,4,PP,,Homoiothermic,Fur,Homoiotherms maintain a constant body temperature despite ambient temperatures.,T,,,,
M17,1,2,II,,,Anteater Tongue,,C,Solar Flare global cooling,7,FALSE,
M18,2,3,AA,,,Saber Tooth,These animals lunged from ambush to eviscerate the belly of their prey.,T,,,,
M19,1,5,A,,Heterodont,Scimitar Incisors Carnassial Molars,,T,,,,
M20,1,1,GI,,Dilambodont Cheek Teeth,"W-shaped teeth used by shrews, moles, and bats.",,T,,,,`

const GenotypeCardSourceData = `G1,cat,Carnivora,Feloids,cats,1,3,PN,dino,Saurischian theropod,Ostrich dinosaurs,"oviraptors, ornithomimids",1,5,PS,MP,,,,HA
G2,cat,Pholidota,Pangolins,,1,2,IA,fin,Crurotarsi,Aetosaurs,,1,3,AN,T,,,,
G3,rhino,Artodactyl ungulate,Swine,"pigs, hippos",1,4,GP,dino,Ornithischian ornithopod,Duckbills,"lambeosaurines, iguanodonts, hadrosaurs",2,5,GG,T,,,,
G4,dolphin,Sirenia,Sea Cows,"dugongs, manatees",1,5,GM,croc,Diapsid reptile,Nothosaurs,,1,3,NM,T,,,,
G5,rhino,Marsupial,Diprotodonts,"kangaroos, wombats",1,5,BS,croc,Crurotarsi,Crocodiles,,1,4,AM,T,,,,
G6,bat,Anthropoidea,Primates,"monkeys, apes",1,2,PPPP,bird,Neornithes strigiformes,Owls,,1,1,NSSS,T,,,,
G7,rhino,Edentate Xenarthra,Glyptodonts,,1,5,GAA,fin,Ornithischian thyracophora,"Ankylosaurs, Nodosaurs",,2,5,GAA,T,,,,
G8,bat,Chiroptera,Bats,,1,1,NSSS,bird,Pterosaur,Flying Reptiles,"rhamphorhynchs, pterodactyls",1,1,MSSS,MP,,,,HAA`
