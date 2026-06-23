package skilltrees

// Skill represents a single survival skill with step-by-step instructions.
type Skill struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Difficulty    string   `json:"difficulty"` // beginner, intermediate, advanced
	TimeEstimate  string   `json:"time_estimate"`
	Prerequisites []string `json:"prerequisites,omitempty"`
	Summary       string   `json:"summary"`
	Steps         []Step   `json:"steps"`
	Tips          []string `json:"tips,omitempty"`
}

// Step is an individual instruction step within a skill.
type Step struct {
	Order       int    `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// SkillCategory groups related survival skills.
type SkillCategory struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Description string  `json:"description"`
	Skills      []Skill `json:"skills"`
}

// SkillTreeDB holds all embedded survival skill data.
type SkillTreeDB struct {
	Categories []SkillCategory
}

// NewSkillTreeDB creates the survival skill reference database.
func NewSkillTreeDB() *SkillTreeDB {
	return &SkillTreeDB{
		Categories: []SkillCategory{
			{
				ID: "fire", Name: "Fire Making", Icon: "flame",
				Description: "Essential fire-starting techniques for warmth, cooking, and signaling",
				Skills: []Skill{
					{
						ID: "friction-fire", Name: "Bow Drill Fire", Difficulty: "advanced",
						TimeEstimate: "30-60 min", Summary: "Start a fire using friction from a bow drill — the most reliable primitive method.",
						Steps: []Step{
							{1, "Gather Materials", "Find a dry, soft wood for the fireboard and spindle (e.g., willow, cedar, cottonwood). Cut a straight, sturdy stick for the bow and a harder piece for the handhold."},
							{2, "Carve the Fireboard", "Flatten one side of a forearm-length piece of wood. Carve a small depression near one edge."},
							{3, "Cut the Notch", "Carve a V-shaped notch from the edge of the board into the depression. This collects the hot dust."},
							{4, "Make the Spindle", "Carve a straight, thumb-thick stick about 18 inches long. Round one end, point the other slightly."},
							{5, "String the Bow", "Tie cordage (paracord, shoelace, or natural fiber) to a curved branch with slight tension."},
							{6, "Prepare the Tinder Bundle", "Gather dry, fine material: inner bark, dried grass, cattail fluff. Form a bird's-nest shape."},
							{7, "Drill", "Loop the spindle in the bowstring. Place the rounded end in the fireboard depression. Press down with the handhold and saw the bow back and forth rapidly."},
							{8, "Catch the Ember", "When smoke pours from the notch, stop. Carefully tip the glowing coal dust into your tinder bundle."},
							{9, "Blow to Flame", "Gently blow on the ember in the tinder bundle, gradually increasing intensity until flames appear."},
						},
						Tips: []string{
							"The fireboard and spindle should be the same type of soft, dry wood.",
							"If the dust is light-colored, add more downward pressure. If dark, add more speed.",
							"Practice this skill BEFORE you need it in an emergency.",
						},
					},
					{
						ID: "ferro-rod", Name: "Ferro Rod Fire Starting", Difficulty: "beginner",
						TimeEstimate: "5-10 min", Summary: "Use a ferrocerium rod to throw sparks onto tinder.",
						Steps: []Step{
							{1, "Prepare Tinder", "Gather dry, fine tinder: cotton balls (with petroleum jelly if available), birch bark, dry grass, or char cloth."},
							{2, "Build a Fire Lay", "Arrange small kindling (pencil-thick sticks) in a teepee or log-cabin shape over the tinder spot."},
							{3, "Position the Rod", "Hold the ferro rod close to and angled toward the tinder bundle."},
							{4, "Strike", "Scrape the striker (or spine of a knife) firmly down the rod, directing sparks onto the tinder."},
							{5, "Nurture the Flame", "Once tinder catches, gently blow and add progressively larger sticks."},
						},
						Tips: []string{
							"Scrape the black coating off a new ferro rod first.",
							"Push the rod BACK instead of the striker forward to avoid disturbing tinder.",
						},
					},
				},
			},
			{
				ID: "shelter", Name: "Shelter Building", Icon: "home",
				Description: "Protection from elements using natural and improvised materials",
				Skills: []Skill{
					{
						ID: "debris-hut", Name: "Debris Hut", Difficulty: "intermediate",
						TimeEstimate: "2-4 hours", Summary: "Build an insulated shelter from forest debris — effective in cold conditions.",
						Steps: []Step{
							{1, "Find a Ridgepole", "Find a sturdy, straight pole about 9-12 feet long. It should support your weight at the thick end."},
							{2, "Set the Framework", "Prop one end on a stump, rock, or Y-shaped stick about 3 feet high. The other end rests on the ground. You should be able to sit up inside."},
							{3, "Add Ribbing", "Lean sticks along both sides of the ridgepole at 45° angles, spaced about a fist apart. This forms the skeleton."},
							{4, "Layer Debris", "Pile leaves, pine needles, grass, and ferns at least 3 feet thick over the ribs. Start from the bottom up like shingles."},
							{5, "Insulate the Floor", "Fill the inside with a thick layer of dry leaves or grass for ground insulation. This is critical — you lose more heat to the ground than to the air."},
							{6, "Block the Entrance", "Stuff a large pile of leaves into the entrance that you can pull behind you. Smaller entrance = warmer shelter."},
						},
						Tips: []string{
							"Build it just big enough for your body — excess space wastes your body heat.",
							"A debris hut with 3 feet of insulation can keep you warm in below-freezing temps with no fire.",
							"Building this shelter takes 2-4 hours. Start EARLY before dark.",
						},
					},
					{
						ID: "tarp-shelter", Name: "Tarp A-Frame Shelter", Difficulty: "beginner",
						TimeEstimate: "15-30 min", Summary: "Quick, effective shelter using a tarp or poncho.",
						Steps: []Step{
							{1, "Find Two Anchor Points", "Locate two trees about 8-10 feet apart."},
							{2, "Tie the Ridgeline", "Run paracord between the trees at about 4 feet high. Pull taut and tie off."},
							{3, "Drape the Tarp", "Hang the tarp centerline over the cord, creating two equal sloping sides."},
							{4, "Stake the Corners", "Stake or weight the tarp corners at 45° angles to the ground. Leave one end open as an entrance."},
							{5, "Angle for Drainage", "Angle the shelter slightly downhill so rain runs away from your sleeping area."},
						},
					},
				},
			},
			{
				ID: "water", Name: "Water Procurement", Icon: "droplets",
				Description: "Finding, collecting, and purifying water in the wild",
				Skills: []Skill{
					{
						ID: "solar-still", Name: "Solar Still", Difficulty: "intermediate",
						TimeEstimate: "1-2 hours setup", Summary: "Extract water from the ground using solar evaporation.",
						Steps: []Step{
							{1, "Dig a Hole", "Dig a bowl-shaped hole about 3 feet wide and 2 feet deep in a sunny location."},
							{2, "Place a Container", "Put a clean cup or container in the center of the hole."},
							{3, "Add Vegetation", "Place non-toxic green vegetation around the container (not in it) to increase moisture output."},
							{4, "Cover with Plastic", "Stretch clear plastic sheeting over the hole. Seal edges with dirt or rocks."},
							{5, "Add Weight", "Place a small stone on the plastic directly above the container so it dips to a point."},
							{6, "Wait", "Sun heats the ground → moisture evaporates → condenses on plastic → drips into your container."},
						},
						Tips: []string{
							"Yields about 1-2 cups per day — not enough alone, but can supplement other sources.",
							"Morning dew collection and plant transpiration bags often yield more water per effort.",
						},
					},
				},
			},
			{
				ID: "navigation", Name: "Land Navigation", Icon: "compass",
				Description: "Finding direction and navigating without GPS or electronics",
				Skills: []Skill{
					{
						ID: "shadow-stick", Name: "Shadow Stick Compass", Difficulty: "beginner",
						TimeEstimate: "30 min", Summary: "Determine east-west direction using the sun and a stick.",
						Steps: []Step{
							{1, "Place a Stick", "Push a straight stick (about 3 feet) vertically into flat, clear ground."},
							{2, "Mark the First Shadow", "Place a stone or mark at the tip of the shadow cast by the stick."},
							{3, "Wait 15-30 Minutes", "The shadow tip will move. Wait at least 15 minutes."},
							{4, "Mark the Second Shadow", "Place another stone at the new shadow tip position."},
							{5, "Draw the Line", "Draw a line between the two marks. This line runs approximately east-west. The first mark is WEST, the second is EAST."},
							{6, "Find North", "Stand with the first mark (west) on your left and the second (east) on your right. You are facing approximately north (in the Northern Hemisphere)."},
						},
						Tips: []string{
							"The longer you wait between marks, the more accurate the reading.",
							"In the Southern Hemisphere, you'll be facing south instead.",
						},
					},
				},
			},
			{
				ID: "signaling", Name: "Signaling & Rescue", Icon: "radio",
				Description: "Attracting attention and communicating distress",
				Skills: []Skill{
					{
						ID: "signal-fire", Name: "Signal Fire", Difficulty: "beginner",
						TimeEstimate: "30 min", Summary: "Build a fire designed to attract rescue attention.",
						Steps: []Step{
							{1, "Choose Location", "Build on high ground, clearings, or near water where smoke is visible from the air."},
							{2, "Build Three Fires", "The international distress signal is THREE fires in a triangle, about 100 feet apart."},
							{3, "Keep Ready to Light", "Prepare the fires with dry kindling so they can be lit quickly when aircraft is spotted."},
							{4, "Create Smoke", "Once lit, add green branches, wet leaves, or rubber to create thick white smoke (daytime) or keep flames visible (nighttime)."},
							{5, "Signal with Pattern", "If you have a single fire, cover and uncover it to create smoke puffs (3 puffs = distress)."},
						},
					},
					{
						ID: "ground-signals", Name: "Ground-to-Air Signals", Difficulty: "beginner",
						TimeEstimate: "1-2 hours", Summary: "Create large visual signals visible to search aircraft.",
						Steps: []Step{
							{1, "Find Contrasting Material", "Use rocks, logs, clothing, or trampled snow to create contrast against the ground."},
							{2, "Make Symbols Large", "Each symbol should be at least 10 feet tall and 3 feet wide. Bigger is always better."},
							{3, "Use Standard Symbols", "V = Need assistance. X = Need medical help. → (arrow) = Traveling this direction. I = Need supplies."},
							{4, "Add Shadow", "Build symbols with depth (stack rocks, pile branches) so shadows make them visible from altitude."},
						},
					},
				},
			},
			{
				ID: "food", Name: "Food Procurement", Icon: "utensils",
				Description: "Finding and preparing food in survival situations",
				Skills: []Skill{
					{
						ID: "snare-trap", Name: "Simple Snare Trap", Difficulty: "intermediate",
						TimeEstimate: "15-30 min per snare", Summary: "Build a basic wire/cordage snare to catch small game.",
						Steps: []Step{
							{1, "Make the Noose", "Form a small loop (fist-sized for rabbits) from wire, paracord, or strong natural cordage. Create a sliding knot."},
							{2, "Find a Game Trail", "Look for animal tracks, droppings, and worn paths between cover and water sources."},
							{3, "Set the Snare", "Hang the noose at head-height of the target animal across the trail, suspended from a branch or stake."},
							{4, "Anchor Securely", "Tie the other end firmly to a stake, tree, or drag log heavy enough that the animal can't escape."},
							{5, "Set Multiple Snares", "Set 6-12 snares to increase your odds. Check every 12-24 hours."},
							{6, "Camouflage", "Use natural materials to funnel the animal into the snare. Remove human scent with mud."},
						},
						Tips: []string{
							"In a survival situation, snares work while you sleep — they're your most efficient food-gathering tool.",
							"Always know local laws. This is for emergency survival situations only.",
						},
					},
				},
			},
		},
	}
}

// GetCategories returns all skill categories.
func (db *SkillTreeDB) GetCategories() []SkillCategory {
	return db.Categories
}

// GetCategory returns a single category by ID.
func (db *SkillTreeDB) GetCategory(id string) *SkillCategory {
	for _, c := range db.Categories {
		if c.ID == id {
			return &c
		}
	}
	return nil
}

// GetSkill returns a single skill by ID.
func (db *SkillTreeDB) GetSkill(id string) *Skill {
	for _, c := range db.Categories {
		for _, s := range c.Skills {
			if s.ID == id {
				return &s
			}
		}
	}
	return nil
}
