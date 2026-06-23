package plantid

// Plant represents a single plant or fungus entry in the field guide.
type Plant struct {
	ID             string   `json:"id"`
	CommonName     string   `json:"common_name"`
	ScientificName string   `json:"scientific_name"`
	Edibility      string   `json:"edibility"` // edible, edible-caution, poisonous, deadly
	Category       string   `json:"category"`  // plant, fungus, berry, root, leaf
	Habitat        string   `json:"habitat"`
	Season         string   `json:"season"`
	Description    string   `json:"description"`
	Identification []string `json:"identification"`
	Preparation    []string `json:"preparation,omitempty"`
	Warnings       []string `json:"warnings,omitempty"`
	LookAlikes     []string `json:"look_alikes,omitempty"`
}

// PlantGroup groups plants by edibility or type.
type PlantGroup struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Description string  `json:"description"`
	Plants      []Plant `json:"plants"`
}

// PlantDB holds the embedded offline plant identification database.
type PlantDB struct {
	Groups []PlantGroup
}

// NewPlantDB creates the plant field guide database.
func NewPlantDB() *PlantDB {
	return &PlantDB{
		Groups: []PlantGroup{
			{
				ID: "edible-plants", Name: "Edible Plants", Icon: "leaf",
				Description: "Common edible wild plants found across temperate regions",
				Plants: []Plant{
					{
						ID: "dandelion", CommonName: "Dandelion", ScientificName: "Taraxacum officinale",
						Edibility: "edible", Category: "plant", Habitat: "Fields, lawns, roadsides, disturbed ground",
						Season: "Year-round (best spring)", Description: "One of the most common and recognizable edible wild plants worldwide.",
						Identification: []string{
							"Rosette of deeply toothed leaves growing from a central point at ground level.",
							"Single bright yellow flower head on a hollow, milky stem.",
							"Leaves are hairless with teeth pointing back toward the base.",
							"Thick, brown taproot that snaps to reveal white interior.",
							"Seeds form the familiar white 'puffball' sphere.",
						},
						Preparation: []string{
							"Young leaves: eat raw in salads or cook like spinach.",
							"Flowers: batter and fry, or make wine/tea.",
							"Roots: roast and grind as a coffee substitute, or boil as a vegetable.",
						},
					},
					{
						ID: "cattail", CommonName: "Cattail", ScientificName: "Typha latifolia",
						Edibility: "edible", Category: "plant", Habitat: "Marshes, pond edges, ditches, wetlands",
						Season: "Year-round (different parts)", Description: "Called the 'supermarket of the swamp' — almost every part is edible.",
						Identification: []string{
							"Tall (4-8 feet), grass-like plant growing in standing water.",
							"Distinctive brown, cigar-shaped seed head on a stiff stalk.",
							"Long, flat, sword-shaped leaves with parallel veins.",
							"Roots are thick, white rhizomes underwater.",
						},
						Preparation: []string{
							"Young shoots (spring): peel and eat raw or cook like asparagus.",
							"Pollen (early summer): shake into bag, use as flour supplement.",
							"Roots: peel, dry, and pound into flour. Or roast like potatoes.",
							"Green flower heads: boil and eat like corn on the cob.",
						},
						Warnings: []string{
							"Do NOT confuse with Iris (poisonous) — Iris leaves have a central ridge, cattail leaves do not.",
						},
					},
					{
						ID: "clover", CommonName: "Red/White Clover", ScientificName: "Trifolium spp.",
						Edibility: "edible", Category: "plant", Habitat: "Lawns, fields, meadows, roadsides",
						Season: "Spring through fall", Description: "Extremely common, easy to identify, and nutritious.",
						Identification: []string{
							"Three-part leaves (trifoliate) with light chevron pattern.",
							"Round flower heads: white (T. repens) or pink/red (T. pratense).",
							"Low-growing, creeping stems.",
						},
						Preparation: []string{
							"Flowers: eat raw, brew as tea, or dry for flour.",
							"Young leaves: eat raw in small amounts or cook to improve digestibility.",
							"Dried flower heads make excellent tea rich in vitamins.",
						},
					},
				},
			},
			{
				ID: "edible-fungi", Name: "Edible Fungi", Icon: "circle-dot",
				Description: "Safe-to-eat mushrooms with reliable identification features",
				Plants: []Plant{
					{
						ID: "chicken-of-woods", CommonName: "Chicken of the Woods", ScientificName: "Laetiporus sulphureus",
						Edibility: "edible-caution", Category: "fungus", Habitat: "Dead or living hardwood trees (oak, cherry)",
						Season: "Late spring through fall", Description: "Large, bright shelf fungus with a taste and texture like chicken when cooked.",
						Identification: []string{
							"Large, fan-shaped shelves growing in overlapping clusters on trees.",
							"Bright orange-yellow upper surface, sulfur-yellow underneath.",
							"Pore surface underneath (no gills).",
							"Flesh is thick, soft, and white when young.",
							"Can grow to enormous size — 50+ pounds on a single tree.",
						},
						Preparation: []string{
							"Only eat young, tender specimens (soft and moist).",
							"Slice into strips and sauté in butter like chicken cutlets.",
							"Always cook thoroughly — never eat raw.",
						},
						Warnings: []string{
							"Avoid specimens growing on conifers, eucalyptus, or locust trees — these can cause GI upset.",
							"Some people have allergic reactions. Try a small amount first.",
						},
					},
					{
						ID: "morel", CommonName: "Morel", ScientificName: "Morchella spp.",
						Edibility: "edible-caution", Category: "fungus", Habitat: "Forests, burned areas, near dying elms and ash trees",
						Season: "Spring (March-May)", Description: "Highly prized edible mushroom. One of the easiest mushrooms to identify correctly.",
						Identification: []string{
							"Honeycomb-like cap with pits and ridges.",
							"Cap is attached directly to the stem at the base.",
							"Interior is COMPLETELY HOLLOW — cut in half to verify.",
							"Color ranges from blonde/yellow to gray/black depending on species.",
						},
						Preparation: []string{
							"Always cook thoroughly. Never eat raw morels.",
							"Soak in salt water for 30 minutes to drive out insects.",
							"Sauté in butter — considered a delicacy.",
						},
						Warnings: []string{
							"FALSE MORELS (Gyromitra) are poisonous! They have brain-like wrinkled caps (not honeycomb pits) and are NOT hollow inside.",
							"Never eat raw morels — they contain hydrazine compounds destroyed by cooking.",
						},
						LookAlikes: []string{
							"False Morel (Gyromitra esculenta) — wrinkled/brain-like cap instead of honeycomb. NOT hollow. POISONOUS.",
						},
					},
				},
			},
			{
				ID: "poisonous", Name: "Poisonous — AVOID", Icon: "skull",
				Description: "Dangerous plants and fungi to recognize and avoid at all costs",
				Plants: []Plant{
					{
						ID: "death-cap", CommonName: "Death Cap", ScientificName: "Amanita phalloides",
						Edibility: "deadly", Category: "fungus", Habitat: "Under oak and other hardwood trees, parks, gardens",
						Season: "Late summer through fall", Description: "Responsible for the majority of fatal mushroom poisonings worldwide. A single cap can kill an adult.",
						Identification: []string{
							"Cap: 3-6 inches, olive-green to yellowish-green, smooth and slightly sticky.",
							"Gills: white, free from the stem.",
							"Stem: white with a large, skirt-like ring (annulus).",
							"Base: sits in a white, cup-like volva (often buried underground).",
							"Spore print: white.",
							"SMELLS PLEASANT — this is what makes it deadly dangerous.",
						},
						Warnings: []string{
							"FATAL. Symptoms are delayed 6-12 hours, creating false hope. By the time symptoms appear, liver damage is already severe.",
							"There is no reliable home treatment. Seek emergency medical care immediately.",
							"Can be confused with edible Paddy Straw mushrooms and Caesar's mushroom.",
						},
					},
					{
						ID: "water-hemlock", CommonName: "Water Hemlock", ScientificName: "Cicuta spp.",
						Edibility: "deadly", Category: "plant", Habitat: "Wet meadows, stream banks, ditches, marshes",
						Season: "Spring through fall", Description: "Widely considered the most violently toxic plant in North America. Causes seizures within minutes.",
						Identification: []string{
							"2-6 feet tall with compound leaves and small white flowers in umbrella-shaped clusters.",
							"Hollow, chambered stem with purple streaking.",
							"Roots: thick, chambered tuberous roots that smell like raw parsnip.",
							"Leaves are doubly compound with toothed leaflets.",
						},
						Warnings: []string{
							"DEADLY. A single bite of the root can kill an adult.",
							"Causes violent seizures, respiratory failure, and death within hours.",
							"Can be confused with edible plants like wild carrot (Queen Anne's Lace) or parsnip.",
						},
						LookAlikes: []string{
							"Wild Carrot/Queen Anne's Lace — has a hairy stem and a single dark flower in center. Carrot smell when crushed.",
							"Wild Parsnip — similar umbrella flowers but causes severe skin burns from sap.",
						},
					},
					{
						ID: "destroying-angel", CommonName: "Destroying Angel", ScientificName: "Amanita virosa",
						Edibility: "deadly", Category: "fungus", Habitat: "Under hardwood and conifer trees in forests",
						Season: "Summer through fall", Description: "Pure white Amanita that is just as deadly as the Death Cap.",
						Identification: []string{
							"Entirely white — cap, gills, stem, and flesh.",
							"Cap: 2-5 inches, smooth, slightly sticky when wet.",
							"Skirt-like ring on the upper stem.",
							"Bulbous base enclosed in a white sac-like volva.",
							"Spore print: white.",
						},
						Warnings: []string{
							"FATAL. Same toxins as Death Cap (amatoxins).",
							"Can be confused with edible white mushrooms like Meadow Mushrooms or Puffballs.",
							"RULE: Never eat any all-white mushroom with a ring and a volva.",
						},
					},
				},
			},
			{
				ID: "berries", Name: "Wild Berries", Icon: "grape",
				Description: "Edible and poisonous berries — know before you eat",
				Plants: []Plant{
					{
						ID: "blackberry", CommonName: "Blackberry", ScientificName: "Rubus fruticosus",
						Edibility: "edible", Category: "berry", Habitat: "Hedgerows, forest edges, disturbed areas, roadsides",
						Season: "Late summer (July-September)", Description: "One of the easiest and safest wild berries to identify and eat.",
						Identification: []string{
							"Thorny, arching canes forming dense thickets.",
							"Compound leaves with 3-5 toothed leaflets.",
							"White to pink 5-petaled flowers.",
							"Berries are clusters of small drupelets, turning from green → red → black.",
							"Berries do NOT separate from a core when picked (unlike raspberries).",
						},
						Preparation: []string{
							"Eat fresh off the bush.",
							"Can be cooked into jam, pies, or dried for later use.",
						},
					},
					{
						ID: "nightshade-berry", CommonName: "Deadly Nightshade", ScientificName: "Atropa belladonna",
						Edibility: "deadly", Category: "berry", Habitat: "Woodlands, waste ground, limestone soils",
						Season: "Late summer", Description: "Extremely poisonous berry that looks deceptively appetizing.",
						Identification: []string{
							"Dull, black, cherry-sized berries growing singly (not in clusters).",
							"Bell-shaped, purple-brown flowers.",
							"Large, oval, pointed leaves.",
							"Bushy plant 2-5 feet tall.",
						},
						Warnings: []string{
							"2-5 berries can kill a child. 10-20 can kill an adult.",
							"Berries are sweet-tasting, making them especially dangerous to children.",
							"Symptoms: dilated pupils, rapid heartbeat, hallucinations, seizures.",
						},
					},
				},
			},
		},
	}
}

// GetGroups returns all plant groups.
func (db *PlantDB) GetGroups() []PlantGroup {
	return db.Groups
}

// GetGroup returns a single group by ID.
func (db *PlantDB) GetGroup(id string) *PlantGroup {
	for _, g := range db.Groups {
		if g.ID == id {
			return &g
		}
	}
	return nil
}

// GetPlant returns a single plant by ID.
func (db *PlantDB) GetPlant(id string) *Plant {
	for _, g := range db.Groups {
		for _, p := range g.Plants {
			if p.ID == id {
				return &p
			}
		}
	}
	return nil
}

// SearchPlants searches plants by name, description, or habitat.
func (db *PlantDB) SearchPlants(query string) []Plant {
	var results []Plant
	q := toLower(query)
	for _, g := range db.Groups {
		for _, p := range g.Plants {
			if contains(toLower(p.CommonName), q) ||
				contains(toLower(p.ScientificName), q) ||
				contains(toLower(p.Description), q) ||
				contains(toLower(p.Habitat), q) {
				results = append(results, p)
			}
		}
	}
	if results == nil {
		results = make([]Plant, 0)
	}
	return results
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
