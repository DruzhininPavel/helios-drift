package game

import "heliosdrift/lang"

// Dialogue is data, not a scripting language: a node has lines and options, an
// option has conditions and effects. The story lives in dialogue_data.go and
// reads like the spec, and the screen builds its text from keys so Russian and
// English come out of the same node.

type Cond struct {
	Key string // a flag, or one of "papers" "boxes" "cash" "wanted" "viperclk"
	Op  string // "flag" "!flag" "==" ">=" "<="
	N   int
}

type Effect struct {
	Kind string // flag, unflag, cash, papers, box, wanted, clock, audit, ending, quest
	Key  string
	N    int
	S    string
}

type Option struct {
	Text string // string key
	Next string // next node, "" ends the conversation
	If   []Cond
	Do   []Effect
}

// LineVariant swaps the body text when its conditions hold: the same screen
// reads differently depending on what the player has done.
type LineVariant struct {
	If    []Cond
	Lines []string
}

type Node struct {
	ID      string
	Kind    string // talk, msg, screen
	Speaker string // "" for system text
	Lines   []string
	// Variants are checked in order before Lines, first match wins.
	Variants []LineVariant
	Options  []Option
	OnEnter  []Effect
}

// Effect kinds.
const (
	FxFlag   = "flag"
	FxUnflag = "unflag"
	FxCash   = "cash"
	FxPapers = "papers"
	FxBox    = "box"
	FxWanted = "wanted"
	FxClock  = "clock"
	FxAudit  = "audit"
	FxEnding = "ending"
	FxQuest  = "quest"
	FxSpawn  = "spawn"
)

var dialogueNodes = map[string]*Node{}

func registerNodes(ns ...*Node) {
	for _, n := range ns {
		dialogueNodes[n.ID] = n
	}
}

// DialogueNode returns a node by id.
func DialogueNode(id string) *Node { return dialogueNodes[id] }

// DialogueNodeIDs lists every node, for the integrity test.
func DialogueNodeIDs() []string {
	out := make([]string, 0, len(dialogueNodes))
	for id := range dialogueNodes {
		out = append(out, id)
	}
	return out
}

// condHolds evaluates a condition against the commander.
func (g *Game) condHolds(c Cond) bool {
	switch c.Key {
	case "papers":
		return compare(g.Cmd.Papers, c.Op, c.N)
	case "boxes":
		return compare(g.Cmd.Boxes, c.Op, c.N)
	case "cash":
		return compare(int(g.Cmd.Cash), c.Op, c.N)
	case "wanted":
		return compare(g.Cmd.Legal, c.Op, c.N)
	case "viperclk":
		return compare(g.Cmd.ViperClk, c.Op, c.N)
	case "rig":
		// the heavy mining rig opens the cabinet by force
		return g.Cmd.Equip["miningrig"]
	}
	switch c.Op {
	case "flag":
		return g.Cmd.Flags[c.Key]
	case "!flag":
		return !g.Cmd.Flags[c.Key]
	}
	return false
}

func compare(v int, op string, n int) bool {
	switch op {
	case "==":
		return v == n
	case ">=":
		return v >= n
	case "<=":
		return v <= n
	case ">":
		return v > n
	case "<":
		return v < n
	}
	return false
}

// conditionsMet reports whether every condition holds.
func (g *Game) conditionsMet(cs []Cond) bool {
	for _, c := range cs {
		if !g.condHolds(c) {
			return false
		}
	}
	return true
}

// VisibleOptions returns the options the player may pick, with their original
// indices, so the screen and the choice agree on numbering.
func (g *Game) VisibleOptions(n *Node) []int {
	var out []int
	if n == nil {
		return out
	}
	for i, o := range n.Options {
		if g.conditionsMet(o.If) {
			out = append(out, i)
		}
	}
	return out
}

