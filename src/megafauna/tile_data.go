package megafauna

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Indices into the CSV file for the biome tile fields.
const (
	biomeTileKeyField = iota
	biomeTileIsMesozoicField
	biomeTileLatitudeKeyField
	biomeTileClimaxNumberField
	biomeTileTypeField
	biomeTileSupertitleField
	biomeTileTitleField
	biomeTileRequirementsField
	biomeTileRooterRequirementsField
	biomeTileNicheField
	biomeTileRedStarField
	biomeTileBlueStarField
	biomeTileWarmingField
	biomeTileCoolingField
)

// Indices into the CSV file for the immigrant tile fields.
const (
	immigrantTileKeyField = iota
	immigrantTileIsMesozoicField
	immigrantTileLatitudeKeyField
	immigrantTileIsLandField
	immigrantTileIsSeaField
	immigrantTileSizeField
	immigrantTileDNAField
	immigrantTileSupertitleField
	immigrantTileTitleField
)

var (
	errInvalidType        = errors.New("Type must be one of " + BiomeTypeKeys + ".")
	errInvalidLatitudeKey = errors.New("LatitudeKey must be one of " + LatitudeKeys + ".")
)

// GetTiles builds the master map of all tiles.
func GetTiles() (map[string]*Tile, error) {
	tiles := make(map[string]*Tile)
	reader := strings.NewReader(biomeTileSourceData)
	err := parsebiomeTiles(reader, tiles)
	if err != nil {
		return nil, err
	}
	reader = strings.NewReader(immigrantTileSourceData)
	err = parseimmigrantTiles(reader, tiles)
	if err != nil {
		return nil, err
	}
	return tiles, nil
}

// parsebiomeTiles takes a Reader containing Tile data in CSV format, parses the data into Tiles, and populates
// the (pre-made) map with the Tiles. 
func parsebiomeTiles(r io.Reader, tiles map[string]*Tile) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		t := new(Tile)
		t.BiomeData = new(BiomeTileData)
		b := t.BiomeData

		t.Key = record[biomeTileKeyField]
		t.IsMesozoic, err = strconv.ParseBool(record[biomeTileIsMesozoicField])
		if err != nil {
			return nil
		}

		t.LatitudeKey = record[biomeTileLatitudeKeyField]

		if len(t.LatitudeKey) != 1 || !strings.Contains(LatitudeKeys, t.LatitudeKey) {
			return errInvalidLatitudeKey
		}

		b.ClimaxNumber, err = strconv.Atoi(record[biomeTileClimaxNumberField])
		if err != nil {
			return err
		}

		t.Supertitle = record[biomeTileSupertitleField]
		t.Title = record[biomeTileTitleField]

		switch record[biomeTileTypeField] {
		case "L":
			t.IsLand = true
			if t.LatitudeKey == "O" {
				b.IsOrogeny = true
			}
		case "S":
			t.IsSea = true
		case "B":
			t.IsLand = true
			t.IsSea = true
		default:
			return errInvalidType
		}

		b.Requirements = MakeDNASpec(record[biomeTileRequirementsField])
		spec := record[biomeTileRooterRequirementsField]
		if spec != "" {
			b.RooterRequirements = MakeDNASpec(record[biomeTileRooterRequirementsField])
		}
		b.Niche, err = MakeNiche(record[biomeTileNicheField])
		if err != nil {
			return err
		}
		b.RedStar, err = strconv.ParseBool(record[biomeTileRedStarField])
		if err != nil {
			return nil
		}
		b.BlueStar, err = strconv.ParseBool(record[biomeTileBlueStarField])
		if err != nil {
			return err
		}
		b.IsWarming, err = strconv.ParseBool(record[biomeTileWarmingField])
		if err != nil {
			return err
		}
		b.IsCooling, err = strconv.ParseBool(record[biomeTileCoolingField])
		if err != nil {
			return err
		}
		tiles[t.Key] = t
	}
	return nil
}

