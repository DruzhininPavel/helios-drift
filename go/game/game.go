// Package game holds the simulation: flight, combat, docking, hyperspace,
// trading and ship AI.
package game

import (
	"math"
	"math/rand"

	"heliosdrift/galaxy"
	"heliosdrift/models"
	"heliosdrift/render"
)

const (
	ViewW = 320
	ViewH = 150

	MaxSpeed    = 420
	PitchRate   = 1.55
	RollRate    = 2.35
	MaxEnergy   = 100
	MaxShield   = 100
	CargoBase   = 20
	FuelMax     = 7.0
	MissilesMax = 4
	RadarRange  = 26000
)

var StationSize = models.StationSize
var SlotW = models.StationSlotW
var SlotH = models.StationSlotH

var LaserRange = [4]float64{2600, 3400, 4600, 2000}
var LaserDamage = [4]float64{9, 16, 26, 7}
var LaserHeat = [4]float64{0.055, 0.10, 0.13, 0.05}
var LaserCooldown = [4]float64{0.14, 0.075, 0.06, 0.2}
var LaserName = [4]string{"PULSE LASER", "BEAM LASER", "MILITARY LASER", "MINING LASER"}

type ShipDef struct {
	Mesh     string
	Radius   float64
	Speed    float64
	Energy   float64
	Shield   float64
	Bounty   int
	Cargo    int
	Laser    int
	FireRate float64
	Turn     float64
	Role     string
}

var Ships = map[string]ShipDef{
	"sidewinder":  {"sidewinder", 26, 380, 34, 0, 5, 2, 1, 0.55, 1.7, "pirate"},
	"krait":       {"krait", 27, 340, 60, 0, 10, 3, 1, 0.5, 1.6, "pirate"},
	"mamba":       {"mamba", 28, 420, 48, 10, 15, 3, 1, 0.45, 1.8, "pirate"},
	"adder":       {"adder", 24, 260, 40, 0, 8, 4, 0, 0.8, 1.3, "trader"},
	"cobra3":      {"cobra3", 34, 350, 80, 20, 20, 6, 1, 0.5, 1.6, "pirate"},
	"python":      {"python", 46, 200, 140, 30, 30, 12, 1, 0.7, 1.0, "trader"},
	"viper":       {"viper", 28, 480, 70, 25, 15, 2, 1, 0.4, 2.0, "police"},
	"asp":         {"asp", 32, 400, 100, 25, 40, 4, 2, 0.4, 1.8, "pirate"},
	"moray":       {"moray", 24, 320, 44, 0, 12, 3, 1, 0.6, 1.5, "pirate"},
	"ferdelance":  {"ferdelance", 30, 480, 100, 30, 40, 4, 2, 0.35, 2.1, "pirate"},
	"thargoid":    {"thargoid", 36, 400, 260, 40, 100, 0, 2, 0.35, 1.6, "thargoid"},
	"shuttle":     {"shuttle", 22, 200, 26, 0, 0, 2, 0, 1.0, 1.0, "trader"},
	"gecko":       {"gecko", 26, 380, 34, 0, 8, 2, 1, 0.55, 1.7, "pirate"},
	"cobra1":      {"cobra1", 30, 340, 55, 10, 15, 5, 1, 0.55, 1.6, "pirate"},
	"worm":        {"worm", 24, 220, 30, 0, 4, 4, 0, 0.9, 1.2, "trader"},
	"boa":         {"boa", 44, 180, 180, 30, 30, 14, 1, 0.8, 0.9, "trader"},
	"anaconda":    {"anaconda", 56, 160, 250, 40, 50, 20, 2, 0.85, 0.8, "trader"},
	"transporter": {"transporter", 26, 200, 40, 0, 0, 4, 0, 1.0, 1.1, "trader"},
	"thargon":     {"thargon", 12, 420, 14, 0, 5, 0, 1, 0.35, 2.4, "thargoid"},
	"constrictor": {"constrictor", 34, 500, 120, 30, 0, 0, 2, 0.4, 2.0, "pirate"},
}

type Message struct {
	Text string
	T    float64
	Life float64
	Kind string
}

type Commander struct {
	Name         string
	Cash         float64
	Fuel         float64
	Missiles     int
	Cargo        map[int]int
	CargoCap     int
	Laser        int
	Kills        int
	Legal        int
	Galaxy       int
	System       int
	TargetSystem int
	Equip        map[string]bool
	Bombs        int
	Deaths       int
	Mission      Mission       `json:"mission"`
	MissionFlags int           `json:"missionFlags"`
	Contracts    []Contract    `json:"contracts"`
	Day          int           `json:"day"`
	Pax          int           `json:"pax"`
	Mail         int           `json:"mail"`
	Ship         string        `json:"ship"`
	Mode         string        `json:"mode"`
	Campaign     CampaignState `json:"campaign"`
}