// apply runs one effect.
func (g *Game) apply(e Effect) {
	cmd := &g.Cmd
	if cmd.Flags == nil {
		cmd.Flags = map[string]bool{}
	}
	switch e.Kind {
	case FxFlag:
		cmd.Flags[e.Key] = true
		if e.Key == "BOUGHT_LICENSE" {
			// remember where the licence was bought: the discount is local
			cmd.LicenseSys = cmd.System
			cmd.LicenseGalaxy = cmd.Galaxy
		}
	case FxUnflag:
		delete(cmd.Flags, e.Key)
	case FxCash:
		cmd.Cash += float64(e.N)
		if cmd.Cash < 0 {
			cmd.Cash = 0
		}
	case FxPapers:
		cmd.Papers += e.N
	case FxBox:
		cmd.Boxes += e.N
		if cmd.Boxes < 0 {
			cmd.Boxes = 0
		}
	case FxWanted:
		cmd.Legal += e.N
		if cmd.Legal < 0 {
			cmd.Legal = 0
		}
	case FxClock:
		cmd.ViperClk = e.N
	case FxAudit:
		cmd.Audit = e.S
	case FxEnding:
		g.StartEnding(e.S)
	case FxQuest:
		g.advanceCampaign()
	case FxSpawn:
		if e.Key == "PIRATES" {
			// the four who were waiting for the shields to drop. One of them
			// takes payment: if the pilot paid the ambush off in chapter two,
			// that one turns away instead of shooting.
			g.SpawnPirates(4)
			if g.Cmd.Flags["PAID_PERMIT"] {
				for i := len(g.World.Entities) - 1; i >= 0; i-- {
					e := g.World.Entities[i]
					if e.Behaviour == "pirate" {
						e.Behaviour = "depart"
						e.Hostile = false
						break
					}
				}
				g.Message("ONE OF THEM TURNS AWAY", "info")
				g.Note(LogPirate, "SPARED", 0, 0)
			}
			break
		}
		g.spawnQuestGroup(e.Key)
	}
}

// TalkStart opens a conversation. It works at a counter and in flight: the
// screen is the same, and EndTalk puts the player back where he was.
func (g *Game) TalkStart(id string) bool {
	n := DialogueNode(id)
	if n == nil {
		return false
	}
	if g.Screen != "talk" {
		g.TalkBack = g.Screen
	}
	g.Talk = id
	g.Screen = "talk"
	g.TalkLines(n)
	return true
}

// TalkLines applies OnEnter once when the node opens.
func (g *Game) TalkLines(n *Node) {
	for _, e := range n.OnEnter {
		g.apply(e)
	}
}

// TalkChoose applies the effects of the picked option and moves on. It returns
// false when the conversation ended.
func (g *Game) TalkChoose(visibleIndex int) bool {
	n := DialogueNode(g.Talk)
	if n == nil {
		g.EndTalk()
		return false
	}
	vis := g.VisibleOptions(n)
	if visibleIndex < 0 || visibleIndex >= len(vis) {
		return true
	}
	opt := n.Options[vis[visibleIndex]]
	for _, e := range opt.Do {
		g.apply(e)
	}
	if opt.Next == "" {
		g.EndTalk()
		return false
	}
	next := DialogueNode(opt.Next)
	if next == nil {
		g.EndTalk()
		return false
	}
	g.Talk = next.ID
	g.TalkLines(next)
	return true
}

// EndTalk closes the conversation and returns to whatever was behind it.
func (g *Game) EndTalk() {
	g.Talk = ""
	g.Screen = g.TalkBack
	g.TalkBack = ""
	if g.Screen == "" && g.State == "docked" {
		g.Screen = "launch"
	}
	// a conversation that settled the campaign ends on the last page - once.
	// Afterwards the sector is still there to fly in, and every counter would
	// otherwise throw the same epitaph at the pilot again.
	if g.Cmd.Ending != "" && !g.Cmd.Flags["ENDING_SHOWN"] {
		g.Screen = "ending"
	}
}

// NodeLines is the text of a node in this commander's circumstances: a variant
// first, the plain lines otherwise.
func (g *Game) NodeLines(n *Node) []string {
	if n == nil {
		return nil
	}
	for _, v := range n.Variants {
		if g.conditionsMet(v.If) {
			return v.Lines
		}
	}
	return n.Lines
}

// TalkText returns the node's lines, translated for the screen.
func (g *Game) TalkText(n *Node) []string {
	src := g.NodeLines(n)
	out := make([]string, 0, len(src))
	for _, ln := range src {
		out = append(out, lang.T(ln))
	}
	return out
}
