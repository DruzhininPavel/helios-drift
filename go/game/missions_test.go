package game

import (
	"testing"
)

func newDocked() *Game {
	g := New()
	g.NewGame("TEST")
	return g
}

// newDockedClassic is the 1984 rule set, where the two original contracts live.
func newDockedClassic() *Game {
	g := New()
	g.NewGameMode(ModeClassic, "TEST")
	return g
}

// The Constrictor contract: offered from AVERAGE rank, spawns its target in the
// right system, pays out on the kill.
func TestConstrictorContract(t *testing.T) {
	g := newDockedClassic()
	g.Cmd.Kills = 40
	g.dockedAt()

	if got := g.MissionOffer(); got != MissionConstrictor {
		t.Fatalf("offer = %q, want %q", got, MissionConstrictor)
	}
	if !g.AcceptMission(MissionConstrictor) {
		t.Fatal("accept failed")
	}
	target := g.Cmd.Mission.System
	if target == g.Cmd.System {
		t.Fatal("mission target is the current system")
	}
	if g.MissionOffer() != "" {
		t.Fatal("offer should be gone once accepted")
	}

	g.EnterSystem(target, "arrive", 9000, false)
	var c *Entity
	for _, e := range g.World.Entities {
		if e.Type == "constrictor" {
			c = e
		}
	}
	if c == nil {
		t.Fatal("constrictor did not spawn in the target system")
	}

	cash := g.Cmd.Cash
	g.DestroyEntity(c, "laser")
	if want := cash + 5000; g.Cmd.Cash != want {
		t.Fatalf("cash = %v, want %v", g.Cmd.Cash, want)
	}
	if g.Cmd.MissionFlags&flagConstrictorDone == 0 {
		t.Fatal("completion flag not set")
	}
	if g.Cmd.Mission.Kind != "" {
		t.Fatal("mission not cleared after completion")
	}
}

// The courier run needs the drive, and the jump consumes it.
func TestDocumentsContract(t *testing.T) {
	g := newDockedClassic()
	g.Cmd.MissionFlags = flagConstrictorDone
	g.dockedAt()

	if got := g.MissionOffer(); got != MissionDocuments {
		t.Fatalf("offer = %q, want %q", got, MissionDocuments)
	}
	if g.AcceptMission(MissionDocuments) {
		t.Fatal("accepted without a galactic hyperdrive")
	}
	g.Cmd.Equip["galdrive"] = true
	if !g.AcceptMission(MissionDocuments) {
		t.Fatal("accept failed with drive fitted")
	}
	wantGalaxy := (g.Cmd.Galaxy + 1) % 8
	if g.Cmd.Mission.Galaxy != wantGalaxy {
		t.Fatalf("mission galaxy = %d, want %d", g.Cmd.Mission.Galaxy, wantGalaxy)
	}

	g.State = "flight"
	if !g.GalacticJump() {
		t.Fatal("galactic jump refused")
	}
	if g.Cmd.Galaxy != wantGalaxy {
		t.Fatalf("galaxy = %d, want %d", g.Cmd.Galaxy, wantGalaxy)
	}
	if g.Cmd.Equip["galdrive"] {
		t.Fatal("drive was not consumed")
	}
	if g.State != "hyper" {
		t.Fatalf("state = %q, want hyper", g.State)
	}
}

// Delivering the documents completes the run at the destination station.
func TestDocumentsDelivery(t *testing.T) {
	g := newDockedClassic()
	g.Cmd.MissionFlags = flagConstrictorDone
	g.dockedAt()
	g.Cmd.Equip["galdrive"] = true
	g.AcceptMission(MissionDocuments)
	m := g.Cmd.Mission
	cash := g.Cmd.Cash

	// the drop-off is in the next galaxy, so the commander has to be there
	g.Cmd.Galaxy = m.Galaxy
	g.EnterSystem(m.System, "arrive", 9000, false)
	dockNow(g)
	if g.Cmd.Cash != cash+float64(m.Pay) {
		t.Fatalf("cash = %v, want %v", g.Cmd.Cash, cash+float64(m.Pay))
	}
	if g.Cmd.MissionFlags&flagDocumentsDone == 0 {
		t.Fatal("documents flag not set")
	}
}

// Every hull in the roster must have a mesh, and vice versa.
func TestShipRosterHasMeshes(t *testing.T) {
	for id, def := range Ships {
		if def.Mesh == "" {
			t.Errorf("%s: empty mesh", id)
			continue
		}
		if meshFor(def.Mesh) == nil {
			t.Errorf("%s: mesh %q not found", id, def.Mesh)
		}
		if def.Radius <= 0 || def.Energy <= 0 {
			t.Errorf("%s: bad stats %+v", id, def)
		}
	}
}

// The original roster, minus craft that only exist as scenery.
func TestOriginalRosterPresent(t *testing.T) {
	want := []string{
		"adder", "anaconda", "asp", "boa", "cobra1", "cobra3", "constrictor",
		"ferdelance", "gecko", "krait", "mamba", "moray", "python", "shuttle",
		"sidewinder", "thargoid", "thargon", "transporter", "viper", "worm",
	}
	for _, id := range want {
		if _, ok := Ships[id]; !ok {
			t.Errorf("missing hull %q", id)
		}
	}
}