func NewCommander() Commander {
	return Commander{
		Name: "JAMESON", Cash: 100, Fuel: 7, Missiles: 3,
		Cargo: map[int]int{}, CargoCap: CargoBase, Equip: map[string]bool{},
		Ship: "cobra3", Mode: ModeAdvanced,
	}
}

type Player struct {
	Cabin        float64
	Pos, Vel     render.Vec3
	Mat          render.Mat3
	Speed        float64
	Energy       float64
	ShieldF      float64
	ShieldA      float64
	LaserTemp    float64
	Cooldown     float64
	Invulnerable float64
}

type Entity struct {
	Kind         string
	Type         string
	Mesh         *render.Mesh
	Pos          render.Vec3
	Mat          render.Mat3
	Speed        float64
	MaxSpeed     float64
	Energy       float64
	MaxEnergy    float64
	Shield       float64
	MaxShield    float64
	Radius       float64
	Bounty       int
	CargoType    int
	Behaviour    string
	Hostile      bool
	Laser        int
	FireRate     float64
	Turn         float64
	Cooldown     float64
	Age          float64
	Life         float64
	CanMissile   bool
	Spin         render.Vec3
	TargetPlayer bool
	TargetEnt    *Entity
	Dead         bool
	Mission      bool
	CampaignTag  string
}

type Planet struct {
	Pos    render.Vec3
	Radius float64
	Seed   uint32
	Kind   int
}

type Sun struct {
	Pos    render.Vec3
	Radius float64
}

type Star struct {
	Dir    render.Vec3
	Bright float64
	Size   int
}

type Dust struct{ Pos render.Vec3 }

type World struct {
	Witchspace   bool
	Hermit       bool
	Sys          *galaxy.System
	Planet       Planet
	Sun          Sun
	Station      *Entity
	Entities     []*Entity
	Stars        []Star
	Dust         []Dust
	Nebula       *render.FB
	PlanetFrames []*render.FB
}

type Approach struct {
	T       float64
	Dur     float64
	From    render.Vec3
	FromMat render.Mat3
}

type Hyper struct {
	T       float64
	Phase   int
	Target  int
	Misjump bool
	Dist    float64
}

type DeadInfo struct {
	Reason string
	T      float64
}

type InputState struct {
	Pitch, Roll, Accel float64
	Fire               bool
}

type Game struct {
	State  string
	Cmd    Commander
	Player Player
	World  *World

	Messages []Message
	Time     float64
	Screen   string
	Overlay  string
	Market   []galaxy.Market
	DockSalt uint32

	Autopilot     bool
	AutopilotPath *Approach
	Hyper         *Hyper
	Dead          *DeadInfo

	Shake, BeamFlash, EcmFlash, BombFlash float64
	EncounterTimer                        float64
	TrafficTimer                          float64

	In      InputState
	OnEvent func(string)

	LastDamageReason string
	EscapeUsed       bool
	DockT            float64

	// contract targets rolled when the station is entered, so the briefing and
	// the accepted mission always agree
	OfferTarget    int
	OfferDocSystem int
	OfferDocGalaxy int
	Offers         []Contract
}

var G *Game

func New() *Game {
	g := &Game{State: "title", Cmd: NewCommander()}
	G = g
	return g
}

func (g *Game) Emit(name string) {
	if g.OnEvent != nil {
		g.OnEvent(name)
	}
}

func (g *Game) Message(text, kind string) {
	if kind == "" {
		kind = "info"
	}
	g.Messages = append(g.Messages, Message{Text: text, Life: 4.5, Kind: kind})
	if len(g.Messages) > 4 {
		g.Messages = g.Messages[1:]
	}
}

func randRange(a, b float64) float64 { return a + rand.Float64()*(b-a) }

func randDir() render.Vec3 {
	z := randRange(-1, 1)
	t := randRange(0, 2*math.Pi)
	r := math.Sqrt(1 - z*z)
	return render.V3(math.Cos(t)*r, math.Sin(t)*r, z)
}

func dangerLevel(sys *galaxy.System) float64 {
	d := 0.15 + float64(7-sys.Government)*0.09
	if sys.Tech < 5 {
		d += 0.12
	}
	return d
}

func makeStars() []Star {
	stars := make([]Star, 90)
	for i := range stars {
		b := rand.Float64()
		size := 1
		if b < 0.12 {
			size = 2
		}
		stars[i] = Star{Dir: randDir(), Bright: rand.Float64(), Size: size}
	}
	return stars
}

func makeDust() []Dust {
	d := make([]Dust, 40)
	for i := range d {
		d[i] = Dust{Pos: render.V3(randRange(-1500, 1500), randRange(-1000, 1000), randRange(-3000, -200))}
	}
	return d
}

