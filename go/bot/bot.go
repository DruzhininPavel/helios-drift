// Package bot drives the game with a policy instead of a keyboard: it plays,
// records what it did and checks the simulation for impossible states. A run
// is reproducible from its seed plus its recorded inputs, which is what makes
// a bug found by the bot debuggable.
package bot

import (
	"encoding/json"
	"fmt"
	"os"

	"heliosdrift/galaxy"
	"heliosdrift/game"
)

// Event is a discrete action taken at a frame: the things a player does with a
// keypress rather than with the stick.
type Event struct {
	Frame int    `json:"frame"`
	Kind  string `json:"kind"`
	A     int    `json:"a,omitempty"`
	B     int    `json:"b,omitempty"`
}

// InputRec is a change of the continuous input; frames between records keep
// the previous value, so a script stays small.
type InputRec struct {
	Frame int     `json:"frame"`
	Pitch float32 `json:"pitch"`
	Roll  float32 `json:"roll"`
	Yaw   float32 `json:"yaw"`
	Accel float32 `json:"accel"`
	Fire  bool    `json:"fire"`
}

// Script is a full recording of a run.
type Script struct {
	Seed   int64      `json:"seed"`
	Mode   string     `json:"mode"`
	Frames int        `json:"frames"`
	Inputs []InputRec `json:"inputs"`
	Events []Event    `json:"events"`
	Hash   uint64     `json:"hash"`
	Hashes []uint64   `json:"hashes,omitempty"`
	Report Report     `json:"report"`
}

// Policy decides what to do at a frame.
type Policy interface {
	Act(g *game.Game, frame int) (game.InputState, []Event)
	Name() string
}

// Violation is one impossible state caught by the checker.
type Violation struct {
	Frame int    `json:"frame"`
	What  string `json:"what"`
	Value string `json:"value"`
}

// Report summarises a run.
type Report struct {
	Policy     string         `json:"policy"`
	Frames     int            `json:"frames"`
	Hash       uint64         `json:"hash"`
	Steps      map[string]int `json:"steps"`
	States     map[string]int `json:"states"`
	Violations []Violation    `json:"violations"`
	Stuck      string         `json:"stuck,omitempty"`
	Death      string         `json:"death,omitempty"`
	DeathFrame int            `json:"death_frame,omitempty"`
	CashStart  float64        `json:"cash_start"`
	CashEnd    float64        `json:"cash_end"`
	Docks      int            `json:"docks"`
	Ended      string         `json:"ended"`
	Phases     map[string]int `json:"phases,omitempty"`
	Distances  []float64      `json:"distances,omitempty"`
	Speeds     []float64      `json:"speeds,omitempty"`
}

const DT = 1.0 / 60.0

// TraceFrame, when set, sees every frame of a run. Tests use it to watch a run
// from the inside without changing it.
var TraceFrame func(frame int, g *game.Game, in game.InputState)

// quantize rounds the stick input to the precision the recording keeps, so a
// replayed run sees exactly the same numbers as the live one. Without this the
// float32 file format silently changes the flight and the replay drifts.
func quantize(in game.InputState) game.InputState {
	return game.InputState{
		Pitch: float64(float32(in.Pitch)),
		Roll:  float64(float32(in.Roll)),
		Yaw:   float64(float32(in.Yaw)),
		Accel: float64(float32(in.Accel)),
		Fire:  in.Fire,
	}
}

