// Package galaxy generates the 8 galaxies of 256 systems and their markets.
package galaxy

import (
	"math"
	"strings"

	"heliosdrift/models"
)

var Pairs = []string{
	"AB", "OU", "SE", "IT", "IL", "ET", "ST", "ON", "LO", "NU", "TH", "NO",
	"AL", "LE", "XE", "GE", "ZA", "CE", "BI", "SO", "US", "ES", "AR", "MA",
	"IN", "DI", "RE", "A", "ER", "AT", "EN", "BE", "RA", "LA", "VE", "TI",
	"ED", "OR", "QU", "AN", "TE", "IS", "RI", "ON",
}

var EconomyNames = []string{"RICH INDUSTRIAL", "AVERAGE INDUSTRIAL", "POOR INDUSTRIAL",
	"MAINLY INDUSTRIAL", "MAINLY AGRICULTURAL", "RICH AGRICULTURAL",
	"AVERAGE AGRICULTURAL", "POOR AGRICULTURAL"}

var GovNames = []string{"ANARCHY", "FEUDAL", "MULTI-GOVERNMENT", "DICTATORSHIP",
	"COMMUNIST", "CONFEDERACY", "DEMOCRACY", "CORPORATE STATE"}

type Commodity struct {
	Name string
	Base float64
	Grad float64
	Qty  float64
	Unit string
}

var Commodities = []Commodity{
	{"FOOD", 19, 0.55, 14, "T"},
	{"TEXTILES", 20, 0.45, 14, "T"},
	{"RADIOACTIVES", 65, 0.30, 8, "T"},
	{"SLAVES", 40, -0.25, 6, "T"},
	{"LIQUOR/WINES", 83, 0.40, 8, "T"},
	{"LUXURIES", 196, -0.35, 6, "T"},
	{"NARCOTICS", 235, -0.30, 5, "T"},
	{"COMPUTERS", 154, -0.60, 8, "T"},
	{"MACHINERY", 117, -0.55, 8, "T"},
	{"ALLOYS", 78, -0.40, 10, "T"},
	{"FIREARMS", 124, 0.20, 6, "T"},
	{"FURS", 176, 0.35, 6, "T"},
	{"MINERALS", 32, -0.20, 12, "T"},
	{"GOLD", 97, 0.15, 6, "KG"},
	{"PLATINUM", 171, 0.10, 5, "KG"},
	{"GEM-STONES", 45, -0.10, 4, "KG"},
	{"ALIEN ITEMS", 53, -0.90, 2, "T"},
}

type System struct {
	ID         int
	Name       string
	X, Y       int
	Economy    int
	Government int
	Tech       int
	Pop        float64
	Prod       int
	Radius     int
	Kind       int
	Seed       uint32
}

type Market struct {
	Name  string
	Unit  string
	Price int
	Sell  int
	Qty   int
}

type Galaxy struct {
	Index   int
	Systems []*System
}

// Equipment kinds. A shop sells three different things and only one of them
// can be bought twice: freight you resell, fittings you keep, ammunition you
// spend. Getting that wrong charged players for a second mining laser they
// could never sell.
const (
	EquipUnique     = "unique"     // bolted on once, never sellable
	EquipConsumable = "consumable" // missiles and bombs, buy as many as you like
	EquipLaser      = "laser"      // one gun in one slot, better replaces worse
)

type Equipment struct {
	ID    string
	Name  string
	Price int
	Tech  int
	Kind  string // one of the Equip* constants
	// Needs is the equipment flag that must be fitted first, and Level is the
	// gun's place in the ladder. Both only matter for some kinds.
	Needs string
	Level int
}

