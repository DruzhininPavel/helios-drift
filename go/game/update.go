package game

import (
	"math"

	"heliosdrift/galaxy"
	"heliosdrift/render"
)

// Update advances the simulation by dt seconds.
// warnYaw tells a Classic pilot once per session why the new axis does nothing.
func (g *Game) warnYaw() {
	if g.YawWarned {
		return
	}
	g.YawWarned = true
	g.Message("NO YAW IN CLASSIC FLIGHT", "info")
}

func (g *Game) Update(dt float64) {
	switch g.State {
	case "dead":
		if g.Dead != nil {
			g.Dead.T += dt
		}
		return
	case "docked":
		g.Time += dt
		return
	case "hyper":
		g.updateHyper(dt)
		return
	case "docking":
		g.Time += dt
		g.updateDocking(dt)
		return
	case "flight":
	default:
		return
	}
	g.Time += dt
	p := &g.Player
	cmd := &g.Cmd
	w := g.World

	// ---- flight controls / docking computer
	if g.Autopilot && g.AutopilotPath != nil {
		g.runAutopilot(dt)
	} else {
		if g.In.Pitch != 0 {
			p.Mat = p.Mat.Mul(render.RotX(g.In.Pitch * PitchRate * g.Ship().Pitch * dt))
		}
		if g.In.Roll != 0 {
			p.Mat = p.Mat.Mul(render.RotZ(-g.In.Roll * RollRate * g.Ship().Roll * dt))
		}
		// yaw did not exist in the 1984 flight model, so Classic ignores it
		if g.In.Yaw != 0 {
			if g.Classic() {
				g.warnYaw()
			} else {
				// negative: a positive yaw input means "nose to the right"
				p.Mat = p.Mat.Mul(render.RotY(-g.In.Yaw * YawRate * g.Ship().Yaw * dt))
			}
		}
		if g.In.Pitch != 0 || g.In.Roll != 0 || (!g.Classic() && g.In.Yaw != 0) {
			p.Mat = p.Mat.Orthonormalize()
		}
		if g.In.Accel > 0 {
			p.Speed = math.Min(g.Ship().Speed, p.Speed+g.Ship().Accel*dt)
		} else if g.In.Accel < 0 {
			// braking is the same for every hull: the retro thrusters are not
			// what a ship is sold on
			p.Speed = math.Max(0, p.Speed-300*dt)
		}
		fwd := p.Mat.Col(2).Mul(-1)
		p.Pos = p.Pos.Add(fwd.Mul(p.Speed * dt))
	}

	// ---- energy / shields / heat
	recharge := 1.7
	if cmd.Equip["extraenergy"] {
		recharge *= 1.6
	}
	p.Energy = math.Min(g.Ship().Energy, p.Energy+recharge*dt)
	p.ShieldF = math.Min(g.Ship().Shield, p.ShieldF+1.9*dt)
	p.ShieldA = math.Min(g.Ship().Shield, p.ShieldA+1.9*dt)
	p.LaserTemp = math.Max(0, p.LaserTemp-0.35*dt)
	if p.Invulnerable > 0 {
		p.Invulnerable -= dt
	}
	if g.In.Fire && !g.Autopilot && p.Cooldown <= 0 {
		p.Cooldown = LaserCooldown[cmd.Laser]
		g.FireLaser()
		g.Emit("laser")
	}
	if p.Cooldown > 0 {
		p.Cooldown -= dt
	}

	// ---- planet
	dPlanet := p.Pos.Sub(w.Planet.Pos).Len()
	if dPlanet < w.Planet.Radius*1.02 {
		g.LastDamageReason = "BURNED UP IN THE ATMOSPHERE"
		g.PlayerDeath()
		return
	} else if dPlanet < w.Planet.Radius*1.18 {
		g.Message("WARNING: ATMOSPHERE", "bad")
		g.DamagePlayer(14*dt, true)
	}
	dSun := p.Pos.Sub(w.Sun.Pos).Len()
	if dSun < w.Sun.Radius*1.05 {
		g.LastDamageReason = "CONSUMED BY THE STAR"
		g.PlayerDeath()
		return
	}
	if cmd.Equip["fuelscoop"] && dSun < w.Sun.Radius*2.2 && cmd.Fuel < FuelMax {
		cmd.Fuel = math.Min(FuelMax, math.Round((cmd.Fuel+0.35*dt)*100)/100)
		if randFloat() < dt*0.8 {
			g.Message("FUEL SCOOPING", "good")
		}
	}

	// ---- cabin temperature: the star cooks you slowly, cooling is radiative
	if w.Sun.Radius > 0 {
		near := w.Sun.Radius * 4
		heat := 0.0
		if dSun < near {
			heat = (1 - dSun/near) * 0.55
		}
		p.Cabin += (heat - 0.11) * dt
		if p.Cabin < 0 {
			p.Cabin = 0
		}
		if p.Cabin > 0.8 {
			g.DamagePlayer(9*dt*(p.Cabin-0.8)/0.2, true)
			if randFloat() < dt*1.5 {
				g.Message("CABIN TEMPERATURE CRITICAL", "bad")
			}
		}
		if p.Cabin >= 1 {
			g.LastDamageReason = "CABIN OVERHEATED"
			g.PlayerDeath()
			return
		}
	}

	// ---- station (none in witchspace)
	if st := w.Station; st != nil {
		st.Mat = st.Mat.Mul(render.RotY(0.05 * dt))
		dst := p.Pos.Sub(st.Pos).Len()
		if dst < StationSize*1.06 {
			if !g.TryDock() {
				local := st.Mat.ApplyT(p.Pos.Sub(st.Pos))
				outsideSlot := math.Abs(local[0]) > SlotW*StationSize+10 || math.Abs(local[1]) > SlotH*StationSize+10
				if dst < StationSize*1.02 && outsideSlot {
					g.LastDamageReason = "CRASHED INTO THE STATION"
					g.PlayerDeath()
					return
				}
			}
		} else if g.TryDock() {
			return
		}
	}

	g.updateEntities(dt)

	// ---- outbound traffic from the station slot
	if w.Station != nil && !w.Witchspace {
		g.TrafficTimer -= dt
		if g.TrafficTimer <= 0 {
			g.TrafficTimer = randRange(9, 22)
			if randFloat() < 0.55 && g.countBehaviour("depart") < 3 {
				g.spawnDeparture()
			}
		}
	}

	// ---- encounters
	g.EncounterTimer -= dt
	if g.EncounterTimer <= 0 {
		g.EncounterTimer = randRange(25, 55)
		danger := dangerLevel(w.Sys)
		switch {
		case g.CarryingContraband() && g.countBehaviour("police") < 2 && randFloat() < 0.55:
			// contraband in the hold makes the police much more interested
			g.spawnPolice(1)
		case randFloat() < danger*0.7 && g.countBehaviour("pirate") < 4:
			g.SpawnPirates(1 + randInt(2))
		case cmd.Legal >= 100 && g.countBehaviour("police") < 3 && randFloat() < 0.5:
			g.spawnPolice(1)
		case randFloat() < 0.25:
			g.SpawnTraffic(1)
		case randFloat() < 0.12 && g.countBehaviour("asteroid") < 6:
			for a := 0; a < 2+randInt(3); a++ {
				g.SpawnAsteroid(randRange(2000, 6000))
			}
		}
	}

	for i := len(g.Messages) - 1; i >= 0; i-- {
		g.Messages[i].T += dt
		if g.Messages[i].T > g.Messages[i].Life {
			g.Messages = append(g.Messages[:i], g.Messages[i+1:]...)
		}
	}
	if g.Shake > 0 {
		g.Shake -= dt * 3
	}
	if g.EcmFlash > 0 {
		g.EcmFlash -= dt
	}
	if g.BombFlash > 0 {
		g.BombFlash -= dt
	}
	if g.BeamFlash > 0 {
		g.BeamFlash -= dt
	}
}

