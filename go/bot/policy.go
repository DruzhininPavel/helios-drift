package bot

import (
	"math"
	"math/rand"
	"sort"

	"heliosdrift/galaxy"
	"heliosdrift/game"
	"heliosdrift/render"
)

// Docking phases: line up on the axis of the slot, face the mouth, creep in.
const (
	phaseTravel = iota
	phaseStandoff
	phaseAlign
	phaseEnter
)

// Trader is the first working policy: it flies, docks, trades, takes work and
// jumps. It is deliberately paced like a person — it decides a few times a
// second, not every frame — and it only uses the controls a player has.
type Trader struct {
	hold     int // frames until the next decision
	in       game.InputState
	rng      *rand.Rand
	pilot    *Pilot
	nextJump int
	phase    int
	stuckFor int
	// ManualDock forces the hand-flown approach even with a computer fitted.
	ManualDock bool
	// Ace makes the pilot a dead shot: the bot still has to fly the attack,
	// but it fires the instant the nose is on the hull. Used to answer "could
	// a player win this fight at all".
	Ace bool
	// questGood is cargo the campaign stage needs, which must not be sold.
	questGood int
	// TradeOnly keeps the bot to buying and selling: the hull benchmark uses
	// it so story payouts do not mask what the cargo hold is worth.
	TradeOnly bool
	// cost is the average price paid for each good in the hold, so the bot
	// does not sell at a loss just to be rid of cargo.
	cost map[int]float64
	// dest is the system the bot is travelling to, -1 when it has no plan.
	dest int
}

// NewTrader builds the policy for a seed.
func NewTrader(seed int64) *Trader {
	return &Trader{rng: rand.New(rand.NewSource(seed ^ 0x5EED)), pilot: DefaultPilot(),
		cost: map[int]float64{}, dest: -1, questGood: -1}
}

// Name is used in reports.
func (t *Trader) Name() string { return "trader" }

// Act is called once per frame; the decision is refreshed twelve times a
// second so the stick is held the way a hand holds it.
func (t *Trader) Act(g *game.Game, frame int) (game.InputState, []Event) {
	var ev []Event
	if t.hold > 0 {
		t.hold--
		return t.in, nil
	}
	t.hold = 4 + t.rng.Intn(4)
	// Classic has no yaw at all: the pilot must bank and pull there
	t.pilot.UseYaw = !g.Classic()

	switch g.State {
	case "docked":
		ev = append(ev, t.docked(g)...)
		t.in = game.InputState{}
	case "flight":
		var in game.InputState
		ev, in = t.flying(g, frame)
		t.in = in
	default:
		t.in = game.InputState{}
	}
	return t.in, ev
}