// EnterSystem builds the in-system scene.
func (g *Game) EnterSystem(systemIndex int, mode string, spawnDist float64, misjump bool) {
	gal := galaxy.Get(g.Cmd.Galaxy)
	sys := gal.Systems[(systemIndex%256+256)%256]
	g.Cmd.System = sys.ID

	rp := 4200 + float64(sys.Seed%2600)
	planetPos := render.V3(0, 0, -rp*4)
	stationPos := render.V3(rp*0.6, rp*0.18, -rp*1.2)

	spawnPos := render.V3(0, 0, 0)
	spawnMat := render.MatIdentity()
	if mode == "launch" {
		sm := render.RotY(randRange(0, 6.28))
		normal := sm.Col(2)
		spawnPos = stationPos.Add(normal.Mul(StationSize + 420))
		spawnMat = render.LookMatrix(normal)
	} else {
		if spawnDist <= 0 {
			spawnDist = 9000
		}
		spawnPos = stationPos.Add(render.V3(0, 0, spawnDist))
		spawnMat = render.LookMatrix(stationPos.Sub(spawnPos).Norm())
	}

	var w *World
	if misjump {
		// witchspace: no planet, no star, no station - only the void and
		// whatever came through with you
		w = &World{Sys: sys, Stars: makeStars(), Dust: makeDust(), Witchspace: true}
		spawnPos = render.V3(0, 0, 0)
		spawnMat = render.MatIdentity()
	} else {
		stMesh := stationMesh(sys)
		w = &World{
			Sys:    sys,
			Planet: Planet{Pos: planetPos, Radius: rp, Seed: sys.Seed, Kind: sys.Kind},
			Sun:    Sun{Pos: render.V3(rp*26, rp*6, -rp*30), Radius: rp * 7.5},
			Stars:  makeStars(),
			Dust:   makeDust(),
			Hermit: stMesh == "hermit",
			Station: &Entity{
				Kind: "station", Pos: stationPos, Mat: render.RotY(randRange(0, 6.28)),
				Radius: StationSize, Energy: 100000, Mesh: models.Get(stMesh),
			},
		}
		// a generation ship now and then, far out and minding its own business
		if sys.Seed%13 == 0 {
			w.Entities = append(w.Entities, &Entity{
				Kind: "genship", Pos: render.V3(-rp*2.5, rp*0.4, -rp*3.4),
				Mat: render.RotY(randRange(0, 6.28)), Radius: 900, Energy: 4000,
				Mesh: models.Get("genship"), Speed: 6,
			})
		}
	}
	g.World = w
	speed := 260.0
	if mode == "launch" {
		speed = 90
	}
	g.Player = Player{
		Pos: spawnPos, Mat: spawnMat, Speed: speed,
		Energy: g.Ship().Energy, ShieldF: MaxShield, ShieldA: MaxShield,
	}
	if mode == "launch" {
		g.Player.Invulnerable = 3
	} else {
		g.Player.Invulnerable = 1.5
	}
	g.State = "flight"
	g.Screen = "flight"
	g.Autopilot = false
	g.AutopilotPath = nil
	g.Messages = nil
	g.Hyper = nil
	if mode == "launch" {
		g.EncounterTimer = 25
		g.Message("LAUNCHED FROM "+sys.Name+" STATION", "info")
	} else {
		g.EncounterTimer = 12
		g.Message("ARRIVED: "+sys.Name, "good")
		if misjump {
			g.Message("HYPERSPACE MISJUMP - WITCHSPACE", "bad")
			g.spawnThargoids()
			g.spawnThargoids()
		} else {
			if rand.Float64() < 0.55 {
				g.SpawnTraffic(1 + rand.Intn(2))
			}
			if rand.Float64() < dangerLevel(sys)*0.6 {
				g.SpawnPirates(1 + rand.Intn(2))
			}
		}
	}
	if !misjump {
		g.missionArrive()
		g.campaignArrive()
	}
}

// stationMesh picks the hull for this system: every fourth chart entry is a
// dodecahedron, the rest are Coriolis.
func stationMesh(sys *galaxy.System) string {
	// rock hermits are a rarity, not a station type: a handful per galaxy, and
	// never the starting system (seed 0)
	if sys.Seed == 0 {
		return "station" // the starting system always gets a classic Coriolis
	}
	if sys.Seed%37 == 0 {
		return "hermit"
	}
	if sys.Seed%4 == 0 {
		return "dodeca"
	}
	return "station"
}