// MouthDepth is how far the docking shaft runs into the hull; it is also the
// distance the ship has to travel inside before the dock counts.
const MouthDepth = 78.0

func (g *Game) updateHyper(dt float64) {
	h := g.Hyper
	if h == nil {
		g.State = "flight"
		return
	}
	h.T += dt
	if h.Phase == 0 && h.T > 1.6 {
		h.Phase = 1
	}
	if h.Phase == 1 && h.T > 3.2 {
		target := h.Target
		misjump := h.Misjump
		g.EnterSystem(target, "arrive", 9000, misjump)
	}
}

func (g *Game) turnTowards(e *Entity, target render.Vec3, dt float64) (ax, ay float64) {
	local := e.Mat.ApplyT(target)
	ax = math.Atan2(local[0], -local[2])
	ay = math.Atan2(local[1], math.Sqrt(local[0]*local[0]+local[2]*local[2]))
	roll := math.Max(-1, math.Min(1, ax*2.6))
	pitch := math.Max(-1, math.Min(1, ay*1.8))
	// roll the target into the vertical plane first, then pull towards it
	if math.Abs(ax) > 0.55 {
		pitch *= 0.25
	}
	e.Mat = e.Mat.Mul(render.RotZ(roll * 1.9 * dt))
	e.Mat = e.Mat.Mul(render.RotX(pitch * 1.7 * dt))
	e.Mat = e.Mat.Orthonormalize()
	return
}

