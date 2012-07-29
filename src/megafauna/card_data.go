package megafauna

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Nicknames of the animal silhouettes; SilhouetteIndex is derived from this.
var DinosaurSilhouettes = []string{"dino", "fin", "bird", "croc"}
var MammalSilhouettes = []string{"cat", "rhino", "bat", "dolphin"}

var (
	ErrInvalidLatitudeKey      = errors.New("Invalid latitude key.")
	ErrInvalidDNASpec          = errors.New("Invalid DNA spec.")
	ErrInvalidInstinctKey      = errors.New("Invalid instinct key.")
	ErrInvalidCatastropheLevel = errors.New("Invalid catastrophe level.")
	ErrInvalidMinSize          = errors.New("Invalid minimum size.")
	ErrInvalidMaxSize          = errors.New("Invalid maximum size.")
	ErrInvalidSilhouetteAbbrev = errors.New("Invalid silhouette abbreviation.")
)

// GetCards builds the master map of all Cards.
func GetCards() (map[string]*Card, error) {
	cards := make(map[string]*Card)
	reader := strings.NewReader(mutationCardSourceData)
	err := parseMutationCards(reader, cards)
	if err != nil {
		return nil, err
	}
	reader = strings.NewReader(genotypeCardSourceData)
	err = parseGenotypeCards(reader, cards)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

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

// parseMutationCards parses CSV data from the Reader and puts the resulting MutationCard objects in the
// cards map.
func parseMutationCards(r io.Reader, cards map[string]*Card) error {
	csvReader := csv.NewReader(r)
	csvReader.TrailingComma = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	fmterr := func(index int, err error) error {
		return fmt.Errorf("Line %v: %v", index+1, err.Error())
	}

	for index, record := range records {
		c := new(Card)
		m := new(MutationCard)
		c.Mutation = m

		c.Key = record[mutationCardKeyField]

		m.MinSize, err = strconv.Atoi(record[mutationCardMinSizeField])
		if err != nil {
			return fmterr(index, ErrInvalidMinSize)
		}
		m.MaxSize, err = strconv.Atoi(record[mutationCardMaxSizeField])
		if err != nil {
			return fmterr(index, ErrInvalidMaxSize)
		}
		instinct := record[mutationCardInstinctField]
		if instinct != "" {
			if len(instinct) != 1 || !strings.Contains(InstinctKeys, instinct) {
				return fmterr(index, ErrInvalidInstinctKey)
			}
			m.InstinctKey = instinct
		}
		m.Mutation = MakeDNASpec(record[mutationCardMutationField])
		if m.Mutation == nil {
			return fmterr(index, ErrInvalidDNASpec)
		}

		m.Supertitle = record[mutationCardSupertitleField]
		m.Title = record[mutationCardTitleField]
		m.Subtitle = record[mutationCardSubtitleField]

		var catLevel int
		if record[mutationCardCatastropheLevelField] != "" {
			catLevel, err = strconv.Atoi(record[mutationCardCatastropheLevelField])
			if err != nil {
				return fmterr(index, ErrInvalidCatastropheLevel)
			}
		}
		catWarming := (record[mutationCardCatastropheIsWarmingField] == "T")
		c.Event, err = makeEvent(record[mutationCardEventTypeField], record[mutationCardMilankovichLatitudesField], catLevel, catWarming)
		if err != nil {
			return fmterr(index, err)
		}
		c.Event.Description = record[mutationCardEventDescriptionField]

		cards[c.Key] = c
	}
	return nil
}

// Indices into the CSV data for the genotype cards
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

// parseGenotypeCards parses CSV data from the Reader passed in into the receiver Card map
func parseGenotypeCards(r io.Reader, cards map[string]*Card) error {
	csvReader := csv.NewReader(r)
	csvReader.TrailingComma = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		var err error
		var data *GenotypeCardData
		c := new(Card)
		c.Key = record[genotypeCardKeyField]

		g := new(GenotypeCard)
		c.Genotype = g

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
		c.Event, err = makeEvent(record[genotypeCardEventTypeField], record[genotypeCardMilankovichLatitudesField], catLevel, catWarming)
		if err != nil {
			return err
		}
		c.Event.Description = record[genotypeCardEventDescriptionField]
		cards[c.Key] = c
	}
	return nil
}