func (g *Game) baseShip(typeKey string, pos render.Vec3, behaviour string) *Entity {
	def := Ships[typeKey]
	e := &Entity{
		Kind: "ship", Type: typeKey, Mesh: models.Get(def.Mesh),
		Pos: pos, Mat: render.LookMatrix(randDir()),
		Speed: def.Speed * 0.4, MaxSpeed: def.Speed,
		Energy: def.Energy, MaxEnergy: def.Energy,
		Shield: def.Shield, MaxShield: def.Shield,
		Radius: def.Radius, Bounty: def.Bounty, CargoType: def.Cargo,
		Behaviour: behaviour, Laser: def.Laser, FireRate: def.FireRate, Turn: def.Turn,
		Cooldown: randRange(0.2, 1.5),
		Hostile:  behaviour == "pirate" || behaviour == "thargoid" || behaviour == "police",
	}
	switch typeKey {
	case "viper", "asp", "ferdelance", "thargoid", "cobra3", "python":
		e.CanMissile = true
	}
	return e
}

func (g *Game) spawnAtDistance(typeKey string, dist float64, behaviour string) *Entity {
	d := randDir()
	d[2] = -math.Abs(d[2])*0.8 - 0.2
	d = d.Norm()
	pos := g.Player.Pos.Add(d.Mul(dist))
	return g.baseShip(typeKey, pos, behaviour)
}

func (g *Game) SpawnPirates(n int) {
	types := []string{"sidewinder", "krait", "mamba", "adder", "cobra3", "moray", "ferdelance", "asp", "gecko", "cobra1"}
	for i := 0; i < n; i++ {
		t := types[rand.Intn(len(types))]
		if g.World.Sys.Government <= 1 && rand.Float64() < 0.3 {
			t = "ferdelance"
		}
		g.World.Entities = append(g.World.Entities, g.spawnAtDistance(t, randRange(2500, 5000), "pirate"))
	}
	g.Message("HOSTILE CONTACT", "bad")
}

func (g *Game) SpawnTraffic(n int) {
	types := []string{"adder", "python", "cobra3", "shuttle", "moray", "worm", "boa", "anaconda", "transporter"}
	for i := 0; i < n; i++ {
		t := types[rand.Intn(len(types))]
		e := g.spawnAtDistance(t, randRange(2200, 6500), "trader")
		e.Hostile = false
		g.World.Entities = append(g.World.Entities, e)
	}
}

func (g *Game) spawnPolice(n int) {
	for i := 0; i < n; i++ {
		g.World.Entities = append(g.World.Entities, g.spawnAtDistance("viper", randRange(4000, 7000), "police"))
	}
	g.Message("POLICE VIPER ON INTERCEPT", "bad")
}

func (g *Game) spawnThargoids() {
	for i := 0; i < 2; i++ {
		g.World.Entities = append(g.World.Entities, g.spawnAtDistance("thargoid", randRange(3000, 6000), "thargoid"))
	}
	// Thargoids never travel alone: each mothership carries its drones
	for i := 0; i < 3; i++ {
		g.World.Entities = append(g.World.Entities, g.spawnAtDistance("thargon", randRange(2600, 5200), "thargoid"))
	}
	g.Message("SWARM ATTACK", "bad")
}

func (g *Game) SpawnCargo(pos render.Vec3, cargoType int) {
	g.World.Entities = append(g.World.Entities, &Entity{
		Kind: "canister", Pos: pos, Mat: render.LookMatrix(randDir()),
		Speed: randRange(10, 40), Radius: 8, CargoType: cargoType,
		Mesh: models.Get("canister"),
	})
}

// SpawnEscapePod launches a pilot's capsule; scooping it yields slaves.
func (g *Game) SpawnEscapePod(pos render.Vec3) {
	g.World.Entities = append(g.World.Entities, &Entity{
		Kind: "capsule", Pos: pos.Add(randDir().Mul(25)),
		Mat: render.LookMatrix(randDir()), Speed: randRange(15, 45),
		Radius: 9, CargoType: 3, Mesh: models.Get("capsule"), Life: 180,
	})
	g.Message("ESCAPE CAPSULE JETTISONED", "info")
}

func (g *Game) SpawnAsteroid(dist float64) {
	d := randDir()
	d[2] = -math.Abs(d[2]) - 0.1
	pos := g.Player.Pos.Add(d.Norm().Mul(dist))
	name := "asteroid"
	if rand.Float64() < 0.5 {
		name = "asteroid2"
	}
	g.World.Entities = append(g.World.Entities, &Entity{
		Kind: "asteroid", Pos: pos, Mat: render.LookMatrix(randDir()),
		Speed: randRange(20, 60), Radius: 20, Energy: 60, MaxEnergy: 60,
		Spin: render.V3(randRange(-0.6, 0.6), randRange(-0.6, 0.6), randRange(-0.6, 0.6)),
		Mesh: models.Get(name),
	})
}