func (g *Game) updateEntities(dt float64) {
	p := &g.Player
	w := g.World
	ents := w.Entities
	for i := len(ents) - 1; i >= 0; i-- {
		e := ents[i]
		e.Age += dt

		if e.Kind == "explosion" {
			if e.Age > e.Life {
				ents = append(ents[:i], ents[i+1:]...)
			}
			continue
		}

		fwd := e.Mat.Col(2).Mul(-1)

		if e.Kind == "missile" {
			g.updateMissile(e, dt)
			if e.Dead {
				ents = append(ents[:i], ents[i+1:]...)
			}
			continue
		}

		if e.Kind == "canister" || e.Kind == "capsule" || e.Kind == "plate" {
			e.Pos = e.Pos.Add(fwd.Mul(e.Speed * dt))
			e.Mat = e.Mat.Mul(render.RotY(0.4 * dt))
			dc := p.Pos.Sub(e.Pos).Len()
			if dc < 40 {
				if g.Cmd.Equip["fuelscoop"] {
					g.scoopCargo(e)
					ents = append(ents[:i], ents[i+1:]...)
					continue
				} else if dc < 22 {
					g.DamagePlayer(2, true)
					ents = append(ents[:i], ents[i+1:]...)
					g.Message("CARGO CANISTER HIT HULL", "bad")
					continue
				}
			}
			if e.Age > 180 {
				ents = append(ents[:i], ents[i+1:]...)
			}
			continue
		}

		if e.Kind == "asteroid" || e.Kind == "boulder" || e.Kind == "splinter" {
			e.Pos = e.Pos.Add(fwd.Mul(e.Speed * dt))
			e.Mat = e.Mat.Mul(render.RotX(e.Spin[0] * dt))
			e.Mat = e.Mat.Mul(render.RotY(e.Spin[1] * dt))
			e.Mat = e.Mat.Mul(render.RotZ(e.Spin[2] * dt))
			if e.Kind == "splinter" {
				// ore chips are scooped like cargo, not rammed
				if p.Pos.Sub(e.Pos).Len() < 45 {
					if g.Cmd.Equip["fuelscoop"] {
						g.scoopCargo(e)
						ents = append(ents[:i], ents[i+1:]...)
						continue
					}
					if p.Pos.Sub(e.Pos).Len() < 24 {
						g.DamagePlayer(1, true)
						ents = append(ents[:i], ents[i+1:]...)
						continue
					}
				}
				if e.Age > 240 {
					ents = append(ents[:i], ents[i+1:]...)
				}
				continue
			}
			if p.Pos.Sub(e.Pos).Len() < e.Radius*0.8 {
				if e.Kind == "boulder" {
					g.LastDamageReason = "COLLIDED WITH A BOULDER"
				} else {
					g.LastDamageReason = "COLLIDED WITH AN ASTEROID"
				}
				g.PlayerDeath()
				return
			}
			continue
		}

		if e.Kind != "ship" {
			continue
		}

		toPlayer := p.Pos.Sub(e.Pos)
		distP := toPlayer.Len()

		switch e.Behaviour {
		case "depart":
			// outbound traffic runs straight out of the slot and leaves
			e.Speed = e.MaxSpeed * 0.9
			if e.Age > 30 {
				ents = append(ents[:i], ents[i+1:]...)
				continue
			}
		case "trader":
			if w.Station == nil {
				e.Speed = e.MaxSpeed * 0.5
				break
			}
			toStation := w.Station.Pos.Sub(e.Pos)
			ds := toStation.Len()
			g.turnTowards(e, toStation.Norm(), dt)
			if ds > 1200 {
				e.Speed = e.MaxSpeed * 0.85
			} else {
				e.Speed = e.MaxSpeed * 0.4
			}
			if ds < StationSize*2.2 {
				ents = append(ents[:i], ents[i+1:]...)
				continue
			}
		case "flee":
			g.turnTowards(e, e.Pos.Sub(p.Pos).Norm(), dt)
			e.Speed = e.MaxSpeed
		case "police":
			ax, ay := g.turnTowards(e, toPlayer.Norm(), dt)
			if distP > 900 {
				e.Speed = e.MaxSpeed
			} else {
				e.Speed = e.MaxSpeed * 0.7
			}
			e.Cooldown -= dt
			if e.Cooldown <= 0 && distP < 3500 && math.Abs(ax) < 0.2 && math.Abs(ay) < 0.2 {
				e.Cooldown = e.FireRate * randRange(0.8, 1.5)
				if randFloat() < math.Max(0.25, 1.05-distP/5200) {
					g.DamagePlayer(randRange(2, 4.5), true)
					g.Emit("hit")
					g.Shake = 0.35
				}
			}
		default:
			ax, ay := g.turnTowards(e, toPlayer.Norm(), dt)
			// keep a stand-off distance instead of ramming
			frac := (distP - 220) / 1800
			if frac > 1 {
				frac = 1
			}
			if frac < 0.18 {
				frac = 0.18
			}
			desired := e.MaxSpeed * frac
			e.Speed += (desired - e.Speed) * math.Min(1, dt*2.2)
			e.Cooldown -= dt
			if e.CanMissile && e.Cooldown <= 0 && distP > 900 && distP < 5000 &&
				math.Abs(ax) < 0.3 && math.Abs(ay) < 0.3 && randFloat() < 0.25 {
				e.CanMissile = false
				g.EnemyMissile(e)
				e.Cooldown = 2.5
			}
			if e.Cooldown <= 0 && distP < 4000 && math.Abs(ax) < 0.22 && math.Abs(ay) < 0.22 {
				e.Cooldown = e.FireRate * randRange(1.0, 2.0)
				if randFloat() < math.Max(0.18, 0.9-distP/6000) {
					dmg := randRange(1.5, 4.0)
					if e.Type == "thargoid" {
						dmg = randRange(6, 11)
					}
					g.DamagePlayer(dmg, true)
					if e.Type == "thargoid" {
						g.LastDamageReason = "DESTROYED BY THE SWARM"
					} else {
						g.LastDamageReason = "DESTROYED BY PIRATE"
					}
					g.Emit("hit")
					g.Shake = 0.35
				}
			}
		}

		e.Pos = e.Pos.Add(e.Mat.Col(2).Mul(-e.Speed * dt))
		if w.Station != nil && e.Pos.Sub(w.Station.Pos).Len() < StationSize*1.3 && e.Speed > 100 {
			e.Mat = e.Mat.Mul(render.RotX(0.5 * dt))
		}
	}
	w.Entities = ents
}