// docked does the station round: sell, buy, refuel, take work, launch.
func (t *Trader) docked(g *game.Game) []Event {
	var ev []Event
	c := &g.Cmd

	// sell first, then spend: the hold is emptied at every port. The order is
	// sorted because map iteration order is random and would make the recorded
	// script differ between two runs of the same seed.
	goods := make([]int, 0, len(c.Cargo))
	for g := range c.Cargo {
		goods = append(goods, g)
	}
	sort.Ints(goods)
	for _, good := range goods {
		qty := c.Cargo[good]
		if qty <= 0 {
			continue
		}
		if good == t.questGood && c.Campaign.Active {
			continue // the stage needs this cargo delivered, not sold
		}
		// sell only if the station pays more than was paid for it
		if good < len(g.Market) && float64(g.Market[good].Sell) > t.cost[good]*1.02 {
			ev = append(ev, Event{Kind: "sell", A: good, B: qty})
			delete(t.cost, good)
		}
	}
	if c.Fuel < 3.5 && c.Cash >= 20 {
		ev = append(ev, Event{Kind: "refuel"})
	}
	// the 1984 contracts only exist in Classic mode
	if kind := g.MissionOffer(); kind != "" && !t.TradeOnly {
		switch kind {
		case "constrictor":
			ev = append(ev, Event{Kind: "mission", A: 0})
		case "documents":
			if c.Equip["galdrive"] {
				ev = append(ev, Event{Kind: "mission", A: 1})
			}
		}
	}
	if g.CampaignOffer() && !t.TradeOnly {
		ev = append(ev, Event{Kind: "campaign"})
		if q := g.CampaignQuest(); q.Kind == "deliver" {
			t.questGood = q.Good
		}
	}
	if len(g.Offers) > 0 && len(c.Contracts) < 3 && t.rng.Intn(2) == 0 && !t.TradeOnly {
		ev = append(ev, Event{Kind: "contract", A: t.rng.Intn(len(g.Offers))})
	}
	// buy what is cheap against its galactic base price, keeping a reserve
	if room := c.CargoCap - game.CargoUsed(*c); room > 0 && len(g.Market) > 0 {
		// the cheapest good against its galactic base price, even when it is
		// not a bargain here: the profit is made at the other end of the jump
		best, bestRatio, price := -1, math.MaxFloat64, 0
		for i, m := range g.Market {
			if i >= len(galaxy.Commodities) || m.Qty <= 0 || float64(m.Price) > c.Cash-20 || m.Price <= 0 {
				continue
			}
			ratio := float64(m.Price) / math.Max(1, galaxy.Commodities[i].Base)
			if ratio < bestRatio {
				best, bestRatio, price = i, ratio, m.Price
			}
		}
		if best >= 0 && price > 0 {
			n := room
			if n > 10 {
				n = 10
			}
			if afford := int(c.Cash-20) / price; n > afford {
				n = afford
			}
			if n > 0 {
				ev = append(ev, Event{Kind: "buy", A: best, B: n})
				t.cost[best] = (t.cost[best]*float64(c.Cargo[best]) + float64(price*n)) / float64(c.Cargo[best]+n)
			}
		}
	}
	// spend spare cash on a sturdier ship, keeping a trading float. The
	// benchmark compares hulls, so it must not buy a different one.
	if !t.TradeOnly {
		ev = append(ev, t.upgrade(g)...)
	}

	// the story comes first: fly the stage target, trade on the way
	t.dest = -1
	if g.Cmd.Campaign.Active {
		t.dest = g.Cmd.Campaign.Target
	} else if g.Cmd.Mission.Kind != "" {
		t.dest = g.Cmd.Mission.System
	} else if id, ok := t.reachable(g); ok {
		t.questGood = -1 // no stage in progress: nothing is reserved
		t.dest = id
	}
	t.phase = phaseTravel
	t.stuckFor = 0
	ev = append(ev, Event{Kind: "launch"})
	return ev
}

// flying handles the approach, the docking run and the odd hyperspace jump.
func (t *Trader) flying(g *game.Game, frame int) ([]Event, game.InputState) {
	var ev []Event

	// a hostile on the scope comes first: nothing else matters while it is
	// shooting, and the docking computer has to be switched off to fight
	if foe := t.nearestFoe(g); foe != nil {
		if g.Autopilot {
			ev = append(ev, Event{Kind: "autopilot"})
		}
		in, fire := t.dogfight(g, foe)
		if fire {
			in.Fire = true
		}
		return ev, t.safety(g, in)
	}

	// with a docking computer the pilot lets the machine do the work, but first
	// the trip: jump to the chosen system, then dock
	if g.Cmd.Equip["dockcomp"] && !t.ManualDock {
		if t.dest >= 0 && g.Cmd.System != t.dest {
			if hop, ok := t.nextHop(g, t.dest); ok {
				if g.Hyperspace(hop) {
					ev = append(ev, Event{Kind: "hyper", A: hop})
				}
			} else {
				t.dest = -1 // out of fuel for the trip: dock and refuel
			}
			return ev, game.InputState{}
		}
		if !g.Autopilot {
			ev = append(ev, Event{Kind: "autopilot"})
		}
		return ev, game.InputState{}
	}

	st := g.World.Station
	if st == nil {
		return ev, game.InputState{}
	}
	in := t.safety(g, t.dockRun(g, st))

	// if the approach drags on, jump somewhere else for a change
	if t.phase == phaseTravel {
		t.stuckFor++
		if t.stuckFor > 60*40 {
			t.stuckFor = 0
			if id, ok := t.reachable(g); ok {
				ev = append(ev, Event{Kind: "hyper", A: id})
			}
		}
	} else {
		t.stuckFor = 0
	}
	return ev, in
}

