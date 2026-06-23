package medical

// Category represents a medical reference category.
type Category struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Description string  `json:"description"`
	Entries     []Entry `json:"entries"`
}

// Entry is a single medical reference entry.
type Entry struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Severity string   `json:"severity"` // critical, warning, info
	Summary  string   `json:"summary"`
	Steps    []string `json:"steps"`
	Warnings []string `json:"warnings,omitempty"`
}

// MedicalDB holds the built-in offline medical reference data.
type MedicalDB struct {
	Categories []Category
}

// NewMedicalDB creates the medical reference database with embedded survival medical knowledge.
func NewMedicalDB() *MedicalDB {
	return &MedicalDB{
		Categories: []Category{
			{
				ID: "bleeding", Name: "Bleeding Control", Icon: "droplet",
				Description: "Emergency hemorrhage control procedures",
				Entries: []Entry{
					{
						ID: "severe-bleeding", Title: "Severe Bleeding / Hemorrhage", Severity: "critical",
						Summary: "Life-threatening blood loss requiring immediate intervention.",
						Steps: []string{
							"Apply direct pressure with a clean cloth or bandage.",
							"If bleeding through, add more layers — do NOT remove the first one.",
							"Elevate the wound above the heart if possible.",
							"If direct pressure fails on a limb, apply a tourniquet 2-3 inches above the wound.",
							"Note the time the tourniquet was applied.",
							"DO NOT remove the tourniquet once applied.",
							"Monitor for shock: pale skin, rapid pulse, confusion.",
						},
						Warnings: []string{
							"A tourniquet is a LAST RESORT for limbs only.",
							"Never apply a tourniquet to the neck, head, or torso.",
						},
					},
					{
						ID: "nosebleed", Title: "Nosebleed", Severity: "info",
						Summary: "Common nosebleed management.",
						Steps: []string{
							"Sit upright and lean slightly forward.",
							"Pinch the soft part of the nose firmly for 10-15 minutes.",
							"Breathe through the mouth.",
							"Apply a cold compress to the bridge of the nose.",
							"Do NOT tilt the head back — blood can enter the airway.",
						},
					},
				},
			},
			{
				ID: "cpr", Name: "CPR & Choking", Icon: "heart-pulse",
				Description: "Cardiopulmonary resuscitation and airway obstruction",
				Entries: []Entry{
					{
						ID: "adult-cpr", Title: "Adult CPR", Severity: "critical",
						Summary: "CPR for unresponsive adult with no normal breathing.",
						Steps: []string{
							"Ensure the scene is safe.",
							"Check responsiveness: tap shoulders and shout.",
							"Call for help if available.",
							"Place the heel of one hand on the center of the chest (lower half of sternum).",
							"Place the other hand on top, interlock fingers.",
							"Push hard and fast: 2 inches deep, 100-120 compressions per minute.",
							"After 30 compressions, give 2 rescue breaths (tilt head back, lift chin, seal mouth).",
							"Continue 30:2 cycle until help arrives or the person recovers.",
						},
						Warnings: []string{
							"Do NOT stop CPR unless the person starts breathing normally.",
							"Compression-only CPR is acceptable if unable to give rescue breaths.",
						},
					},
					{
						ID: "choking-adult", Title: "Choking (Adult)", Severity: "critical",
						Summary: "Airway obstruction in a conscious adult.",
						Steps: []string{
							"Ask: 'Are you choking?' — if they cannot speak or cough, act immediately.",
							"Stand behind the person, wrap your arms around their waist.",
							"Make a fist with one hand, place it above the navel, below the ribcage.",
							"Grasp the fist with the other hand.",
							"Deliver quick, upward thrusts (Heimlich maneuver).",
							"Repeat until the object is expelled or the person becomes unresponsive.",
							"If unresponsive: begin CPR, checking the airway each cycle.",
						},
					},
				},
			},
			{
				ID: "fractures", Name: "Fractures & Splinting", Icon: "bone",
				Description: "Broken bone identification and immobilization",
				Entries: []Entry{
					{
						ID: "fracture-general", Title: "General Fracture Management", Severity: "warning",
						Summary: "Suspected broken bone in a limb.",
						Steps: []string{
							"Do NOT try to realign the bone.",
							"Immobilize the injury in the position found.",
							"Splint the joint above AND below the fracture.",
							"Use rigid materials: sticks, boards, rolled magazines.",
							"Pad the splint with cloth for comfort.",
							"Secure with bandages or strips of cloth — not too tight.",
							"Check circulation below the splint: pulse, color, sensation.",
							"Apply cold pack if available (wrapped in cloth) for swelling.",
						},
						Warnings: []string{
							"If bone is protruding (open fracture), cover with a sterile dressing. Do NOT push it back in.",
							"Check for signs of compartment syndrome: severe pain, numbness, pale/blue skin.",
						},
					},
				},
			},
			{
				ID: "burns", Name: "Burn Treatment", Icon: "flame",
				Description: "Thermal, chemical, and electrical burn first aid",
				Entries: []Entry{
					{
						ID: "thermal-burn", Title: "Thermal Burns", Severity: "warning",
						Summary: "Burns from fire, hot liquids, or hot surfaces.",
						Steps: []string{
							"Remove the person from the heat source.",
							"Cool the burn with cool (not cold) running water for at least 20 minutes.",
							"Remove jewelry and clothing near the burn BEFORE swelling begins.",
							"Cover with a clean, non-stick dressing or cling wrap.",
							"Do NOT apply butter, oil, toothpaste, or ice.",
							"Do NOT burst blisters.",
							"For large burns or burns to face/hands/genitals, seek urgent medical aid.",
						},
						Warnings: []string{
							"Burns covering >10% of body surface area are life-threatening.",
							"Inhalation burns (from smoke/steam) can cause airway swelling — monitor breathing closely.",
						},
					},
				},
			},
			{
				ID: "environmental", Name: "Environmental Emergencies", Icon: "thermometer",
				Description: "Hypothermia, heat stroke, dehydration, and exposure",
				Entries: []Entry{
					{
						ID: "hypothermia", Title: "Hypothermia", Severity: "critical",
						Summary: "Core body temperature drops below 35°C (95°F).",
						Steps: []string{
							"Move to shelter or out of the wind/cold.",
							"Remove wet clothing gently.",
							"Insulate from the ground with a sleeping pad, branches, or dry material.",
							"Warm the core FIRST: warm blankets, body heat from another person.",
							"If conscious and able to swallow, give warm (not hot) sweet drinks.",
							"Do NOT rub the skin, apply direct heat, or give alcohol.",
							"Handle gently — rough movement can cause cardiac arrest in severe hypothermia.",
						},
						Warnings: []string{
							"Severe hypothermia victims may appear dead. Continue warming — 'They're not dead until they're warm and dead.'",
						},
					},
					{
						ID: "heat-stroke", Title: "Heat Stroke", Severity: "critical",
						Summary: "Core body temperature rises above 40°C (104°F). Life-threatening.",
						Steps: []string{
							"Move to shade or the coolest area available.",
							"Remove excess clothing.",
							"Cool aggressively: wet the skin and fan, apply cold packs to neck, armpits, and groin.",
							"If conscious, give small sips of water.",
							"Do NOT give fluids if unconscious.",
							"Monitor breathing and be prepared to start CPR.",
						},
					},
					{
						ID: "dehydration", Title: "Dehydration", Severity: "warning",
						Summary: "Fluid loss from sweating, vomiting, diarrhea, or insufficient intake.",
						Steps: []string{
							"Drink small, frequent sips of water.",
							"If available, use oral rehydration salts (ORS).",
							"DIY ORS: 1 liter water + 6 teaspoons sugar + ½ teaspoon salt.",
							"Rest in shade and avoid exertion.",
							"Monitor urine color — aim for pale yellow.",
						},
					},
				},
			},
			{
				ID: "bites-stings", Name: "Bites & Stings", Icon: "bug",
				Description: "Snake bites, insect stings, and animal bites",
				Entries: []Entry{
					{
						ID: "snake-bite", Title: "Snake Bite", Severity: "critical",
						Summary: "Suspected venomous snake bite.",
						Steps: []string{
							"Keep the person calm and still — movement spreads venom faster.",
							"Immobilize the bitten limb with a splint, keep at or below heart level.",
							"Remove rings, watches, tight clothing near the bite before swelling.",
							"Note the time of the bite and the snake's appearance if safely possible.",
							"Do NOT cut, suck, tourniquet, or apply ice to the wound.",
							"Do NOT give alcohol or aspirin.",
							"Transport to medical care as quickly and calmly as possible.",
						},
						Warnings: []string{
							"Many 'non-venomous' bites still carry infection risk.",
							"Pressure immobilization bandage is recommended for some snake types (e.g. elapids).",
						},
					},
				},
			},
			{
				ID: "water", Name: "Water Safety", Icon: "droplets",
				Description: "Purification, sourcing, and water-borne illness prevention",
				Entries: []Entry{
					{
						ID: "water-purification", Title: "Water Purification Methods", Severity: "info",
						Summary: "Making water safe to drink in the field.",
						Steps: []string{
							"BOILING: Bring water to a rolling boil for at least 1 minute (3 minutes above 2000m).",
							"CHEMICAL: Add 2 drops of household bleach (5-6% sodium hypochlorite) per liter. Wait 30 minutes.",
							"FILTER: Use a commercial filter rated to 0.2 microns, or improvise with layers of sand, charcoal, and gravel.",
							"UV: If you have a UV purifier, follow manufacturer's instructions.",
							"SOLAR: Fill clear PET bottles, place in direct sunlight for 6+ hours (SODIS method).",
							"Always filter cloudy water through cloth before chemical/UV treatment.",
						},
					},
				},
			},
		},
	}
}

// GetCategories returns all medical reference categories.
func (mdb *MedicalDB) GetCategories() []Category {
	return mdb.Categories
}

// GetCategory returns a single category by ID.
func (mdb *MedicalDB) GetCategory(id string) *Category {
	for _, c := range mdb.Categories {
		if c.ID == id {
			return &c
		}
	}
	return nil
}

// GetEntry returns a single entry by ID, searching all categories.
func (mdb *MedicalDB) GetEntry(id string) *Entry {
	for _, c := range mdb.Categories {
		for _, e := range c.Entries {
			if e.ID == id {
				return &e
			}
		}
	}
	return nil
}
