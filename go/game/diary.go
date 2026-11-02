package game

import "sort"

// The diary: the story as the pilot writes it down.
//
// The log records what happened in numbers - a tonne sold, a pirate killed. The
// diary records what it meant: a page per turn of the story, written when the
// turn happens, and different depending on which way the pilot went. Entries are
// never edited or removed: what you decided is what you keep.

// DiaryEntry is one page.
type DiaryEntry struct {
	ID    string
	Order int      // the page's place in the book
	Title string   // the heading on the page
	Lines []string // the text, wrapped by the screen
	// When is what has to be true for the page to exist: the story reached a
	// stage, or a particular choice was made.
	When []Cond
}

// diaryEntries is the whole book, in order. The conditions are read from the
// same flags the story sets, so a page appears exactly when its turn happened.
var diaryEntries = []DiaryEntry{
	{
		ID: "prologue", Order: 0, Title: "ONE WAY DOWN",
		Lines: []string{
			"THE STATION DOES NOT CARE THAT YOU ARRIVED. IT COUNTS.",
			"EIGHTEEN DOCKS, ONE BREAKER, A CLERK WITH A LEDGER - THAT IS THE WHOLE OF HELIOS STATION, AND THE LEDGER IS THE POINT. EVERY HULL THAT DOCKS OWES SOMETHING, AND THE ONE YOU FLEW IN ON OWES MORE THAN MOST.",
			"THE DISPATCHER HAS WORK. IT IS NOT GOOD WORK, BUT IT IS PAID, AND A PILOT WITH AN EMPTY HOLD AND A FULL DEBT DOES NOT ASK WHAT THE WORK IS.",
		},
	},
	{
		ID: "ch1_crate", Order: 10, Title: "THE CRATE",
		When: []Cond{{Key: "HAS_SEAL", Op: "flag"}},
		Lines: []string{
			"THE CRATE IS WARM AND LIGHTER THAN IT SHOULD BE. FOOD RATIONS, THE MANIFEST SAYS. THE SEAL ON THE MANIFEST IS NOT THE CLINIC'S.",
			"I TOOK IT ANYWAY. A HOLD WITH NOTHING IN IT EARNS NOTHING, AND THE DISPATCHER PAYS ON DELIVERY.",
		},
	},
	{
		ID: "ch1_receipt", Order: 11, Title: "THE RECEIPT",
		When: []Cond{{Key: "HAS_RECEIPT", Op: "flag"}},
		Lines: []string{
			"THE CLERK DID NOT WANT TO STAMP IT. WHEN A CLERK DOES NOT WANT TO STAMP SOMETHING, THE PAPER IS WORTH MORE THAN THE CARGO.",
			"I KEPT THE RECEIPT IN THE LOG BOOK, FOLDED TWICE. IT WILL COST ME LATER. IT MAY ALSO SAVE ME.",
		},
	},
	{
		ID: "ch1_refused", Order: 12, Title: "AN EMPTY HOLD",
		When: []Cond{{Key: "REFUSED_SEAL", Op: "flag"}},
		Lines: []string{
			"I DID NOT TAKE THE CRATE. THE DISPATCHER LOOKED AT ME THE WAY A MAN LOOKS AT A BROKEN TOOL AND SAID NOTHING AT ALL.",
			"THE HOLD IS EMPTY AND THE DEBT IS NOT. THAT IS THE WHOLE OF THE DECISION.",
		},
	},
	{
		ID: "ch2_paid", Order: 20, Title: "THE COLLECTORS",
		When: []Cond{{Key: "PAID_PERMIT", Op: "flag"}},
		Lines: []string{
			"THEY CAME OUT OF THE DARK WITH THEIR GUNS COLD AND THEIR PRICE READY. THREE THOUSAND, OR THE HULL SOLD FOR SCRAP.",
			"I PAID. THERE IS NO PAGE IN A LOG BOOK FOR THAT, SO I AM WRITING ONE HERE: I PAID, AND ONE OF THEM WILL REMEMBER IT.",
		},
	},
	{
		ID: "ch2_fought", Order: 21, Title: "THE COLLECTORS, OTHERWISE",
		When: []Cond{{Key: "DISP_OUT", Op: "flag"}},
		Lines: []string{
			"THE COLLECTORS WANTED THE CARGO AND THE HULL BOTH. I GAVE THEM NEITHER - THE RUN WAS NOT DELIVERED, THE DEBT WENT TO THE GUARANTOR, AND I LEFT THE SYSTEM WITH MY NAME ON A LIST.",
			"THE COMPANY WAS POLITE ABOUT IT. POLITE IS HOW A COMPANY WRITES A THREAT.",
		},
	},
	{
		ID: "ch3_burned", Order: 30, Title: "THE RECORDER",
		When: []Cond{{Key: "BURNED_RECORDER", Op: "flag"}},
		Lines: []string{
			"THIRTY SEVEN NAMES ON A CARGO ROSTER AND THREE OF THEM STILL NEEDED A BERTH. THE RECORDER HAD THE REST OF IT: THE LEDGER, THE DATES, THE MEN WHO SIGNED.",
			"I BURNED IT. WHAT IS LEFT IS THE WRECK AND A CREW LIST I CAN STILL SEE WHEN I CLOSE MY EYES.",
		},
	},
	{
		ID: "ch3_copied", Order: 31, Title: "THE RECORDER, COPIED",
		When: []Cond{{Key: "COPIED_RECORDER", Op: "flag"}},
		Lines: []string{
			"THIRTY SEVEN NAMES, THREE VACANCIES, AND A RECORDER THAT KEPT TALKING AFTER THE CREW STOPPED.",
			"I COPIED IT BEFORE I FIRED. THE COPY IS WORTH NOTHING TODAY. TODAY IS NOT THE DAY IT WILL MATTER.",
		},
	},
	{
		ID: "ch4_list", Order: 40, Title: "THE LIST",
		When: []Cond{{Key: "HAS_LIST", Op: "flag"}},
		Lines: []string{
			"THE COURIER HAD A LIST OF EIGHTEEN STATIONS AND WHAT EACH OF THEM OWED. HE SOLD IT FOR SEVEN THOUSAND AND LEFT SMILING, WHICH IS HOW I KNEW THE PRICE WAS TOO LOW.",
			"EVERY NAME ON IT IS A DOOR. SOME OF THEM WILL OPEN.",
		},
	},
	{
		ID: "ch4_read", Order: 41, Title: "THE LIST, READ",
		When: []Cond{{Key: "READ_LIST", Op: "flag"}},
		Lines: []string{
			"I READ IT AND LET HIM GO. HE DID NOT THANK ME; HE TOLD ME WHICH STATION WOULD KILL ME FIRST.",
			"EIGHTEEN STATIONS. THE LIST IS IN MY HEAD NOW, WHICH IS A WORSE PLACE FOR IT THAN A SAFE.",
		},
	},
	{
		ID: "ch4_killed", Order: 42, Title: "THE COURIER'S SHIP",
		When: []Cond{{Key: "COURIER_DEAD", Op: "flag"}},
		Lines: []string{
			"HE RAN, WHICH WAS FAIR, AND HE WAS SLOW, WHICH WAS NOT. THE LIST WAS STILL IN THE WRECK WHEN THE WARNING SHOT WENT THROUGH IT.",
			"THE COMPANY PAID FOR THE COURIER. NOBODY PAID FOR WHAT HE KNEW.",
		},
	},
	{
		ID: "ch5_broke", Order: 50, Title: "THE CONVOY",
		When: []Cond{{Key: "CONVOY_ATTACKED", Op: "flag"}},
		Lines: []string{
			"THREE HAULERS, NO ESCORT WORTH THE NAME, AND CRATES OF LICENCES THAT WOULD HAVE BEEN WORTH FORTY THOUSAND EACH IF THE MARKET HAD STAYED CALM.",
			"IT DID NOT STAY CALM. I BROKE THE CONVOY AND THE PRICE OF EVERYTHING I OWN HALVED BEFORE I DOCKED.",
		},
	},
	{
		ID: "ch5_spared", Order: 51, Title: "THE CONVOY, LET PAST",
		When: []Cond{{Key: "CONVOY_SPARED", Op: "flag"}},
		Lines: []string{
			"THE CAPTAIN TOLD ME WHAT WAS IN THE CRATES AND WHAT BREAKING THEM OPEN WOULD DO TO EVERY PRICE IN THE SYSTEM. THEN HE WAITED FOR ME TO DECIDE.",
			"I LET THEM PASS. THEY WILL NOT THANK ME AND THEY WILL NOT PAY ME, AND THE RECORD WILL SAY SO.",
		},
	},
	{
		ID: "ch5_licence", Order: 52, Title: "A LICENCE",
		When: []Cond{{Key: "BOUGHT_LICENSE", Op: "flag"}},
		Lines: []string{
			"FOUR THOUSAND FOR A LICENCE BOUGHT FROM THE MAN WHO WAS CARRYING TWO HUNDRED OF THEM. IT BUYS FIVE PERCENT OFF AT THIS STATION AND NOTHING ANYWHERE ELSE.",
			"IT IS THE FIRST HONEST THING I HAVE BOUGHT IN THIS SECTOR, AND IT WAS SOLD TO ME BY A SMUGGLER.",
		},
	},
	{
		ID: "ch6_papers", Order: 60, Title: "THE ARCHIVE",
		When: []Cond{{Key: "HAS_COPY", Op: "flag"}},
		Lines: []string{
			"THE CLERK WAS RIGHT: THE CABINET WAS LOCKED AND THE KEY WAS WITH THE CHIEF. THE COMBINE WAS RIGHT TOO: THE COPY CONFIRMS THE CARGO WAS RECALLED, NOT STOLEN, AND CLEARS ME OF NOTHING.",
			"I HAVE A PAGE WITH A CERTIFIED STAMP ON IT. IT IS THE MOST EXPENSIVE USELESS THING I OWN.",
		},
	},
	{
		ID: "ch6_auditor", Order: 61, Title: "THE AUDITOR",
		When: []Cond{{Key: "MET_AUDITOR", Op: "flag"}},
		Lines: []string{
			"THE COMBINE DOES NOT NEGOTIATE. IT AUDITS. SOMEONE, SOMEWHERE, AGREED TO LISTEN TO ME ANYWAY.",
			"AN AUDITOR IS A PERSON WITH A LEDGER AND NO AUTHORITY WHO IS STILL ALLOWED TO ASK QUESTIONS. THAT IS THE FIRST HOPE THIS SECTOR HAS OFFERED ME.",
		},
	},
	{
		ID: "ch6_forged", Order: 62, Title: "A SIGNATURE",
		When: []Cond{{Key: "FORGED", Op: "flag"}},
		Lines: []string{
			"THE CABINET OPENED ON A SIGNATURE THAT WAS NOT THE CLERK'S AND NOT MINE. WHAT IS INSIDE IS NOT WORTH THE RISK YET, BUT IT WILL BE.",
			"THE CLERK WILL NOTICE. THE CLERK WILL SAY NOTHING. THAT IS WHAT A FORGED SIGNATURE IS FOR.",
		},
	},
	{
		ID: "ch7_cut", Order: 70, Title: "THE FIELD",
		When: []Cond{{Key: "SHIELDS_CUT", Op: "flag"}},
		Lines: []string{
			"THE MOTHERSHIP SAT IN A FIELD FED BY THE STATION'S OWN SHIELDS. CUTTING IT MADE THE HIVE NAKED AND THE STATION NAKED WITH IT - TWO HUNDRED PEOPLE AND ONE BREAKER.",
			"THE NAVY WILL HAVE MY NAME FOR THIS. THE ALTERNATIVE WAS TO WAIT TEN MINUTES FOR THE FIELD TO FALL AND FIGHT A MOTHERSHIP THAT HAD ALREADY WOKEN UP.",
		},
	},
	{
		ID: "ch7_wait", Order: 71, Title: "THE FIELD, KEPT",
		When: []Cond{{Key: "SHIELDS_UP", Op: "flag"}},
		Lines: []string{
			"I LEFT THE STATION'S SHIELDS UP AND FLEW AROUND A MOTHERSHIP THAT COULD NOT BE TOUCHED. TEN MINUTES OF WATCHING IT WAIT FOR ME.",
			"WHEN THE FIELD FELL IT FELL ON ITS OWN, AND WHAT CAME OUT OF IT HAD STOPPED BEING PATIENT.",
		},
	},
	{
		ID: "ch7_reform", Order: 72, Title: "A REWRITING",
		When: []Cond{{Key: "AUDIT_REFORMED", Op: "flag"}},
		Lines: []string{
			"THE DOCUMENT WAS ACCEPTED. NOT A CLOSURE - A REWRITING. SEVENTEEN STATIONS GO TO DIRECT CONTRACTS, THE DEBT IS WRITTEN OFF, THE COMBINE LOSES THE LICENCES.",
			"FORTY TWO YEARS OF STAGNATION, CANCELLED BY A PILE OF PAPER SOMEONE ELSE WROTE. I ONLY CARRIED IT HERE.",
		},
	},
	{
		ID: "ending_combine", Order: 80, Title: "THE SEAL, KEPT",
		When: []Cond{{Key: "END_COMBINE", Op: "flag"}},
		Lines: []string{
			"I KEPT THE SEAL. THE COMBINE IS NOT GONE - IT HAS AN OWNER NOW, AND THE OWNER HAS A SHIP AND A HOLD AND A DEBT HE STILL REMEMBERS.",
			"THE SECTOR TRADES ON. IT TRADES BETTER THAN BEFORE, AND IT TRADES BECAUSE OF ME, WHICH IS A SENTENCE I DO NOT SAY OUT LOUD.",
		},
	},
	{
		ID: "ending_council", Order: 81, Title: "THE SEAL, GIVEN",
		When: []Cond{{Key: "END_COUNCIL", Op: "flag"}},
		Lines: []string{
			"THE STATIONS TOOK THE SEAL AND LOCKED IT IN A SHARED SAFE WITH EIGHTEEN KEYS. RULES ARE WRITTEN BY THE PEOPLE WHO LIVE UNDER THEM.",
			"IT IS DULL. IT WORKS. I HAVE NO PART IN IT ANY MORE, AND FOR THE FIRST TIME IN A LONG TIME NOBODY IS LOOKING FOR ME.",
		},
	},
	{
		ID: "ending_chaos", Order: 82, Title: "NO SEAL AT ALL",
		When: []Cond{{Key: "END_CHAOS", Op: "flag"}},
		Lines: []string{
			"THERE ARE NO LICENCES IN THE SECTOR NOW. NO COMBINE, NO ORDER, AND PRICES ON ONE STATION TWICE THOSE ON THE NEXT FOR NO REASON AT ALL.",
			"CARRYING CARGO PAYS AND IS FRIGHTENING. WITHIN A YEAR SOMEONE WILL PROPOSE THAT EVERYONE AGREE ON THE RULES, AND IT WILL NOT BE ME.",
		},
	},
}