// dockRun flies the ship onto the axis of the slot and in through the mouth.
//
// The station turns, so the mouth sweeps along the hull at roughly fifteen
// units a second. Standing still and lining up is hopeless: the slot is gone
// by the time the nose arrives. The only approach that works is pursuit —
// always aim at a point on the live axis just outside the mouth, and let the
// ship curve onto the axis as it closes. Speed is limited by how far off the
// axis the ship still is, so it arrives slowly enough to fit the slot.
func (t *Trader) dockRun(g *game.Game, st *game.Entity) game.InputState {
	p := &g.Player
	normal := st.Mat.Col(2)
	local := st.Mat.ApplyT(p.Pos.Sub(st.Pos))
	lateral := math.Hypot(local[0], local[1])
	max := g.Ship().Speed

	// where on the axis to aim: just outside the mouth, never further out than
	// the ship already is, and well clear of the hull while it is off the axis
	along := local[2]
	if along > game.StationSize+80 {
		along = game.StationSize + 80
	}
	if lateral > 100 {
		along = game.StationSize + 700
	}
	if along < game.StationSize-30 {
		along = game.StationSize - 30
	}
	aimPoint := st.Pos.Add(normal.Mul(along))

	// slower the further off the axis the ship still is
	cap := max * clamp(lateral/220, 0.10, 0.35)
	if lateral < 45 {
		cap = max * 0.12
	}

	switch {
	case lateral < 45 && local[2] < game.StationSize+250:
		t.phase = phaseEnter
	case lateral < 100:
		t.phase = phaseAlign
	default:
		t.phase = phaseTravel
	}

	to := aimPoint.Sub(p.Pos)
	if to.Len() < 1 {
		return game.InputState{}
	}
	want := to.Norm()
	in := t.pilot.Aim(g, want)
	in.Accel = holdSpeed(g, want, cap)
	return in
}

// holdSpeed keeps the ship at roughly the given speed, and refuses to push the
// throttle while the nose is still off the carrot.
func holdSpeed(g *game.Game, want render.Vec3, cap float64) float64 {
	fwd := g.Player.Mat.Col(2).Mul(-1)
	aim := fwd.Dot(want)
	switch {
	case g.Player.Speed > cap*1.15:
		return -1
	case aim < 0.90:
		return -1
	case g.Player.Speed < cap*0.85:
		return 1
	}
	return 0
}

// reachable picks a system inside jump range. It prefers a place with the
// opposite kind of economy, which is where the goods in the hold are worth
// most: agricultural worlds pay for machinery and industrial worlds pay for
// food, the same rule a trader learns on the first run.
func (t *Trader) reachable(g *game.Game) (int, bool) {
	gal := galaxy.Get(g.Cmd.Galaxy)
	from := gal.Systems[g.Cmd.System]
	best, bestScore := -1, -1.0
	for _, s := range gal.Systems {
		d := galaxy.Distance(from, s)
		if d <= 0.5 || d > g.Cmd.Fuel || s.ID == g.Cmd.System {
			continue
		}
		score := math.Abs(float64(s.Economy)-float64(from.Economy)) + t.rng.Float64()*2
		if score > bestScore {
			best, bestScore = s.ID, score
		}
	}
	if best < 0 {
		return 0, false
	}
	return best, true
}

// PhaseName reports the docking phase, so runs can be read from a log.
func (t *Trader) PhaseName() string {
	switch t.phase {
	case phaseStandoff:
		return "standoff"
	case phaseAlign:
		return "align"
	case phaseEnter:
		return "enter"
	}
	return "travel"
}

// nearestFoe finds the closest hostile within engagement range.
func (t *Trader) nearestFoe(g *game.Game) *game.Entity {
	var best *game.Entity
	bestD := 1600.0
	// near the station the bot has business to finish: only fight what is
	// actually in the way, otherwise it never docks
	if st := g.World.Station; st != nil && g.Player.Pos.Sub(st.Pos).Len() < 2500 {
		bestD = 900
	}
	for _, e := range g.World.Entities {
		if e.Energy <= 0 || e.Kind == "station" || e.Kind == "planet" {
			continue
		}
		if e.CampaignTag != "" {
			return e // the stage target outranks everything else
		}
		if !e.Hostile && e.Behaviour != "pirate" && e.Behaviour != "thargoid" && e.Behaviour != "police" {
			continue
		}
		if e.Behaviour == "police" && g.Cmd.Legal >= 100 {
			continue // clean record: the law is not the enemy
		}
		if st := g.World.Station; st != nil && e.Pos.Sub(st.Pos).Len() < game.StationSize*3 {
			continue // do not chase anything into the station's lap
		}
		if d := e.Pos.Sub(g.Player.Pos).Len(); d < bestD {
			best, bestD = e, d
		}
	}
	return best
}