func (g *Game) updateMissile(m *Entity, dt float64) {
	steer := func(target render.Vec3) {
		aim := target.Sub(m.Pos)
		d := aim.Len()
		if d < 1e-6 {
			return
		}
		local := m.Mat.ApplyT(aim.Norm())
		ax := math.Atan2(local[0], -local[2])
		ay := math.Atan2(local[1], math.Sqrt(local[0]*local[0]+local[2]*local[2]))
		m.Mat = m.Mat.Mul(render.RotZ(math.Max(-1, math.Min(1, ax*2.4)) * 2.2 * dt))
		m.Mat = m.Mat.Mul(render.RotX(math.Max(-1, math.Min(1, ay*2.4)) * 2.2 * dt))
		m.Mat = m.Mat.Orthonormalize()
		m.Speed = math.Min(m.MaxSpeed, m.Speed+700*dt)
		m.Pos = m.Pos.Add(m.Mat.Col(2).Mul(-m.Speed * dt))
		m.Age += dt
		if m.Age > m.Life {
			m.Dead = true
		}
	}

	if m.TargetPlayer {
		if m.Pos.Sub(g.Player.Pos).Len() < 70 {
			g.DamagePlayer(70, false)
			g.addExplosion(m.Pos, 18)
			g.LastDamageReason = "DESTROYED BY A MISSILE"
			g.Emit("explode")
			m.Dead = true
			return
		}
		steer(g.Player.Pos)
		return
	}
	if m.TargetEnt == nil || m.TargetEnt.Dead {
		m.TargetEnt = g.pickMissileTarget()
	}
	if m.TargetEnt != nil {
		if m.Pos.Sub(m.TargetEnt.Pos).Len() < 80 {
			g.HitEntity(m.TargetEnt, 120, "missile")
			g.addExplosion(m.Pos, 20)
			m.Dead = true
			return
		}
		steer(m.TargetEnt.Pos)
		return
	}
	// no target: fly straight
	m.Speed = math.Min(m.MaxSpeed, m.Speed+700*dt)
	m.Pos = m.Pos.Add(m.Mat.Col(2).Mul(-m.Speed * dt))
	m.Age += dt
	if m.Age > m.Life {
		m.Dead = true
	}
}