// DiaryAvailable is the book as far as this pilot has lived it, in order.
func (g *Game) DiaryAvailable() []DiaryEntry {
	var out []DiaryEntry
	for _, e := range diaryEntries {
		if len(e.When) == 0 || g.conditionsMet(e.When) {
			out = append(out, e)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Order < out[j].Order })
	return out
}

// DiaryUnread is how many pages have not been opened yet.
func (g *Game) DiaryUnread() int {
	n := 0
	for _, e := range g.DiaryAvailable() {
		if !g.diaryRead(e.ID) {
			n++
		}
	}
	return n
}

// DiaryReadIDs is the set of pages already opened, for the screen's marks.
func (g *Game) DiaryReadIDs() map[string]bool {
	out := make(map[string]bool, len(g.Cmd.DiaryRead))
	for _, id := range g.Cmd.DiaryRead {
		out[id] = true
	}
	return out
}

func (g *Game) diaryRead(id string) bool {
	for _, seen := range g.Cmd.DiaryRead {
		if seen == id {
			return true
		}
	}
	return false
}

// DiaryMarkRead remembers that the pilot has opened a page.
func (g *Game) DiaryMarkRead(id string) {
	if g.diaryRead(id) {
		return
	}
	g.Cmd.DiaryRead = append(g.Cmd.DiaryRead, id)
}