// parseimmigrantTiles takes a Reader containing Tile data in CSV format, parses the data into Tiles, and populates
// the (pre-made) map with the Tiles. 
func parseimmigrantTiles(r io.Reader, tiles map[string]*Tile) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	fmterr := func(index int, err error) error {
		return fmt.Errorf("Line %v: %v", index+1, err.Error())
	}

	for index, record := range records {
		t := new(Tile)
		t.ImmigrantData = new(ImmigrantTileData)
		i := t.ImmigrantData

		t.Key = record[immigrantTileKeyField]
		t.IsMesozoic, err = strconv.ParseBool(record[immigrantTileIsMesozoicField])
		if err != nil {
			return fmterr(index, err)
		}

		t.LatitudeKey = record[immigrantTileLatitudeKeyField]

		if len(t.LatitudeKey) != 1 || !strings.Contains(LatitudeKeys, t.LatitudeKey) {
			return fmterr(index, errInvalidLatitudeKey)
		}

		t.Supertitle = record[immigrantTileSupertitleField]
		t.Title = record[immigrantTileTitleField]

		t.IsLand, err = strconv.ParseBool(record[immigrantTileIsLandField])
		if err != nil {
			return fmterr(index, err)
		}
		t.IsSea, err = strconv.ParseBool(record[immigrantTileIsSeaField])
		if err != nil {
			return fmterr(index, err)
		}

		size := record[immigrantTileSizeField]
		if size == "" {
			i.IsHerbivore = false
		} else {
			i.IsHerbivore = true
			i.Size, err = strconv.Atoi(size)
			if err != nil {
				return fmterr(index, err)
			}
		}

		i.DNA = MakeDNASpec(record[immigrantTileDNAField])

		tiles[t.Key] = t
	}
	return nil
}

// GetInheritanceTiles parses the inheritance tile data and returns a slice containing the 5 inheritance tiles. 
func GetInheritanceTiles() []*InheritanceTile {
	r := strings.NewReader(inheritanceTileSourceData)
	csvReader := csv.NewReader(r)
	records, _ := csvReader.ReadAll()
	result := make([]*InheritanceTile, 0)
	for _, record := range records {
		t := new(InheritanceTile)
		ob := new(InheritanceTileData)
		re := new(InheritanceTileData)
		t.Obverse = ob
		t.Reverse = re

		ob.MinSize, _ = strconv.Atoi(record[0])
		ob.MaxSize, _ = strconv.Atoi(record[1])
		ob.DNA = MakeDNASpec(record[2])

		re.MinSize, _ = strconv.Atoi(record[3])
		re.MaxSize, _ = strconv.Atoi(record[4])
		re.DNA = MakeDNASpec(record[5])

		result = append(result, t)
	}
	return result
}

const biomeTileSourceData = `MA16,TRUE,A,16,L,Deciduous Gymnosperm,Polar Forest,BB,H,N,TRUE,FALSE,FALSE,FALSE
MA20,TRUE,A,20,L,Cordaites,Broadleaf Conifer Forest,BB,,size,FALSE,FALSE,FALSE,FALSE
MA30,TRUE,A,30,L,Gingkophytes,Ginkgo Woodland,B,,P,TRUE,FALSE,FALSE,FALSE
MA53,TRUE,A,53,L,Mycophycophytes,Lichen Tundra,GG,,size,FALSE,TRUE,FALSE,FALSE
MA72,TRUE,A,72,L,Deciduous Larix,Larch Forest,BG,,size,FALSE,FALSE,FALSE,FALSE
MH12,TRUE,H,12,L,Cordaitalies,"Araucarites ""Petrified"" Forest",BB,,A,FALSE,TRUE,FALSE,FALSE
MH14,TRUE,H,14,S,Ganoid,Palaeoniscid Ray-Fins,M,,M,FALSE,FALSE,FALSE,FALSE
MH24,TRUE,H,24,L,Gymnosperm Bennititalean,Cycadoid,B,H,S,FALSE,FALSE,FALSE,FALSE
MH38,TRUE,H,38,S,Crinoid,Sea Lily Beds,AM,,M,FALSE,FALSE,FALSE,FALSE
MH43,TRUE,H,43,S,Cephalopod,Mediterranean Ammonites,AM,,size,FALSE,FALSE,FALSE,FALSE
MH44,TRUE,H,44,S,Phaeophyta,Coastal Kelp Forest,MM,,size,FALSE,FALSE,FALSE,FALSE
MH48,TRUE,H,48,S,Cephalopod,Belemnites,M,,N,FALSE,FALSE,FALSE,FALSE
MH49,TRUE,H,49,S,Pelecypod,Clam and Oyster Beds,AAM,,P,FALSE,FALSE,FALSE,FALSE
MJ15,TRUE,J,15,L,Pteridospermophytes,Seed Ferns,B,,I,TRUE,FALSE,FALSE,FALSE
MJ17,TRUE,J,17,L,Cycadofilicales,Cycad Woods,BG,H,size,TRUE,FALSE,FALSE,FALSE
MJ33,TRUE,J,33,S,Coccolithophore,Calcareous Plankton Bloom,MM,,size,FALSE,FALSE,FALSE,FALSE
MJ34,TRUE,J,34,L,Pteridophytes,Tree Ferns,BB,,A,FALSE,FALSE,FALSE,FALSE
MJ40,TRUE,J,40,L,Cypressales,Dawn Redwood Forest,BB,H,A,FALSE,FALSE,FALSE,FALSE
MJ51,TRUE,J,51,S,Ceratids,Euro-Boreal Ammonites,AM,,N,TRUE,FALSE,FALSE,FALSE
MJ52,TRUE,J,52,L,Pinaceae Pagiophylum,Conifer Taiga,BB,H,I,FALSE,FALSE,FALSE,FALSE
MJ54,TRUE,J,54,L,Coniferales,Boreal Spruce Tree Forest,BB,H,A,FALSE,FALSE,FALSE,FALSE
MO67,TRUE,O,67,L,Variscan orogeny,African Podocarp High Forest,B,,size,TRUE,FALSE,FALSE,FALSE
MO69,TRUE,O,69,L,DeGeer bridge,Iberian Sphagnum Bog,I,,G,TRUE,FALSE,FALSE,FALSE
MO70,TRUE,O,70,L,Siberian traps,Alpine Liverwort Meadow,G,,P,FALSE,FALSE,TRUE,FALSE
MO94,TRUE,O,94,L,CAMP,Disaster Fern Spike,I,,S,FALSE,FALSE,TRUE,FALSE
MT10,TRUE,T,10,L,Lycopodiophytes,Lepidodendrales Coal Swamp,BB,,size,FALSE,FALSE,FALSE,FALSE
MT18,TRUE,T,18,L,Equisetophyta,Calamites Thicket,B,N,I,FALSE,FALSE,FALSE,FALSE
MT26,TRUE,T,26,L,Lycopodiophytes,Lycopod Meadow,I,,N,FALSE,TRUE,FALSE,FALSE
MT27,TRUE,T,27,L,Equisetophyta,Horsetail Swamp,G,N,size,FALSE,FALSE,FALSE,FALSE
MT28,TRUE,T,28,S,Lycophyta,Quillwort Coastal Marsh,M,,G,FALSE,TRUE,FALSE,FALSE
MT32,TRUE,T,32,S,Brachiopod,Lampshell Shoals,AAM,,P,FALSE,FALSE,FALSE,FALSE
MT36,TRUE,T,36,S,Actinistian,Coelacanth Lobe-Fins,MM,,P,FALSE,FALSE,FALSE,FALSE
MT41,TRUE,T,41,L,Cycadales,Cycad Thicket,B,H,S,FALSE,TRUE,FALSE,FALSE
MT42,TRUE,T,42,L,Isoptera,Termite Mounds,II,,A,FALSE,TRUE,FALSE,FALSE
MT45,TRUE,T,45,S,Rudistid Reef,Aspidorhyneid Ray-Fins,M,,M,FALSE,FALSE,FALSE,FALSE
MT46,TRUE,T,46,S,Sponge Reef,Pycnodontid Ray-Fins,M,,P,FALSE,FALSE,FALSE,FALSE
MT50,TRUE,T,50,S,Decapod,Lobsters,MM,,N,FALSE,FALSE,FALSE,FALSE`