// Run plays the game for the given number of frames with the policy, recording
// every input change and event.
func Run(g *game.Game, frames int, pol Policy) *Script {
	sc := &Script{Seed: game.Seed(), Mode: g.Cmd.Mode, Frames: frames}
	sc.Report = Report{Policy: pol.Name(), Steps: map[string]int{}, States: map[string]int{},
		Phases: map[string]int{}, CashStart: g.Cmd.Cash}
	last := game.InputState{Pitch: 1e9} // force the first record
	quiet := 0
	prevState := g.State
	lastProgress := g.Hash64()

	for f := 0; f < frames; f++ {
		in, events := pol.Act(g, f)
		in = quantize(in)
		if in != last {
			sc.Inputs = append(sc.Inputs, InputRec{
				Frame: f, Pitch: float32(in.Pitch), Roll: float32(in.Roll),
				Yaw: float32(in.Yaw), Accel: float32(in.Accel), Fire: in.Fire,
			})
			last = in
		}
		for _, e := range events {
			e.Frame = f
			sc.Events = append(sc.Events, e)
			sc.Report.Steps[e.Kind]++
			applyEvent(g, e)
		}
		g.In = in
		g.Update(DT)
		if TraceFrame != nil {
			TraceFrame(f, g, in)
		}
		sc.Report.States[g.State]++
		sc.Hashes = append(sc.Hashes, g.StateHash())

		if g.State == "docked" && prevState != "docked" {
			sc.Report.Docks++
		}
		prevState = g.State
		if g.State == "dead" && sc.Report.Death == "" {
			sc.Report.Death = g.Dead.Reason
			sc.Report.DeathFrame = f
		}
		if f%300 == 0 {
			if st := g.World.Station; st != nil {
				sc.Report.Distances = append(sc.Report.Distances, float64(int(g.Player.Pos.Sub(st.Pos).Len())))
				sc.Report.Speeds = append(sc.Report.Speeds, float64(int(g.Player.Speed)))
			}
		}
		if tp, ok := pol.(interface{ PhaseName() string }); ok {
			sc.Report.Phases[tp.PhaseName()]++
		}
		if v := check(g, f); v != nil {
			sc.Report.Violations = append(sc.Report.Violations, *v)
			if len(sc.Report.Violations) > 20 {
				break
			}
		}
		// a long stretch with no change in the world is a softlock candidate
		if h := g.Hash64(); h != lastProgress {
			lastProgress = h
			quiet = 0
		} else if quiet++; quiet == 60*180 {
			sc.Report.Stuck = fmt.Sprintf("no change since frame %d", f-60*180)
			break
		}
	}
	sc.Report.Frames = len(sc.Hashes)
	sc.Hash = g.StateHash()
	sc.Report.Hash = sc.Hash
	sc.Report.CashEnd = g.Cmd.Cash
	sc.Report.Ended = g.State
	return sc
}

// Replay applies a recorded script to a fresh game. It returns the final hash
// and, when the run drifts from the recording, the frame where it happened.
func Replay(g *game.Game, sc *Script) (uint64, int) {
	diverge := -1
	in := game.InputState{}
	i, e := 0, 0
	for f := 0; f < sc.Frames; f++ {
		for i < len(sc.Inputs) && sc.Inputs[i].Frame == f {
			r := sc.Inputs[i]
			in = game.InputState{Pitch: float64(r.Pitch), Roll: float64(r.Roll),
				Yaw: float64(r.Yaw), Accel: float64(r.Accel), Fire: r.Fire}
			i++
		}
		for e < len(sc.Events) && sc.Events[e].Frame == f {
			applyEvent(g, sc.Events[e])
			e++
		}
		g.In = in
		g.Update(DT)
		if diverge < 0 && f < len(sc.Hashes) && g.StateHash() != sc.Hashes[f] {
			diverge = f
		}
	}
	return g.StateHash(), diverge
}

func applyEvent(g *game.Game, e Event) {
	switch e.Kind {
	case "launch":
		g.Launch()
	case "dock":
		g.Dock()
	case "refuel":
		g.Refuel()
	case "jump":
		g.JumpInSystem()
	case "hyper":
		g.Hyperspace(e.A)
	case "galactic":
		g.GalacticJump()
	case "missile":
		g.FireMissile()
	case "ecm":
		g.FireECM()
	case "bomb":
		g.FireBomb()
	case "autopilot":
		g.ToggleAutopilot()
	case "jettison":
		g.Jettison()
	case "buy":
		n := e.B
		if n == 0 {
			n = 1
		}
		g.Buy(e.A, n)
	case "sell":
		n := e.B
		if n == 0 {
			n = 1
		}
		g.Sell(e.A, n)
	case "campaign":
		g.AcceptCampaign()
	case "contract":
		g.AcceptOffer(e.A)
	case "equip":
		g.BuyEquipment(equipID(e.A))
	}
}

// equipID turns the equipment index used by events into an id.
func equipID(i int) string {
	if i < 0 || i >= len(galaxy.EquipmentList) {
		return ""
	}
	return galaxy.EquipmentList[i].ID
}

// Save writes the script as JSON.
func (sc *Script) Save(path string) error {
	raw, err := json.MarshalIndent(sc, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

// Load reads a script from JSON.
func Load(path string) (*Script, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sc Script
	if err := json.Unmarshal(raw, &sc); err != nil {
		return nil, err
	}
	return &sc, nil
}
