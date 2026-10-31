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

type Node struct {
	ID      string
	Kind    string // talk, msg, screen
	Speaker string // "" for system text
	Lines   []string
	Options []Option
	OnEnter []Effect
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
		cmd.Ending = e.S
	case FxQuest:
		g.advanceCampaign()
	}
}

// TalkStart opens a conversation, applying the node's entry effects.
func (g *Game) TalkStart(id string) bool {
	n := DialogueNode(id)
	if n == nil {
		return false
	}
	g.Talk = id
	g.TalkLines(n)
	g.State = "docked"
	g.Screen = "talk"
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

// EndTalk closes the conversation and returns to the station menu.
func (g *Game) EndTalk() {
	g.Talk = ""
	g.Screen = "launch"
}

// TalkText returns the node's lines, translated for the screen.
func (g *Game) TalkText(n *Node) []string {
	out := make([]string, 0, len(n.Lines))
	for _, ln := range n.Lines {
		out = append(out, lang.T(ln))
	}
	return out
}