var EquipmentList = []Equipment{
	{ID: "missile", Name: "MISSILE", Price: 30, Tech: 1, Kind: EquipConsumable},
	{ID: "largebay", Name: "LARGE CARGO BAY", Price: 400, Tech: 2, Kind: EquipUnique},
	{ID: "ecm", Name: "E.C.M. SYSTEM", Price: 600, Tech: 3, Kind: EquipUnique},
	{ID: "fuelscoop", Name: "FUEL SCOOP", Price: 525, Tech: 4, Kind: EquipUnique},
	{ID: "laserpulse", Name: "PULSE LASER", Price: 400, Tech: 2, Kind: EquipLaser, Level: 0},
	{ID: "laserbeam", Name: "BEAM LASER", Price: 1000, Tech: 5, Kind: EquipLaser, Level: 1},
	{ID: "lasermilitary", Name: "MILITARY LASER", Price: 6000, Tech: 10, Kind: EquipLaser, Level: 2},
	{ID: "laserplasma", Name: "PLASMA ACCELERATOR", Price: 18000, Tech: 12, Kind: EquipLaser, Level: 3},
	{ID: "mininglaser", Name: "MINING LASER", Price: 800, Tech: 4, Kind: EquipUnique},
	{ID: "miningrig", Name: "HEAVY MINING RIG", Price: 2600, Tech: 5, Kind: EquipUnique, Needs: "mining"},
	{ID: "shieldboost1", Name: "SHIELD BOOSTER", Price: 4000, Tech: 6, Kind: EquipUnique},
	{ID: "shieldboost2", Name: "HEAVY SHIELDING", Price: 12000, Tech: 11, Kind: EquipUnique, Needs: "shieldboost1"},
	{ID: "cargobay2", Name: "CARGO EXPANSION", Price: 9000, Tech: 7, Kind: EquipUnique, Needs: "largebay"},
	{ID: "bomb", Name: "ENERGY BOMB", Price: 900, Tech: 7, Kind: EquipConsumable},
	{ID: "escape", Name: "ESCAPE CAPSULE", Price: 1000, Tech: 6, Kind: EquipUnique},
	{ID: "extraenergy", Name: "EXTRA ENERGY UNIT", Price: 1500, Tech: 8, Kind: EquipUnique},
	{ID: "dockcomp", Name: "DOCKING COMPUTER", Price: 1500, Tech: 9, Kind: EquipUnique},
	{ID: "galdrive", Name: "GALACTIC HYPERDRIVE", Price: 5000, Tech: 11, Kind: EquipUnique},
}

// EquipByID finds a catalogue entry.
func EquipByID(id string) *Equipment {
	for i := range EquipmentList {
		if EquipmentList[i].ID == id {
			return &EquipmentList[i]
		}
	}
	return nil
}

// EquipNameForFlag gives the shop its words: "mining" is the flag, and the
// thing the player has to buy first is a MINING LASER.
func EquipNameForFlag(flag string) string {
	switch flag {
	case "mining":
		return "MINING LASER"
	}
	if e := EquipByID(flag); e != nil {
		return e.Name
	}
	return strings.ToUpper(flag)
}

const FuelPricePerLY = 2

var cache = map[int]*Galaxy{}

// Galaxies is the number of galactic charts in a saved commander's universe.
const Galaxies = 8

func Get(index int) *Galaxy {
	index = ((index % Galaxies) + Galaxies) % Galaxies
	if g, ok := cache[index]; ok {
		return g
	}
	g := &Galaxy{Index: index, Systems: generate(index)}
	cache[index] = g
	return g
}

func makeName(r *models.Rng) string {
	n := 2 + int(r.Float()*3)
	s := ""
	for i := 0; i < n; i++ {
		s += Pairs[int(r.Float()*float64(len(Pairs)))]
	}
	if len(s) > 9 {
		s = s[:9]
	}
	if len(s) < 4 {
		s += Pairs[int(r.Float()*float64(len(Pairs)))]
	}
	return s
}