func (g *Game) scoopCargo(e *Entity) {
	if CargoUsed(g.Cmd) >= g.Cmd.CargoCap {
		g.Message("CARGO HOLD FULL", "bad")
		return
	}
	g.Cmd.Cargo[e.CargoType]++
	name := "CARGO"
	if e.CargoType >= 0 && e.CargoType < len(galaxy.Commodities) {
		name = galaxy.Commodities[e.CargoType].Name
	}
	g.Message("SCOOPED: "+name, "good")
}

func CargoUsed(c Commander) int {
	n := 0
	for _, q := range c.Cargo {
		n += q
	}
	return n
}

// ---------------------------------------------------------------- trade

func (g *Game) Buy(idx, qty int) {
	cmd := &g.Cmd
	if idx < 0 || idx >= len(g.Market) {
		return
	}
	m := &g.Market[idx]
	free := cmd.CargoCap - CargoUsed(*cmd)
	if qty > m.Qty {
		qty = m.Qty
	}
	if qty > free {
		qty = free
	}
	if m.Price > 0 {
		if maxAfford := int(cmd.Cash) / m.Price; qty > maxAfford {
			qty = maxAfford
		}
	}
	if qty <= 0 {
		g.Message("CANNOT BUY", "bad")
		return
	}
	cmd.Cash -= float64(qty * m.Price)
	cmd.Cargo[idx] += qty
	m.Qty -= qty
	g.Message("BOUGHT "+itoa(qty)+"T "+m.Name, "good")
}

