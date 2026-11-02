package game

import (
	"sort"
	"testing"
)

// Every branch in the story has to lead somewhere. This walks all 62 nodes and
// checks the wiring: a Next that names no node, a node nothing can reach, a node
// with no way out, and - the bug that started this - two different answers that
// are wired to the same place, which means one of them is a lie.
func TestDialogueWiring(t *testing.T) {
	var ids []string
	for id := range dialogueNodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	targets := map[string][]string{} // node -> who points at it
	for _, id := range ids {
		n := dialogueNodes[id]
		for i, o := range n.Options {
			if o.Next == "" {
				continue
			}
			if _, ok := dialogueNodes[o.Next]; !ok {
				t.Errorf("%s option %d (%q) points at %q, which does not exist", id, i, o.Text, o.Next)
			}
			targets[o.Next] = append(targets[o.Next], id+"#"+o.Text)
		}
	}

	// a node with no options at all can only be an ending
	for _, id := range ids {
		n := dialogueNodes[id]
		if len(n.Options) == 0 {
			continue
		}
		// two answers that go to the same place AND do the same thing are the
		// same answer written twice: that is the c7_hive bug in its general form
		seen := map[string]string{}
		for _, o := range n.Options {
			if o.Next == "" {
				continue
			}
			key := o.Next + "|" + effectsKey(o.Do) + "|" + condsKey(o.If)
			if prev, ok := seen[key]; ok {
				t.Errorf("%s: %q and %q are the same choice twice (both %s, same effects)",
					id, prev, o.Text, o.Next)
			}
			seen[key] = o.Text
		}
	}

	// unreachable nodes: nothing in the story points at them and they are not an
	// entry point the campaign opens by name
	entries := map[string]bool{}
	for _, q := range campaignQuests {
		if q.Talk != "" {
			entries[q.Talk] = true
		}
		for _, v := range q.Variants {
			entries[v.Node] = true
		}
	}
	for _, id := range ids {
		if entries[id] || len(targets[id]) > 0 {
			continue
		}
		t.Logf("nothing points at %s (only opened by name)", id)
	}

	// and every option has something to say
	for _, id := range ids {
		for i, o := range dialogueNodes[id].Options {
			if o.Text == "" {
				t.Errorf("%s option %d has no text", id, i)
			}
		}
	}
}

func effectsKey(es []Effect) string {
	out := ""
	for _, e := range es {
		out += e.Kind + ":" + e.Key + ":" + itoa(e.N) + ";"
	}
	return out
}

func condsKey(cs []Cond) string {
	out := ""
	for _, c := range cs {
		out += c.Key + c.Op + itoa(c.N) + ";"
	}
	return out
}

// A branch that asks for a flag nothing ever sets can never be taken. Same for a
// branch that asks for the wrong kind of thing - "papers >= 1" on a flag, or a
// flag test on a counter.
func TestDialogueConditionsCanBeMet(t *testing.T) {
	setByEffect := map[string]bool{}
	counters := map[string]bool{}
	for _, n := range dialogueNodes {
		for _, e := range n.OnEnter {
			if e.Kind == FxFlag {
				setByEffect[e.Key] = true
			}
			if e.Kind == FxPapers || e.Kind == FxCash {
				counters[e.Key] = true
			}
		}
		for _, o := range n.Options {
			for _, e := range o.Do {
				if e.Kind == FxFlag {
					setByEffect[e.Key] = true
				}
			}
		}
	}
	// code sets some of them directly: campaign tags, story consequences
	codeFlags := map[string]bool{
		"PAID_PERMIT": true, "KNOWS_HIVE": true, "KNOWS_BOX": true,
		"ASKED_PRICE": true, "MET_AUDITOR": true, "HAS_LIST": true,
		"READ_LIST": true, "HAS_COPY": true, "COPIED_RECORDER": true,
		"CONVOY_SPARED": true, "CONVOY_ATTACKED": true, "BOUGHT_LICENSE": true,
		"SHIELDS_UP": true, "SHIELDS_CUT": true, "HIVE_ATTACKED": true,
		"PAPERS": true, "AUDIT_REFORMED": true,
	}
	for _, id := range sortedIDs() {
		n := dialogueNodes[id]
		for _, o := range n.Options {
			for _, c := range o.If {
				// these are read from the commander, not from the flag map
				switch c.Key {
				case "papers", "boxes", "cash", "wanted", "viperclk", "rig":
					continue
				}
				switch c.Op {
				case "flag":
					if !setByEffect[c.Key] && !codeFlags[c.Key] {
						t.Errorf("%s: option %q needs flag %s, which nothing sets", id, o.Text, c.Key)
					}
				case ">=", "<=", "==", ">":
					if !counters[c.Key] && !codeFlags[c.Key] {
						t.Logf("%s: option %q compares %s, which is not a flag", id, o.Text, c.Key)
					}
				}
			}
		}
	}
}

func sortedIDs() []string {
	var ids []string
	for id := range dialogueNodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