// dogfight points at the target and holds a firing distance: close fast, then
// match speed so the nose stays on the hull long enough to shoot.
func (t *Trader) dogfight(g *game.Game, foe *game.Entity) (game.InputState, bool) {
	to := foe.Pos.Sub(g.Player.Pos)
	dist := to.Len()

	// shields going down: run for it, a live trader docks and trades again
	if g.Player.ShieldF < 25 && g.Player.ShieldA < 25 {
		away := to.Mul(-1).Norm()
		in := t.pilot.Aim(g, away)
		in.Accel = 1
		return in, false
	}

	want := to.Norm()
	in := t.pilot.Aim(g, want)
	fwd := g.Player.Mat.Col(2).Mul(-1)
	align := fwd.Dot(want)

	max := g.Ship().Speed
	switch {
	case dist > 2500:
		in.Accel = 1
	case dist > 1200:
		if g.Player.Speed > max*0.6 {
			in.Accel = -1
		} else {
			in.Accel = 1
		}
	case dist > 500:
		if g.Player.Speed > math.Min(190, max*0.4) {
			in.Accel = -1
		}
	case dist < 220:
		in.Accel = -1
	default:
		if g.Player.Speed > math.Min(95, max*0.22) {
			in.Accel = -1
		} else {
			in.Accel = 0
		}
	}
	if align < 0.9 {
		in.Accel = -1
	}
	// never chase a target through the station: keep a low speed near the hull
	if st := g.World.Station; st != nil {
		if d := g.Player.Pos.Sub(st.Pos).Len(); d < 1200 && g.Player.Speed > max*0.25 {
			in.Accel = -1
		}
	}

	if t.Ace {
		angle := math.Acos(clamp(align, -1, 1))
		onTarget := angle*dist < foe.Radius*1.6 || align > 0.995
		return in, onTarget && dist < 4000
	}
	return in, align > 0.99 && dist < 3200
}

// safety is the rule that keeps the bot alive: close to the station and not
// lined up with the slot means point away and brake, whatever else it wanted
// to do. Without it the bot happily chases a pirate into the hull.
func (t *Trader) safety(g *game.Game, in game.InputState) game.InputState {
	st := g.World.Station
	if st == nil {
		return in
	}
	to := st.Pos.Sub(g.Player.Pos)
	d := to.Len()
	if d > game.StationSize*1.9 {
		return in
	}
	local := st.Mat.ApplyT(g.Player.Pos.Sub(st.Pos))
	if math.Hypot(local[0], local[1]) < 45 {
		return in // lined up with the mouth: let it through
	}
	want := g.Player.Pos.Sub(st.Pos)
	if want.Len() < 1 {
		return in
	}
	out := t.pilot.Aim(g, want.Norm())
	out.Accel = -1
	out.Fire = false
	return out
}

// upgrade buys what keeps a trader alive, cheapest first, and only with money
// that is not needed for cargo.
func (t *Trader) upgrade(g *game.Game) []Event {
	var ev []Event
	c := &g.Cmd
	float := 300.0
	buy := func(id string, price float64) bool {
		if c.Equip[id] || c.Cash < price+float {
			return false
		}
		ev = append(ev, Event{Kind: "equip", A: equipIndex(id)})
		return true
	}
	// a beam laser first: the pulse one barely scratches a pirate
	if c.Laser < 1 && !c.Equip["laserbeam"] {
		buy("laserbeam", 1000)
	}
	if c.Missiles < 3 {
		if c.Cash > 30+float {
			ev = append(ev, Event{Kind: "equip", A: equipIndex("missile")})
		}
	}
	buy("galdrive", 5000)
	buy("extraenergy", 1500)
	buy("escape", 1000)
	buy("ecm", 600)
	buy("fuelscoop", 525)
	buy("largebay", 400)
	return ev
}

// equipIndex maps an equipment id to its position in the market list.
func equipIndex(id string) int {
	for i := range galaxy.EquipmentList {
		if galaxy.EquipmentList[i].ID == id {
			return i
		}
	}
	return -1
}

// nextHop picks the jump that gets the ship closest to the destination, so a
// target several jumps away is reached hop by hop.
func (t *Trader) nextHop(g *game.Game, dest int) (int, bool) {
	gal := galaxy.Get(g.Cmd.Galaxy)
	from := gal.Systems[g.Cmd.System]
	want := gal.Systems[((dest%256)+256)%256]
	best, bestScore := -1, math.MaxFloat64
	for _, s := range gal.Systems {
		d := galaxy.Distance(from, s)
		if d <= 0.1 || d > g.Cmd.Fuel || s.ID == g.Cmd.System {
			continue
		}
		score := galaxy.Distance(s, want) + d*0.05
		if score < bestScore {
			best, bestScore = s.ID, score
		}
	}
	if best < 0 {
		return 0, false
	}
	return best, true
}