func generate(index int) []*System {
	r := models.NewRng(uint32(0x5a4a27) + uint32(index)*0x9e3779b9)
	systems := make([]*System, 0, 256)
	used := map[string]bool{}
	for i := 0; i < 256; i++ {
		name := makeName(r)
		guard := 0
		for used[name] && guard < 40 {
			name = makeName(r)
			guard++
		}
		if used[name] {
			name = name + string(rune('0'+i%10))
		}
		used[name] = true

		economy := int(r.Float() * 8)
		government := int(r.Float() * 8)
		if government >= 6 && r.Float() < 0.5 {
			economy = int(r.Float() * 4)
		}
		if government <= 1 && r.Float() < 0.5 {
			economy = 4 + int(r.Float()*4)
		}
		tech := int(math.Round(1 + float64(7-economy)*0.6 + float64(government)*0.35 + r.Float()*4))
		if tech < 1 {
			tech = 1
		}
		if tech > 15 {
			tech = 15
		}
		pop := math.Round((0.4+r.Float()*11)*10) / 10
		prodEconomy := 0.6
		if economy <= 3 {
			prodEconomy = 1
		}
		prod := int(prodEconomy*float64(tech)*pop*4 + r.Float()*400)
		radius := 2000 + int(r.Float()*7000)
		systems = append(systems, &System{
			ID: i, Name: name,
			X: int(r.Float() * 256), Y: int(r.Float() * 256),
			Economy: economy, Government: government, Tech: tech,
			Pop: pop, Prod: prod, Radius: radius,
			Kind: models.PlanetKindFor(uint32(i*977 + index*31 + 5)),
			Seed: uint32(i*2654435761 + index*40503),
		})
	}
	// spread out
	for k := 0; k < 256; k++ {
		s0 := systems[k]
		for j := 0; j < 8; j++ {
			other := systems[int(r.Float()*256)]
			if other == s0 {
				continue
			}
			dx := float64(s0.X - other.X)
			dy := float64(s0.Y - other.Y)
			if math.Hypot(dx, dy) < 2.4 {
				s0.X = (other.X + 3 + int(r.Float()*9)) % 256
				s0.Y = (other.Y + 3 + int(r.Float()*9)) % 256
			}
		}
	}
	// de-dupe
	for pass := 0; pass < 6; pass++ {
		moved := 0
		for a := 0; a < len(systems); a++ {
			for b := a + 1; b < len(systems); b++ {
				A, B := systems[a], systems[b]
				dx, dy := A.X-B.X, A.Y-B.Y
				if dx*dx+dy*dy < 4 {
					B.X = (B.X + 5 + int(r.Float()*7)) % 256
					B.Y = (B.Y + 5 + int(r.Float()*7)) % 256
					moved++
				}
			}
		}
		if moved == 0 {
			break
		}
	}
	ensureConnected(systems, r)
	return systems
}

func lyBetween(a, b *System) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Hypot(dx, dy) / 4
}

func ensureConnected(systems []*System, r *models.Rng) {
	for pass := 0; pass < 4; pass++ {
		seen := map[int]bool{0: true}
		queue := []int{0}
		for len(queue) > 0 {
			cur := systems[queue[0]]
			queue = queue[1:]
			for _, s := range systems {
				if seen[s.ID] {
					continue
				}
				if lyBetween(cur, s) <= 6.8 {
					seen[s.ID] = true
					queue = append(queue, s.ID)
				}
			}
		}
		if len(seen) == len(systems) {
			return
		}
		for _, s := range systems {
			if seen[s.ID] {
				continue
			}
			host := systems[int(r.Float()*float64(len(systems)))]
			if !seen[host.ID] {
				host = systems[0]
			}
			ang := r.Float() * 2 * math.Pi
			dist := 3 + r.Float()*3
			s.X = int((float64(host.X) + math.Cos(ang)*dist*4 + 256)) % 256
			s.Y = int((float64(host.Y) + math.Sin(ang)*dist*4 + 256)) % 256
		}
	}
}

// Distance returns the jump distance in light years.
func Distance(a, b *System) float64 { return lyBetween(a, b) }

func FindByName(g *Galaxy, name string) *System {
	for _, s := range g.Systems {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// MakeMarket builds the price list for a system; salt varies per docking.
func MakeMarket(sys *System, salt uint32) []Market {
	r := models.NewRng(sys.Seed ^ (salt * 2654435761))
	ind := float64(7-sys.Economy) / 7
	out := make([]Market, len(Commodities))
	for i, c := range Commodities {
		level := 1 + c.Grad*(ind-0.45)*1.15
		if level < 0.35 {
			level = 0.35
		}
		fluct := 0.88 + r.Float()*0.26
		price := int(math.Max(1, math.Round(c.Base*level*fluct)))
		q := int(math.Round(c.Qty * (0.45 + r.Float()*1.5) * (1.15 - 0.45*c.Grad*(ind-0.45))))
		if q < 0 {
			q = 0
		}
		sell := int(math.Max(1, math.Round(float64(price)*0.80)))
		out[i] = Market{Name: c.Name, Unit: c.Unit, Price: price, Sell: sell, Qty: q}
	}
	return out
}

func AvailableEquipment(sys *System) []Equipment {
	var out []Equipment
	for _, e := range EquipmentList {
		if sys.Tech >= e.Tech {
			out = append(out, e)
		}
	}
	return out
}

var Ranks = []string{"HARMLESS", "MOSTLY HARMLESS", "POOR", "AVERAGE", "ABOVE AVERAGE",
	"COMPETENT", "DANGEROUS", "DEADLY", "LEGEND"}

var rankThresholds = []int{0, 8, 16, 32, 64, 128, 512, 2560, 6400}

func RankFor(kills int) string {
	r := 0
	for i, t := range rankThresholds {
		if kills >= t {
			r = i
		}
	}
	return Ranks[r]
}
