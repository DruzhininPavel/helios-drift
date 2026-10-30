package game

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"sort"
)

// StateHash is a stable fingerprint of everything that matters for gameplay.
// Two runs fed the same seed and the same inputs must produce the same hash;
// when they do not, something in the simulation is not deterministic.
func (g *Game) StateHash() uint64 {
	h := fnv.New64a()
	var buf [8]byte

	f := func(v float64) {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			binary.LittleEndian.PutUint64(buf[:], 0xDEADBEEFDEADBEEF)
			h.Write(buf[:])
			return
		}
		binary.LittleEndian.PutUint64(buf[:], math.Float64bits(v))
		h.Write(buf[:])
	}
	i := func(v int) {
		binary.LittleEndian.PutUint64(buf[:], uint64(uint32(v)))
		h.Write(buf[:])
	}
	s := func(v string) {
		h.Write([]byte(v))
		h.Write([]byte{0})
	}

	s(g.State)
	s(g.Screen)
	f(g.Time)

	c := &g.Cmd
	s(c.Name)
	s(c.Ship)
	s(c.Mode)
	f(c.Cash)
	f(c.Fuel)
	i(c.System)
	i(c.Galaxy)
	i(c.Kills)
	i(c.Legal)
	i(c.Missiles)
	i(c.Deaths)
	i(c.Day)
	i(c.Pax)
	i(c.Mail)
	i(c.CargoCap)
	s(c.Mission.Kind)
	i(c.Mission.System)
	i(c.Mission.Galaxy)
	i(c.MissionFlags)
	if c.Campaign.Active {
		i(1)
	} else {
		i(0)
	}
	i(c.Campaign.Stage)
	i(c.Campaign.Target)
	i(c.Campaign.Galaxy)

	// maps are ordered before hashing, or the result would change run to run
	goods := make([]int, 0, len(c.Cargo))
	for k := range c.Cargo {
		goods = append(goods, k)
	}
	sort.Ints(goods)
	for _, k := range goods {
		i(k)
		i(c.Cargo[k])
	}
	equip := make([]string, 0, len(c.Equip))
	for k := range c.Equip {
		equip = append(equip, k)
	}
	sort.Strings(equip)
	for _, k := range equip {
		s(k)
	}

	p := &g.Player
	f(p.Pos[0])
	f(p.Pos[1])
	f(p.Pos[2])
	f(p.Speed)
	f(p.Energy)
	f(p.ShieldF)
	f(p.ShieldA)
	f(p.Cabin)
	for col := 0; col < 3; col++ {
		v := p.Mat.Col(col)
		f(v[0])
		f(v[1])
		f(v[2])
	}

	if g.World != nil {
		i(len(g.World.Entities))
		for _, e := range g.World.Entities {
			s(e.Type)
			s(e.Behaviour)
			f(float64(e.Energy))
			f(e.Pos[0])
			f(e.Pos[1])
			f(e.Pos[2])
		}
	}
	return h.Sum64()
}

// Hash64 is a cheap per-frame fingerprint, used to notice that the world has
// stopped changing. It covers far less than StateHash on purpose: it runs every
// frame of a bot run.
func (g *Game) Hash64() uint64 {
	h := fnv.New64a()
	var buf [8]byte
	add := func(v float64) {
		binary.LittleEndian.PutUint64(buf[:], math.Float64bits(math.Round(v*1000)/1000))
		h.Write(buf[:])
	}
	h.Write([]byte(g.State))
	add(g.Time)
	add(g.Cmd.Cash)
	add(g.Cmd.Fuel)
	add(float64(g.Cmd.System))
	add(float64(g.Cmd.Kills))
	add(float64(len(g.Cmd.Cargo)))
	if g.World != nil {
		add(float64(len(g.World.Entities)))
		for _, e := range g.World.Entities {
			h.Write([]byte(e.Type))
			add(e.Pos[0])
			add(e.Pos[1])
			add(e.Pos[2])
			add(float64(e.Energy))
		}
	}
	add(g.Player.Pos[0])
	add(g.Player.Pos[1])
	add(g.Player.Pos[2])
	add(g.Player.Speed)
	add(g.Player.Energy)
	return h.Sum64()
}