func (g *Game) Sell(idx, qty int) {
	cmd := &g.Cmd
	if idx < 0 || idx >= len(g.Market) {
		return
	}
	m := &g.Market[idx]
	have := cmd.Cargo[idx]
	if qty > have {
		qty = have
	}
	if qty <= 0 {
		g.Message("NOTHING TO SELL", "bad")
		return
	}
	cmd.Cash += float64(qty * m.Sell)
	cmd.Cargo[idx] = have - qty
	m.Qty += qty
	g.Message("SOLD "+itoa(qty)+"T "+m.Name+" FOR "+itoa(qty*m.Sell)+" CR", "good")
}

func (g *Game) Refuel() {
	cmd := &g.Cmd
	need := FuelMax - cmd.Fuel
	if need <= 0.01 {
		g.Message("TANKS ALREADY FULL", "bad")
		return
	}
	cost := int(math.Ceil(need * galaxy.FuelPricePerLY))
	if cmd.Cash < float64(cost) {
		canBuy := cmd.Cash / galaxy.FuelPricePerLY
		if canBuy < 0.1 {
			g.Message("NOT ENOUGH CREDITS", "bad")
			return
		}
		cmd.Fuel = math.Round((cmd.Fuel+canBuy)*10) / 10
		cmd.Cash = 0
		g.Message("PARTIAL REFUEL", "info")
		return
	}
	cmd.Cash -= float64(cost)
	cmd.Fuel = FuelMax
	g.Message("REFUELLED FOR "+itoa(cost)+" CR", "good")
}

func (g *Game) BuyEquipment(id string) {
	cmd := &g.Cmd
	var item *galaxy.Equipment
	for i := range galaxy.EquipmentList {
		if galaxy.EquipmentList[i].ID == id {
			item = &galaxy.EquipmentList[i]
		}
	}
	if item == nil {
		return
	}
	if g.World != nil && g.World.Sys.Tech < item.Tech {
		g.Message("NOT AVAILABLE HERE", "bad")
		return
	}
	if cmd.Cash < float64(item.Price) {
		g.Message("NOT ENOUGH CREDITS", "bad")
		return
	}
	switch id {
	case "missile":
		if cmd.Missiles >= MissilesMax {
			g.Message("MISSILE RACKS FULL", "bad")
			return
		}
		cmd.Missiles++
	case "laserpulse", "laserbeam", "lasermilitary":
		lvl := 0
		if id == "laserbeam" {
			lvl = 1
		} else if id == "lasermilitary" {
			lvl = 2
		}
		if cmd.Laser >= lvl {
			g.Message("ALREADY FITTED", "bad")
			return
		}
		if lvl > g.Ship().MaxLaser {
			g.Message("THIS HULL CANNOT MOUNT IT", "bad")
			return
		}
		cmd.Laser = lvl
	case "mininglaser":
		cmd.Equip["mining"] = true
		cmd.Laser = 3
	case "largebay":
		if cmd.Equip["largebay"] {
			g.Message("ALREADY FITTED", "bad")
			return
		}
		cmd.Equip["largebay"] = true
		cmd.CargoCap = 35
	case "bomb":
		cmd.Bombs++
		cmd.Equip["bomb"] = true
	default:
		if cmd.Equip[id] {
			g.Message("ALREADY FITTED", "bad")
			return
		}
		cmd.Equip[id] = true
	}
	cmd.Cash -= float64(item.Price)
	g.Message("PURCHASED: "+item.Name, "good")
}

// NewGame starts a fresh commander docked at the first system.
func (g *Game) NewGame(name string) { g.NewGameMode(ModeAdvanced, name) }

// NewGameMode starts a commander under the given rule set.
func (g *Game) NewGameMode(mode, name string) {
	g.Cmd = NewCommander()
	g.Cmd.Mode = mode
	if name != "" {
		g.Cmd.Name = name
	}
	g.Cmd.System = 0
	g.Cmd.Galaxy = 0
	g.EscapeUsed = false
	g.EnterSystem(0, "arrive", 9000, false)
	g.State = "docked"
	g.dockedAt()
}