func (g *Game) addExplosion(pos render.Vec3, radius float64) {
	g.World.Entities = append(g.World.Entities, &Entity{
		Kind: "explosion", Pos: pos, Radius: radius, Life: 0.6 + radius/90,
	})
}

func (g *Game) countBehaviour(b string) int {
	n := 0
	for _, e := range g.World.Entities {
		if e.Kind == "ship" && e.Behaviour == b {
			n++
		}
	}
	return n
}

// ------------------------------------------------------------- combat

func (g *Game) FireLaser() {
	if g.State != "flight" {
		return
	}
	p := &g.Player
	if p.LaserTemp > 0.98 {
		return
	}
	cmd := &g.Cmd
	dmg := LaserDamage[cmd.Laser]
	rng := LaserRange[cmd.Laser]
	p.LaserTemp = math.Min(1, p.LaserTemp+LaserHeat[cmd.Laser])
	fwd := p.Mat.Col(2).Mul(-1)
	var best *Entity
	bestDepth := math.Inf(1)
	for _, e := range g.World.Entities {
		if e.Kind == "planet" || e.Kind == "sun" || e.Kind == "explosion" {
			continue
		}
		rel := e.Pos.Sub(p.Pos)
		depth := rel.Dot(fwd)
		if depth < 30 || depth > rng {
			continue
		}
		lateral := rel.Sub(fwd.Mul(depth)).Len()
		angRad := e.Radius / depth
		spread := math.Max(angRad*1.25, 0.012)
		if lateral < depth*spread && depth < bestDepth {
			best = e
			bestDepth = depth
		}
	}
	if best == nil {
		return
	}
	g.HitEntity(best, dmg, "laser")
}

func (g *Game) HitEntity(e *Entity, dmg float64, source string) {
	if e.Kind == "missile" {
		g.removeEntity(e)
		return
	}
	if e.Kind == "station" {
		return
	}
	if e.Shield > 0 {
		e.Shield -= dmg
		if e.Shield < 0 {
			e.Energy += e.Shield
			e.Shield = 0
		}
	} else {
		e.Energy -= dmg
	}
	if e.Kind == "canister" {
		e.Energy = 0
	}
	if e.Energy <= 0 {
		g.DestroyEntity(e, source)
	} else if e.Kind == "ship" && e.Behaviour == "trader" {
		e.Behaviour = "flee"
		e.Hostile = false
	}
}

func (g *Game) removeEntity(e *Entity) {
	for i, x := range g.World.Entities {
		if x == e {
			g.World.Entities = append(g.World.Entities[:i], g.World.Entities[i+1:]...)
			return
		}
	}
}