const immigrantTileSourceData = `M1,TRUE,A,TRUE,FALSE,1,HIIN,Primitive mammal,Multituberculates
M2,TRUE,A,TRUE,FALSE,2,BGGA,Primitive ceratopsian,Protoceratops
M3,TRUE,A,TRUE,FALSE,5,BBGA,Ornithopod,Igunaodonts
M4,TRUE,A,TRUE,FALSE,,PAASS,Saurischian theropod,Velociraptors
M5,TRUE,J,TRUE,TRUE,3,BHHAA,Primitive thyreophoran,Scelidosaurs
M6,TRUE,J,TRUE,FALSE,5,BBA,Theropod maniraptora,Therizinosaurs
M7,TRUE,T,TRUE,FALSE,1,BII,Sphenosuchian,Crocodiles
M8,TRUE,T,TRUE,FALSE,1,BIS,Primitive ornothschian,Fabrosaurs
M9,TRUE,T,FALSE,TRUE,2,AMM,Diapsid Chorstoderan,Champosaurs
M10,TRUE,T,TRUE,TRUE,2,HAAM,Euryapside reptile,Placodonts
M11,TRUE,T,TRUE,FALSE,6,BBA,Armored sauropod,Titanosaurs
M12,TRUE,T,FALSE,TRUE,6,MMN,Euryapside reptile,Ichthyosaurs
M13,TRUE,T,FALSE,TRUE,,PAMM,Diapsid thracophra,Pliosaurs
M14,TRUE,T,FALSE,TRUE,,AAMM,Diapsid lizard,Mosasaurs
M15,TRUE,T,TRUE,TRUE,,MNN,Lissamphibian capotosaur,Labyrinthodonts
M16,TRUE,T,TRUE,TRUE,,AAM,Mesosuchian sebecid,Terrestrial crocodiles`

const inheritanceTileSourceData = `1,2,H,1,4,B
1,6,G,1,6,GG
1,2,I,1,6,G
1,4,B,4,6,BB
1,6,P,1,6,PP`