// convertSilhouette returns the silhouette index for an abbreviation.
func convertSilhouette(abbrev string, isDinosaur bool) (int, error) {
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
	silhouette, err = convertSilhouette(record[startField], isDinosaur)
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

// makeEvent creates an event from the appropriate fields in the card data.  Returns error if the data is invalid.
func makeEvent(eventType string, milankovichLatitude string, catastropheLevel int, catastropheIsWarming bool) (*Event, error) {

	var e Event
	eventKey := eventType[0]
	err := fmt.Errorf("\"%v\" is an invalid event type.", eventType)

	switch {
	default:
		return nil, err
	case eventKey == 'T':
		e.IsDrawTwo = true
	case eventKey == 'G': // can be "GW" or "GC" for global warming or cooling
		if len(eventType) < 2 {
			return nil, err
		}
		e.IsWarming = eventType[1] == 'W'
		e.IsCooling = eventType[1] == 'C'
		if !e.IsWarming && !e.IsCooling {
			return nil, err
		}
	case eventKey == 'C':
		e.IsCatastrophe = true
		e.CatastropheLevel = catastropheLevel
		e.IsWarming = catastropheIsWarming
		e.IsCooling = !catastropheIsWarming
	case eventKey == 'M':
		e.IsMilankovich = true
		e.MilankovichLatitudeKeys = make([]string, len(milankovichLatitude))
		for index, key := range milankovichLatitude {
			if !strings.Contains(LatitudeKeys, string(key)) {
				return nil, ErrInvalidLatitudeKey
			}
			e.MilankovichLatitudeKeys[index] = string(key)
		}
	}
	return &e, nil
}

// CSV data for the mutation cards.
const mutationCardSourceData = `M1,1,6,S,,Breathing while running,Carrier's Constant Diaphragm,,T,,,,
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
M20,1,1,GI,,Dilambodont Cheek Teeth,"W-shaped teeth used by shrews, moles, and bats.",,T,,,,
M21,1,6,P,L,Encephalization,Forebrain,,T,,,,
M22,1,6,P,L,Intraspecies,Communication,,T,,,,
M23,1,4,SS,,Unguligrade,Hooves,Keratin-reinforced toe-tips are used by swift animals from deer to land crocodiles.,GC,Erosion global cooling,,,
M24,1,1,AA,,Expendable,Spines,,C,Clathrate Gun global warming,4,TRUE,
M25,1,6,,S,Scent gland marking,Courtship & Territorialism,,T,,,,
M26,1,6,PS,,Endocrine hormones,Pituitary and Thyroid Glands,,T,,,,
M27,1,6,P,,Fatty,Hump,,T,,,,
M28,2,6,BG,,,Hindgut Digestion,"To digest leaves, and elongated intestine and colon plus an enlarged cecum to host bacteria are needed.",T,,,,
M29,1,6,BG,,,"Cheeks, Palate, and Tongue",Cheeks retain food during chewing.  Palate allows simultaneous breathing & chewing.,T,,,,
M30,1,3,IN,L,Land or sea sonar,Echolocation,,T,,,,
M31,1,1,HG,,Crop,Foregut Digestion,,T,,,,
M32,1,6,MM,,Natatorial lunate,Caudal Fin,This shape is optimized for cruising.,T,,,,
M33,1,6,AS,,Fight or Flight hormones,Adrenal Glands,,T,,,,
M34,1,2,BI,M,,Opposable Thumb,Shown is the hand of a panda.,T,,,,
M35,3,6,PM,,Homoiotherm,Blubber,,T,,,,
M36,1,6,G,,,Lophodont Cheek Teeth,The addition of hard enamel ridges to teeth improves their grinding action.,T,,,,
M37,1,6,,L,Cooperative,Flushing or Mobbing,"One hunter flushes, the other pounces.  Prey may harass predators in mobs.",ME,,,,T
M38,1,2,P,N,Seasonal,Food Storage,Squirrels and woodpeckers act as seed-dispersing agents for oak trees.,T,,,,
M39,1,6,M,,,Salt-excreting Tubenose,The tubes on some seabird beaks allows them to drink seawater.,GC,Erosion global cooling,,,
M40,1,4,N,,"Hammer, anvil, and stirrup bones",Directional Ears,Mammal evolution reshaped 3 jawbones to form a chain of auditory bones.,GC,Erosion global cooling,,,
M41,1,2,N,,Chromatophore,Camouflage,,T,,,,
M42,2,5,A,,Disemboweling,Switchblade Claw,,T,,,,
M43,1,6,,S,Mating Call,Vocal Amplifier,,C,X-Ray Burdeter global cooling,6,FALSE,
M44,1,2,IN,N,Rodent den excavation,Digging Claws,,MT,,,,J
M45,1,6,P,S,"Herds, Flocks, or Packs",Seasonal Migration,"As the world got more seasonal, animals responded by going south for the winter.",T,,,,
M46,2,6,,S,"Antlers, boneheads, & horns",Head Butting,,T,,,,
M47,1,4,M,,Cardiform,Needle Teeth,These teeth are better adapted for seizing slippery fish than tackling a land animal.,T,,,,
M48,1,6,N,,Light sensitive rods,Night Binocular Vision,,ME,,,,T
M49,1,6,M,,Natatorial,Web Feet,,GC,,,,
M50,1,6,BG,,Gastric Mill,Gizzard Stones,"Some animals, from brontosaurs to pigeons, swallow stones to grind food.",T,,,,
M51,1,2,IM,N,Electroreceptive,Rooting Snout,,T,,,,
M52,1,6,,S,,Lek,Leks are social arenas where males assemble to perform courtship displays.,T,,,,
M53,2,6,B,M,Prehensile,Trunk,,GC,,,,
M54,2,5,P,,Thermoregulation,Sail Back,,T,,,,
M55,1,6,N,N,Olfactory receptor,Bloodhound Nose,,T,,,,
M56,1,1,H,,,Grawing Incisors,Rodents have a set of continuously growing incisors that must be kept short by gnawing.,T,,,,
M57,1,3,AA,,,Plate Armor,A turtle's spine and ribs are fused to interlocking bony plates under the skin.,GC,,,,
M58,3,6,AA,,,Tail Club,,T,,,,
M59,1,6,S,L,Sentry,Warning Cry,,T,,,,
M60,1,6,B,M,,Prehensile Tongue,,T,,,,
M61,1,5,M,M,Water vascular,Tentacles,,T,,,,
M62,1,5,P,M,,Nest Building,,MP,,,,HA
M63,2,6,AA,,,Nose Horn,,T,,,,
M64,2,3,BGG,,,Cud-chewing,Ruminents are mammals that return food from the reticulorumen foregut to be chewed for additional processing.,T,,,,
M65,1,6,BA,,Horny,Beak,,T,,,,
M66,1,2,MN,,Serpentine,Vermiform,,T,,,,`

const genotypeCardSourceData = `G1,cat,Carnivora,Feloids,cats,1,3,PN,dino,Saurischian theropod,Ostrich dinosaurs,"oviraptors, ornithomimids",1,5,PS,MP,,,,HA
G2,cat,Pholidota,Pangolins,,1,2,IA,fin,Crurotarsi,Aetosaurs,,1,3,AN,T,,,,
G3,rhino,Artodactyl ungulate,Swine,"pigs, hippos",1,4,GP,dino,Ornithischian ornithopod,Duckbills,"lambeosaurines, iguanodonts, hadrosaurs",2,5,GG,T,,,,
G4,dolphin,Sirenia,Sea Cows,"dugongs, manatees",1,5,GM,croc,Diapsid reptile,Nothosaurs,,1,3,NM,T,,,,
G5,rhino,Marsupial,Diprotodonts,"kangaroos, wombats",1,5,BS,croc,Crurotarsi,Crocodiles,,1,4,AM,T,,,,
G6,bat,Anthropoidea,Primates,"monkeys, apes",1,2,PPPP,bird,Neornithes strigiformes,Owls,,1,1,NSSS,T,,,,
G7,rhino,Edentate Xenarthra,Glyptodonts,,1,5,GAA,fin,Ornithischian thyracophora,"Ankylosaurs, Nodosaurs",,2,5,GAA,T,,,,
G8,bat,Chiroptera,Bats,,1,1,NSSS,bird,Pterosaur,Flying Reptiles,"rhamphorhynchs, pterodactyls",1,1,MSSS,MP,,,,HAA`