func (g *Game) DestroyEntity(e *Entity, source string) {
	cmd := &g.Cmd
	switch e.Kind {
	case "ship":
		def := Ships[e.Type]
		if e.CargoType > 0 && rand.Float64() < 0.6 {
			n := 1 + rand.Intn(maxInt(1, minInt(3, e.CargoType)))
			for k := 0; k < n; k++ {
				g.SpawnCargo(e.Pos.Add(randDir().Mul(20)), rand.Intn(12))
			}
		}
		if e.Behaviour == "pirate" || e.Behaviour == "thargoid" {
			bounty := def.Bounty
			if e.Mission {
				bounty = e.Bounty // contract ships are paid through the mission
			}
			cmd.Cash += float64(bounty)
			cmd.Kills++
			g.countContracts()
			if def.Bounty > 0 {
				g.Message("BOUNTY "+itoa(def.Bounty)+" CR", "good")
			}
		} else if e.Behaviour == "police" {
			cmd.Legal += 50
			g.Message("POLICE SHIP DESTROYED  +50 CR BOUNTY ON YOU", "bad")
		} else {
			cmd.Legal += 20
			g.Message("TRADER DESTROYED  LEGAL RECORD WORSENED", "bad")
		}
		g.missionKill(e)
		g.campaignKill(e.CampaignTag)
		if (e.Behaviour == "trader" || e.Behaviour == "pirate") && rand.Float64() < 0.25 {
			g.SpawnEscapePod(e.Pos)
		}
		// wreckage: alloy plates are worth scooping
		for i := 0; i < 1+rand.Intn(2); i++ {
			if rand.Float64() < 0.6 {
				g.SpawnPlate(e.Pos)
			}
		}
		g.addExplosion(e.Pos, e.Radius)
		g.Emit("explode")
	case "asteroid":
		g.addExplosion(e.Pos, e.Radius)
		g.Emit("explode")
		if cmd.Equip["mining"] || e.Mission {
			// a mining laser cuts the rock into pieces you can scoop
			g.breakRock(e, "boulder", 3)
		} else if rand.Float64() < 0.4 {
			g.SpawnCargo(e.Pos.Add(randDir().Mul(30)), 12)
		}
	case "boulder":
		g.addExplosion(e.Pos, e.Radius*0.8)
		g.Emit("hit")
		g.breakRock(e, "splinter", 3)
	case "canister", "plate":
		g.addExplosion(e.Pos, 10)
	case "genship":
		g.addExplosion(e.Pos, 400)
		g.Emit("explode")
		for i := 0; i < 6; i++ {
			g.SpawnCargo(e.Pos.Add(randDir().Mul(200)), rand.Intn(12))
		}
	}
	g.removeEntity(e)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [12]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func (g *Game) FireMissile() {
	cmd := &g.Cmd
	if cmd.Missiles <= 0 {
		g.Message("NO MISSILES", "bad")
		return
	}
	target := g.pickMissileTarget()
	if target == nil {
		g.Message("NO TARGET LOCKED", "bad")
		return
	}
	cmd.Missiles--
	g.World.Entities = append(g.World.Entities, &Entity{
		Kind: "missile", Pos: g.Player.Pos, Mat: g.Player.Mat,
		Speed: 300, MaxSpeed: 1100, Radius: 6, Life: 9, Energy: 10,
		TargetEnt: target, Mesh: models.Get("missile"),
	})
	g.Message("MISSILE AWAY", "good")
	g.Emit("missile")
}

func (g *Game) pickMissileTarget() *Entity {
	p := &g.Player
	fwd := p.Mat.Col(2).Mul(-1)
	var best *Entity
	bestScore := -1.0
	for _, e := range g.World.Entities {
		if e.Kind != "ship" {
			continue
		}
		rel := e.Pos.Sub(p.Pos)
		depth := rel.Dot(fwd)
		if depth < 0 || depth > 30000 {
			continue
		}
		ang := rel.Sub(fwd.Mul(depth)).Len() / math.Max(depth, 1)
		score := math.Cos(math.Min(math.Pi/2, ang))
		if e.Hostile {
			score *= 1.5
		}
		if score > bestScore {
			bestScore = score
			best = e
		}
	}
	return best
}

func (g *Game) EnemyMissile(from *Entity) {
	back := from.Mat.Col(2)
	g.World.Entities = append(g.World.Entities, &Entity{
		Kind: "missile", Pos: from.Pos.Add(back.Mul(from.Radius)),
		Mat: from.Mat, Speed: 260, MaxSpeed: 900, Radius: 6, Life: 10, Energy: 6,
		TargetPlayer: true, Mesh: models.Get("missile"),
	})
	g.Message("INCOMING MISSILE", "bad")
	g.Emit("warning")
}

func (g *Game) FireECM() {
	if !g.Cmd.Equip["ecm"] {
		g.Message("NO E.C.M. FITTED", "bad")
		return
	}
	killed := 0
	keep := g.World.Entities[:0]
	for _, e := range g.World.Entities {
		if e.Kind == "missile" {
			killed++
			continue
		}
		keep = append(keep, e)
	}
	g.World.Entities = keep
	g.Player.Energy = math.Max(0, g.Player.Energy-4)
	if killed > 0 {
		g.Message("E.C.M. ACTIVE - MISSILES DESTROYED", "info")
	} else {
		g.Message("E.C.M. ACTIVE", "info")
	}
	g.EcmFlash = 0.4
	g.Emit("ecm")
}

func (g *Game) FireBomb() {
	cmd := &g.Cmd
	if cmd.Bombs <= 0 {
		g.Message("NO ENERGY BOMB", "bad")
		return
	}
	cmd.Bombs--
	killed := 0
	var keep []*Entity
	for _, e := range g.World.Entities {
		if e.Kind == "ship" && e.Type != "thargoid" {
			g.addExplosion(e.Pos, e.Radius)
			killed++
			continue
		}
		keep = append(keep, e)
	}
	g.World.Entities = keep
	g.BombFlash = 1.2
	g.Message("ENERGY BOMB DETONATED  ("+itoa(killed)+" DESTROYED)", "good")
	g.Emit("bomb")
}

func (g *Game) DamagePlayer(dmg float64, fromFront bool) {
	p := &g.Player
	if p.Invulnerable > 0 {
		return
	}
	var sh float64
	if fromFront {
		sh = p.ShieldF
	} else {
		sh = p.ShieldA
	}
	if sh > 0 {
		sh -= dmg
		if sh < 0 {
			p.Energy += sh
			sh = 0
		}
		if fromFront {
			p.ShieldF = sh
		} else {
			p.ShieldA = sh
		}
	} else {
		p.Energy -= dmg
	}
	if p.Energy < 0 {
		p.Energy = 0
	}
	if p.Energy <= 0 {
		g.PlayerDeath()
	}
}

func (g *Game) PlayerDeath() {
	cmd := &g.Cmd
	if cmd.Equip["escape"] && !g.EscapeUsed {
		g.EscapeUsed = true
		cmd.Equip["escape"] = false
		cmd.Cash = math.Max(0, math.Round(cmd.Cash*0.5))
		cmd.Cargo = map[int]int{}
		cmd.Missiles = 0
		g.Message("ESCAPE CAPSULE LAUNCHED - RESCUED", "bad")
		// the rescue rebuilds the ship: without this the commander sat docked
		// with the negative energy of the hit that killed them
		g.Player.Energy = g.Ship().Energy
		g.Player.ShieldF = MaxShield
		g.Player.ShieldA = MaxShield
		g.Player.Cabin = 0
		g.Player.Speed = 0
		g.State = "docked"
		g.Screen = "launch"
		g.dockedAt()
		return
	}
	g.State = "dead"
	g.Dead = &DeadInfo{Reason: g.LastDamageReason}
	g.Emit("explode")
	cmd.Deaths++
}

func (g *Game) Docked() { g.dockedAt() }

func (g *Game) dockedAt() {
	g.Screen = "launch"
	g.Market = galaxy.MakeMarket(g.World.Sys, g.DockSalt)
	if g.World != nil && g.World.Hermit {
		// a rock hermit has no competition and knows it
		for i := range g.Market {
			g.Market[i].Price = int(float64(g.Market[i].Price) * 1.35)
			g.Market[i].Qty = int(float64(g.Market[i].Qty) * 0.6)
		}
	}
	g.DockSalt++
	g.Autopilot = false
	g.AutopilotPath = nil
	g.Player.Speed = 0
	if g.World.Nebula == nil {
		g.World.Nebula = models.MakeNebula(g.World.Sys.Seed, g.Cmd.Galaxy)
	}
	g.rollOffers()
	g.expireContracts()
	g.deliverContracts()
	g.rollContracts()
}

// TryDock checks alignment with the rotating slot.
func (g *Game) TryDock() bool {
	w := g.World
	if w == nil {
		return false
	}
	st := w.Station
	if st == nil {
		return false
	}
	p := &g.Player
	local := st.Mat.ApplyT(p.Pos.Sub(st.Pos))
	slotW := SlotW * StationSize
	slotH := SlotH * StationSize
	// dock only when the hull has reached the mouth; the old window started 90
	// units short of the face, so the station screen appeared while the ship
	// was still outside watching itself approach
	if local[2] > StationSize-MouthDepth && local[2] < StationSize+8 {
		if math.Abs(local[0]) < slotW+10 && math.Abs(local[1]) < slotH+10 {
			normal := st.Mat.Col(2)
			fwd := p.Mat.Col(2).Mul(-1)
			align := fwd.Dot(normal.Mul(-1))
			if align > 0.72 && p.Speed < g.Ship().Speed*0.85 {
				g.Dock()
				return true
			}
		}
	}
	return false
}

// Dock starts the final glide into the shaft. The station screen comes up once
// the ship has actually travelled inside.
func (g *Game) Dock() {
	if g.State == "docking" {
		return
	}
	g.State = "docking"
	g.DockT = 0
	g.Player.Speed = 45
	g.Player.Invulnerable = 3
}

// finishDock is the end of the glide.
func (g *Game) finishDock() {
	g.State = "docked"
	g.Screen = "launch"
	g.Player.Speed = 0
	g.missionDock()
	g.campaignDock()
	g.dockedAt()
	g.Message("DOCKED", "good")
	g.Emit("dock")
}

func (g *Game) updateDocking(dt float64) {
	st := g.World.Station
	if st == nil {
		g.finishDock()
		return
	}
	g.DockT += dt
	normal := st.Mat.Col(2)
	fwd := normal.Mul(-1)
	g.Player.Mat = render.LookMatrix(fwd)
	g.Player.Pos = g.Player.Pos.Add(fwd.Mul(120 * dt))
	st.Mat = st.Mat.Mul(render.RotY(0.05 * dt))
	if g.DockT > 0.9 {
		g.finishDock()
	}
}

func (g *Game) Launch() {
	g.EscapeUsed = false
	g.EnterSystem(g.Cmd.System, "launch", 0, false)
}

func (g *Game) Hyperspace(target int) bool {
	cmd := &g.Cmd
	gal := galaxy.Get(cmd.Galaxy)
	from := gal.Systems[cmd.System]
	to := gal.Systems[((target%256)+256)%256]
	d := galaxy.Distance(from, to)
	if to.ID == from.ID {
		g.Message("ALREADY IN SYSTEM", "bad")
		return false
	}
	if d > 7 {
		g.Message("TARGET OUT OF RANGE (MAX 7.0 LY)", "bad")
		return false
	}
	if cmd.Fuel < d {
		g.Message("INSUFFICIENT FUEL", "bad")
		return false
	}
	cmd.Fuel = math.Round((cmd.Fuel-d)*10) / 10
	cmd.Day++
	cmd.TargetSystem = to.ID
	g.State = "hyper"
	g.Hyper = &Hyper{Target: to.ID, Misjump: rand.Float64() < 0.03, Dist: d}
	g.Message("HYPERSPACE TO "+to.Name, "info")
	g.Emit("hyper")
	return true
}

func (g *Game) JumpInSystem() {
	if g.State != "flight" {
		return
	}
	if g.Classic() {
		g.Message("NOT AVAILABLE IN CLASSIC MODE", "bad")
		return
	}
	if g.Cmd.Fuel < 0.1 {
		g.Message("NOT ENOUGH FUEL FOR JUMP", "bad")
		return
	}
	st := g.World.Station
	dist := st.Pos.Sub(g.Player.Pos).Len()
	if dist < 2500 {
		g.Message("ALREADY AT THE STATION", "bad")
		return
	}
	g.Cmd.Fuel = math.Round((g.Cmd.Fuel-0.1)*10) / 10
	dir := st.Pos.Sub(g.Player.Pos).Norm()
	g.Player.Pos = st.Pos.Sub(dir.Mul(3000))
	g.Player.Mat = render.LookMatrix(dir)
	g.Player.Speed = 200
	g.Message("IN-SYSTEM JUMP", "info")
	if rand.Float64() < dangerLevel(g.World.Sys)*0.5 {
		g.SpawnPirates(1 + rand.Intn(2))
	}
}

// -------------------------------------------------------------- docking aid

func (g *Game) autopilotState() (aim render.Vec3, inLine bool, dist float64, local render.Vec3, normal render.Vec3, gate render.Vec3) {
	st := g.World.Station
	p := &g.Player
	normal = st.Mat.Col(2)
	rel := p.Pos.Sub(st.Pos)
	local = st.Mat.ApplyT(rel)
	dist = rel.Len()
	lateral := math.Hypot(local[0], local[1])
	inLine = local[2] > StationSize+60 && lateral < 55
	gate = st.Pos.Add(normal.Mul(StationSize + 1400))
	if inLine {
		aim = normal.Mul(-1)
	} else {
		aim = gate.Sub(p.Pos).Norm()
	}
	return
}

func (g *Game) ToggleAutopilot() {
	if !g.Cmd.Equip["dockcomp"] {
		g.Message("DOCKING COMPUTER NOT FITTED", "bad")
		return
	}
	if g.Autopilot {
		g.Autopilot = false
		g.AutopilotPath = nil
		g.Message("MANUAL CONTROL", "info")
		return
	}
	aim, _, dist, local, normal, gate := g.autopilotState()
	_ = aim
	if dist > 3000 || local[2] < 0 || math.Hypot(local[0], local[1]) > 300 {
		g.Player.Pos = gate
		g.Player.Speed = 70
		g.Player.Mat = render.LookMatrix(normal.Mul(-1))
	}
	g.Autopilot = true
	g.AutopilotPath = &Approach{T: 0, Dur: 6.0, From: g.Player.Pos, FromMat: g.Player.Mat}
	g.Message("DOCKING COMPUTER ON", "info")
}

func smoothstep(t float64) float64 { return t * t * (3 - 2*t) }

func (g *Game) runAutopilot(dt float64) {
	p := &g.Player
	st := g.World.Station
	path := g.AutopilotPath
	path.T += dt / path.Dur
	k := 1.0
	if path.T < 1 {
		k = smoothstep(path.T)
	}
	normal := st.Mat.Col(2)
	end := st.Pos.Add(normal.Mul(StationSize - MouthDepth*0.5))
	old := p.Pos
	p.Pos = render.V3(
		path.From[0]+(end[0]-path.From[0])*k,
		path.From[1]+(end[1]-path.From[1])*k,
		path.From[2]+(end[2]-path.From[2])*k,
	)
	p.Speed = old.Sub(p.Pos).Len() / math.Max(dt, 0.0001)
	target := render.LookMatrix(normal.Mul(-1))
	var m render.Mat3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m[i][j] = path.FromMat[i][j]*(1-k) + target[i][j]*k
		}
	}
	p.Mat = m.Orthonormalize()
	if path.T >= 1 {
		g.Autopilot = false
		g.AutopilotPath = nil
		g.Dock()
	}
}

// meshFor exposes the mesh table to tests without importing models there.
func meshFor(name string) *render.Mesh { return models.Get(name) }
