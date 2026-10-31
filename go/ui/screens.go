package ui

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"heliosdrift/galaxy"
	"heliosdrift/game"
	"heliosdrift/input"
	"heliosdrift/lang"
	"heliosdrift/models"
	"heliosdrift/render"
)

func randFloat() float64 { return rand.Float64() }

func clip(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func chrome(fb *render.FB, title string) {
	fb.ClearGradient(0xff0e0704, 0xff0e0704)
	fb.FillRect(0, 0, 319, 13, render.RGB(0x08131f))
	fb.HLine(0, 319, 13, ColFrame)
	fb.Text(title, 6, 4, ColYellow, 1)
	g := game.G
	fb.TextR("GALAXY "+itoa(g.Cmd.Galaxy+1)+"   "+money(g.Cmd.Cash)+" CR", 314, 4, ColText, 1)
}

func footer(fb *render.FB) {
	fb.FillRect(0, 211, 319, 223, render.RGB(0x08131f))
	fb.HLine(0, 319, 211, ColFrame)
	fb.TextC(Hint, 160, 215, ColDim, 1)
}

func docked(fb *render.FB) {
	g := game.G
	title := "STATION"
	if g.World != nil {
		title = g.World.Sys.Name + " STATION"
	}
	chrome(fb, title)
	switch g.Screen {
	case "market":
		marketScreen(fb)
	case "equip":
		equipScreen(fb)
	case "chart":
		chartScreen(fb)
	case "status":
		statusScreen(fb)
	case "mission":
		missionScreen(fb)
	case "contracts":
		contractsScreen(fb)
	case "journal":
		journalScreen(fb)
	case "shipyard":
		shipyardScreen(fb)
	case "slots":
		slotsScreen(fb)
	default:
		stationMenu(fb)
	}
	footer(fb)
}

func stationMenu(fb *render.FB) {
	g := game.G
	cmd := g.Cmd
	Hint = "1-8 SELECT   UP/DOWN   F REFUEL   ENTER CONFIRM"

	// left: commander sheet
	plate(fb, 8, 16, 152, 188, "COMMANDER "+cmd.Name)
	row := func(label, value string, y int, col uint32) {
		fb.Text(label, 16, y, ColDim, 1)
		fb.TextR(value, 150, y, col, 1)
	}
	row("CASH", money(cmd.Cash)+" CR", 34, ColYellow)
	row("FUEL", strconv.FormatFloat(cmd.Fuel, 'f', 1, 64)+" LY", 50, ColText)
	row("HOLD", itoa(game.CargoUsed(cmd))+"/"+itoa(cmd.CargoCap)+" T", 66, ColText)
	row("MISSILES", itoa(cmd.Missiles), 82, ColText)
	row("LASER", game.LaserName[cmd.Laser], 98, ColText)
	row("RANK", galaxy.RankFor(cmd.Kills), 114, ColGreen)
	legal := "CLEAN"
	lc := ColGreen
	if cmd.Legal >= 500 {
		legal, lc = "FUGITIVE", ColRed
	} else if cmd.Legal > 0 {
		legal, lc = "OFFENDER", ColRed
	}
	row("LEGAL", legal, 130, lc)
	gal := galaxy.Get(cmd.Galaxy)
	row("TARGET", gal.Systems[((cmd.TargetSystem%256)+256)%256].Name, 146, ColText)

	// system block at the bottom of the same plate
	fb.HLine(16, 150, 160, render.RGB(0x1b3348))
	if g.World != nil {
		sys := g.World.Sys
		fb.Text(sys.Name, 16, 168, ColYellow, 1)
		fb.Text(galaxy.EconomyNames[sys.Economy], 16, 178, ColDim, 1)
		fb.Text(galaxy.GovNames[sys.Government], 16, 188, ColDim, 1)
	}

	// right: services
	plate(fb, 168, 16, 144, 188, "STATION SERVICES")
	items := menuItems()
	step := 16
	switch {
	case len(items) > 10:
		step = 11
	case len(items) > 9:
		step = 12
	case len(items) > 8:
		step = 14
	case len(items) > 7:
		step = 15
	}
	for i, it := range items {
		y := 34 + i*step
		sel := S.Menu == i
		if sel {
			fb.FillRect(174, y-4, 306, y+9, render.RGB(0x12384f))
			fb.FillRect(174, y-4, 176, y+9, ColYellow)
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.Text(itoa(i+1)+"  "+it, 180, y, col, 1)
	}
	fb.HLine(174, 306, 158, render.RGB(0x1b3348))
	fb.Text("REFUEL  [F]", 180, 168, ColDim, 1)
	fb.TextR(itoa(galaxy.FuelPricePerLY)+" CR/LY", 300, 168, ColGreen, 1)
	if g.World != nil {
		fb.Text("TECH "+itoa(g.World.Sys.Tech), 180, 186, ColDim, 1)
		fb.TextR("DOCKED", 300, 186, render.RGB(0x4c6a80), 1)
	}
}

// menuItems is the station service list, in hotkey order.
func menuItems() []string {
	items := []string{"BUY / SELL CARGO", "EQUIP SHIP"}
	if !game.G.Classic() {
		items = append(items, "SHIPYARD")
	}
	items = append(items, "GALACTIC CHART",
		"COMMANDER STATUS", "LAUNCH", "SAVE COMMANDER", "SETTINGS")
	if game.G.MissionOffer() != "" || game.G.Cmd.Mission.Kind != "" ||
		game.G.CampaignOffer() || game.G.Cmd.Campaign.Active {
		items = append(items, "MISSION")
	}
	if !game.G.Classic() && (len(game.G.Offers) > 0 || len(game.G.Cmd.Contracts) > 0) {
		items = append(items, "CONTRACTS")
	}
	if game.G.Cmd.Mission.Kind != "" || len(game.G.Cmd.Contracts) > 0 {
		items = append(items, "JOURNAL")
	}
	return items
}

func openMission() {
	g := game.G
	if g.MissionOffer() == "" && g.Cmd.Mission.Kind == "" && !g.CampaignOffer() && !g.Cmd.Campaign.Active {
		g.Message("NO CONTRACT ON OFFER", "bad")
		return
	}
	g.Screen = "mission"
}

func missionScreen(fb *render.FB) {
	g := game.G
	cmd := &g.Cmd
	plate(fb, 24, 30, 272, 150, "MISSION")
	if !g.Classic() {
		q := g.CampaignQuest()
		active := g.Cmd.Campaign.Active
		if q.Title == "" && !active {
			fb.TextC("CAMPAIGN COMPLETE", 160, 62, ColGreen, 1)
			Hint = "ESC BACK"
			return
		}
		fb.TextC(q.Title, 160, 44, ColYellow, 1)
		y := 60
		for _, ln := range q.Brief {
			if ln == "" {
				y += 4
				continue
			}
			fb.TextC(clip(ln, 42), 160, y, ColText, 1)
			y += 10
		}
		if active {
			if st := g.CampaignStatus(); st != "" {
				fb.TextC(clip(st, 40), 160, 148, ColGreen, 1)
			}
			Hint = "ESC BACK"
		} else {
			fb.TextC("ENTER ACCEPT     ESC DECLINE", 160, 148, ColYellow, 1)
			Hint = "ENTER ACCEPT   ESC DECLINE"
		}
		if q.Pay > 0 {
			fb.TextC("FEE: "+itoa(q.Pay)+" CR", 160, 128, ColDim, 1)
		}
		return
	}
	kind := g.MissionOffer()
	active := cmd.Mission.Kind != ""
	if active {
		kind = cmd.Mission.Kind
		fb.TextC("CONTRACT IN PROGRESS", 160, 44, ColDim, 1)
	} else {
		fb.TextC("CONTRACT OFFER", 160, 44, ColDim, 1)
	}
	y := 62
	for _, ln := range g.MissionBrief(kind) {
		if ln == "" {
			y += 5
			continue
		}
		fb.TextC(ln, 160, y, ColText, 1)
		y += 11
	}
	if active {
		Hint = "ESC BACK"
		if kind == game.MissionDocuments {
			fb.TextC("DOCK AT THE DESTINATION TO DELIVER", 160, 158, ColGreen, 1)
		} else {
			fb.TextC("JUMP TO THE SYSTEM AND DESTROY IT", 160, 158, ColGreen, 1)
		}
		return
	}
	Hint = "ENTER ACCEPT   ESC DECLINE"
	fb.TextC("ENTER ACCEPT     ESC DECLINE", 160, 158, ColYellow, 1)
}

func marketScreen(fb *render.FB) {
	g := game.G
	cmd := g.Cmd
	Hint = "UP/DN  B BUY  S SELL  SHIFT+ 10  ESC"
	plate(fb, 4, 18, 312, 190, "GALACTIC MARKET")
	fb.Text("COMMODITY", 12, 30, ColDim, 1)
	fb.TextR("SELL", 178, 30, ColDim, 1)
	fb.TextR("BUY", 220, 30, ColDim, 1)
	fb.TextR("QTY", 254, 30, ColDim, 1)
	fb.TextR("HOLD", 294, 30, ColDim, 1)
	fb.HLine(10, 310, 37, ColFrame)
	for i, m := range g.Market {
		y := 42 + i*9
		sel := S.Market == i
		if sel {
			fb.FillRect(8, y-1, 311, y+7, render.RGB(0x12384f))
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.Text(m.Name, 12, y, col, 1)
		if m.Unit == "KG" {
			fb.Text("KG", 104, y, ColDim, 1)
		}
		sc := ColDim
		if m.Sell > m.Price {
			sc = ColGreen
		}
		fb.TextR(itoa(m.Sell), 178, y, sc, 1)
		pc := ColGreen
		if sel {
			pc = ColWhite
		}
		fb.TextR(itoa(m.Price), 220, y, pc, 1)
		fb.TextR(itoa(m.Qty), 254, y, ColDim, 1)
		have := cmd.Cargo[i]
		hc := ColDim
		if have > 0 {
			hc = ColYellow
		}
		fb.TextR(itoa(have), 294, y, hc, 1)
	}
	fb.Text("HOLD "+itoa(game.CargoUsed(cmd))+"/"+itoa(cmd.CargoCap)+" T", 240, 197, ColYellow, 1)
}

func equipScreen(fb *render.FB) {
	g := game.G
	cmd := g.Cmd
	Hint = "UP/DOWN SELECT   ENTER BUY   ESC BACK"
	tech := 0
	if g.World != nil {
		tech = g.World.Sys.Tech
	}
	plate(fb, 4, 18, 312, 190, "EQUIP SHIP  -  TECH LEVEL "+itoa(tech))
	items := galaxy.AvailableEquipment(g.World.Sys)
	for i, it := range items {
		y := 32 + i*12
		sel := S.Equip == i
		if sel {
			fb.FillRect(8, y-2, 311, y+8, render.RGB(0x12384f))
		}
		owned := false
		if it.ID == "missile" {
			owned = cmd.Missiles >= game.MissilesMax
		} else if len(it.ID) >= 5 && it.ID[:5] == "laser" {
			owned = false
		} else {
			owned = cmd.Equip[it.ID]
		}
		col := ColText
		if owned {
			col = ColDim
		} else if sel {
			col = ColYellow
		}
		fb.Text(it.Name, 14, y, col, 1)
		pc := ColRed
		if cmd.Cash >= float64(it.Price) {
			pc = ColGreen
		}
		fb.TextR(itoa(it.Price)+" CR", 240, y, pc, 1)
		if it.ID == "missile" {
			fb.TextR("x"+itoa(cmd.Missiles), 290, y, ColText, 1)
		} else if owned {
			fb.TextR("FITTED", 290, y, ColGreen, 1)
		}
	}
}

func chartScreen(fb *render.FB) {
	g := game.G
	cmd := g.Cmd
	Hint = "ARROWS SELECT SYSTEM   H HYPERSPACE   ESC BACK"
	gal := galaxy.Get(cmd.Galaxy)
	cur := gal.Systems[cmd.System]
	selIdx := S.ChartCursor
	if selIdx < 0 || selIdx >= len(gal.Systems) {
		selIdx = cur.ID
	}
	sel := gal.Systems[selIdx]

	plate(fb, 4, 18, 196, 190, "GALACTIC CHART "+itoa(cmd.Galaxy+1))
	ox, oy, cw, ch := 14, 34, 176, 164
	fb.Rect(ox, oy, ox+cw, oy+ch, render.RGB(0x16394a))
	cxp := float64(ox) + float64(cur.X)/256*float64(cw)
	cyp := float64(oy+ch) - float64(cur.Y)/256*float64(ch)
	rr := int(28.0 / 256 * float64(cw))
	fb.Ellipse(int(cxp), int(cyp), rr, rr, render.RGB(0x3fbcd8))
	fb.Ellipse(int(cxp), int(cyp), rr-1, rr-1, render.RGB(0x1f5e6e))

	for _, s := range gal.Systems {
		x := ox + s.X*cw/256
		y := oy + ch - s.Y*ch/256
		switch {
		case s.ID == cur.ID:
			fb.FillRect(x-1, y-1, x+1, y+1, ColWhite)
		case s.ID == sel.ID:
			fb.FillRect(x-1, y-1, x+1, y+1, ColYellow)
			fb.Rect(x-3, y-3, x+3, y+3, ColOrange)
		default:
			fb.Px(x, y, render.Color("teal", 3))
		}
	}

	// mission and contract destinations, so a job is findable on the chart
	if sys, gal, ok := g.MissionTarget(); ok && gal == cmd.Galaxy {
		s := galSystems(gal)[sys]
		mx := ox + s.X*cw/256
		my := oy + ch - s.Y*ch/256
		r := 3 + int(2*math.Abs(math.Sin(g.Time*2.2)))
		fb.Rect(mx-r, my-r, mx+r, my+r, render.Color("pink", 4))
	}
	for _, id := range g.ContractTargets() {
		s := galSystems(cmd.Galaxy)[id]
		mx := ox + s.X*cw/256
		my := oy + ch - s.Y*ch/256
		fb.FillRect(mx-1, my-1, mx+1, my+1, render.Color("green", 4))
	}

	plate(fb, 200, 18, 118, 190, "SYSTEM DATA")
	ty := 32
	line := func(label, value string, col uint32) {
		fb.Text(label, 206, ty, ColDim, 1)
		ty += 8
		fb.Text(clip(value, 18), 206, ty, col, 1)
		ty += 11
	}
	fb.Text(sel.Name, 206, ty, ColYellow, 1)
	ty += 11
	line("ECONOMY", galaxy.EconomyNames[sel.Economy], ColText)
	line("GOVERNMENT", galaxy.GovNames[sel.Government], ColText)
	fb.Text("TECH LEVEL "+itoa(sel.Tech), 206, ty, ColText, 1)
	ty += 9
	fb.Text("POP "+strconv.FormatFloat(sel.Pop, 'f', 1, 64)+" B", 206, ty, ColText, 1)
	ty += 9
	fb.Text("PROD "+itoa(sel.Prod)+" MC", 206, ty, ColText, 1)
	ty += 9
	fb.Text("RADIUS "+itoa(sel.Radius)+" KM", 206, ty, ColText, 1)
	ty += 11
	d := galaxy.Distance(cur, sel)
	dc := ColRed
	if d <= 7 {
		dc = ColGreen
	}
	fb.Text("DIST "+strconv.FormatFloat(d, 'f', 1, 64)+" LY", 206, ty, dc, 1)
	ty += 9
	fb.Text("FUEL "+strconv.FormatFloat(cmd.Fuel, 'f', 1, 64)+" LY", 206, ty, ColText, 1)
	ty += 11
	if sel.ID != cur.ID {
		if d <= 7 && cmd.Fuel >= d {
			fb.Text("H - HYPERSPACE", 206, ty, ColYellow, 1)
		} else {
			fb.Text("OUT OF RANGE", 206, ty, ColDim, 1)
		}
	}
}

func statusScreen(fb *render.FB) {
	g := game.G
	cmd := g.Cmd
	Hint = "ESC BACK"
	plate(fb, 4, 18, 152, 190, "COMMANDER")
	rows := [][2]string{
		{"NAME", cmd.Name},
		{"SHIP", g.Ship().Name},
		{"MODE", g.ModeName()},
		{"CASH", money(cmd.Cash) + " CR"},
		{"FUEL", strconv.FormatFloat(cmd.Fuel, 'f', 1, 64) + " LY"},
		{"CARGO", itoa(game.CargoUsed(cmd)) + "/" + itoa(cmd.CargoCap) + " T"},
		{"MISSILES", itoa(cmd.Missiles)},
		{"LASER", game.LaserName[cmd.Laser]},
		{"RANK", galaxy.RankFor(cmd.Kills)},
		{"KILLS", itoa(cmd.Kills)},
	}
	legal := "CLEAN"
	if cmd.Legal >= 500 {
		legal = "FUGITIVE"
	} else if cmd.Legal > 0 {
		legal = "OFFENDER"
	}
	rows = append(rows, [2]string{"LEGAL", legal})
	rows = append(rows, [2]string{"GALAXY", itoa(cmd.Galaxy + 1)})
	rows = append(rows, [2]string{"DAY", itoa(cmd.Day)})
	if cmd.Pax > 0 || cmd.Mail > 0 {
		rows = append(rows, [2]string{"PAX/MAIL", itoa(cmd.Pax) + "/" + itoa(cmd.Mail)})
	}
	y := 32
	for _, r := range rows {
		fb.Text(r[0], 12, y, ColDim, 1)
		col := ColText
		if r[0] == "MODE" && g.Classic() {
			col = ColYellow
		}
		fb.TextR(r[1], 150, y, col, 1)
		y += 11
	}

	if ms := g.MissionStatus(); ms != "" {
		fb.HLine(12, 148, 176, render.RGB(0x1b3348))
		fb.Text("MISSION", 12, 182, ColDim, 1)
		fb.Text(clip(ms, 23), 12, 192, ColGreen, 1)
	}

	plate(fb, 160, 18, 156, 190, "EQUIPMENT")
	eq := [][2]interface{}{
		{"E.C.M.", cmd.Equip["ecm"]},
		{"FUEL SCOOP", cmd.Equip["fuelscoop"]},
		{"LARGE BAY", cmd.Equip["largebay"]},
		{"ESCAPE CAPSULE", cmd.Equip["escape"]},
		{"ENERGY UNIT", cmd.Equip["extraenergy"]},
		{"DOCKING COMP", cmd.Equip["dockcomp"]},
		{"GAL. HYPERDRIVE", cmd.Equip["galdrive"]},
		{"ENERGY BOMB", cmd.Bombs > 0},
	}
	yy := 32
	for _, e := range eq {
		on := e[1].(bool)
		col := render.RGB(0x2d4a5c)
		if on {
			col = ColGreen
		}
		fb.Text(e[0].(string), 168, yy, col, 1)
		v := "-"
		if on {
			v = "YES"
		}
		fb.TextR(v, 308, yy, col, 1)
		yy += 11
	}
	yy += 6
	fb.Text("CARGO:", 168, yy, ColDim, 1)
	yy += 10
	any := false
	more := 0
	for i := 0; i < len(galaxy.Commodities); i++ {
		if cmd.Cargo[i] <= 0 {
			continue
		}
		if yy > 184 {
			more++
			continue
		}
		any = true
		fb.Text(galaxy.Commodities[i].Name+" x"+itoa(cmd.Cargo[i]), 168, yy, ColText, 1)
		yy += 9
	}
	if !any {
		fb.Text("EMPTY", 168, yy, ColDim, 1)
	} else if more > 0 {
		fb.Text("+"+itoa(more)+" MORE", 168, yy, ColDim, 1)
	}
}

func pauseScreen(fb *render.FB) {
	chrome(fb, "PAUSED")
	plate(fb, 40, 26, 240, 170, "FLIGHT CONTROLS")
	steer := "PITCH AND ROLL"
	if !game.G.Classic() {
		steer = "PITCH / ROLL / YAW"
	}
	rows := [][2]string{
		{"ARROWS / MOUSE", steer},
		{"X / Z", "ACCELERATE / BRAKE"},
		{"SPACE or LMB", "FIRE LASERS"},
		{"M", "LAUNCH MISSILE"},
		{"E", "E.C.M."},
		{"TAB", "ENERGY BOMB"},
		{"L", "DOCKING COMPUTER"},
		{"J", "IN-SYSTEM JUMP"},
		{"C", "GALACTIC CHART"},
		{"V", "COMMANDER STATUS"},
		{"ESC", "RESUME"},
		{"Q", "MAIN MENU"},
	}
	for i, r := range rows {
		fb.Text(r[0], 52, 40+i*12, ColYellow, 1)
		fb.TextR(r[1], 268, 40+i*12, ColText, 1)
	}
}

// Key handles a key press; returns true when consumed.
func Key(k string) bool {
	g := game.G
	// the settings screen sits on top of whatever state opened it, so it gets
	// first refusal on every key
	if g.Screen == "settings" {
		return settingsKey(k)
	}
	if g.Screen == "controls" {
		return controlsKey(k)
	}
	if g.Screen == "padtest" {
		if k == "Escape" || k == "Enter" {
			game.G.Screen = S.PadTestFrom
		}
		return true
	}
	if g.Screen == "slots" && g.State == "title" {
		return slotsKey(k)
	}
	if g.Screen == "mode" {
		return modeKey(k)
	}
	switch g.State {
	case "title":
		items := TitleItems()
		if S.Title >= len(items) {
			S.Title = 0
		}
		switch k {
		case "ArrowUp":
			S.Title = (S.Title + len(items) - 1) % len(items)
			return true
		case "ArrowDown":
			S.Title = (S.Title + 1) % len(items)
			return true
		case "Enter":
			switch items[S.Title] {
			case "NEW GAME":
				openMode()
			case "CONTINUE FLIGHT":
				ContinueFlight()
			case "LOAD COMMANDER":
				if game.SaveExists() {
					openSlots(false)
				} else {
					g.Message("NO SAVED COMMANDER", "bad")
				}
			case "SETTINGS":
				OpenSettings()
			case "QUIT":
				quitRequested = true
			}
			return true
		}
		// digits pick an entry, matching the number shown in the hint
		if n, err := strconv.Atoi(k); err == nil && n >= 1 && n <= len(items) {
			S.Title = n - 1
		}
		return true
	case "dead":
		if k == "Enter" {
			DropSession()
			g.State = "title"
			g.Screen = "launch"
		}
		return true
	case "docked":
		return dockedKey(k)
	case "flight":
		return flightKey(k)
	}
	return false
}

func dockedKey(k string) bool {
	g := game.G
	if g.Screen == "launch" {
		switch k {
		case "ArrowUp":
			n := len(menuItems())
			S.Menu = (S.Menu + n - 1) % n
			return true
		case "ArrowDown":
			n := len(menuItems())
			S.Menu = (S.Menu + 1) % n
			return true
		case "f":
			g.Refuel()
			return true
		case "1":
			g.Screen = "market"
			return true
		case "2":
			g.Screen = "equip"
			return true
		case "3":
			if g.Classic() {
				g.Screen = "chart"
			} else {
				g.Screen = "shipyard"
				S.Ship = 0
			}
			return true
		case "4":
			g.Screen = "chart"
			S.ChartCursor = g.Cmd.System
			return true
		case "5":
			g.Screen = "status"
			return true
		case "6":
			g.Launch()
			return true
		case "7":
			openSlots(true)
			return true
		case "8":
			OpenSettings()
			return true
		case "9":
			openMission()
			return true
		case "0":
			openContracts()
			return true
		case "Enter":
			it := menuItems()
			if S.Menu < 0 || S.Menu >= len(it) {
				return true
			}
			switch it[S.Menu] {
			case "BUY / SELL CARGO":
				g.Screen = "market"
			case "EQUIP SHIP":
				g.Screen = "equip"
			case "GALACTIC CHART":
				g.Screen = "chart"
				S.ChartCursor = g.Cmd.System
			case "COMMANDER STATUS":
				g.Screen = "status"
			case "SHIPYARD":
				g.Screen = "shipyard"
				S.Ship = 0
			case "LAUNCH":
				g.Launch()
			case "SAVE COMMANDER":
				openSlots(true)
			case "SETTINGS":
				OpenSettings()
			case "MISSION":
				openMission()
			case "CONTRACTS":
				openContracts()
			case "JOURNAL":
				openJournal()
			}
			return true
		}
		return true
	}
	if g.Screen == "journal" {
		if k == "Escape" {
			g.Screen = "launch"
		}
		return true
	}
	if g.Screen == "contracts" {
		return contractsKey(k)
	}
	if g.Screen == "shipyard" {
		return shipyardKey(k)
	}
	if g.Screen == "slots" {
		return slotsKey(k)
	}
	if g.Screen == "mission" {
		switch k {
		case "Enter":
			if !g.Classic() {
				g.AcceptCampaign()
			} else if kind := g.MissionOffer(); kind != "" {
				g.AcceptMission(kind)
			}
			g.Screen = "launch"
			return true
		case "Escape":
			g.Screen = "launch"
			return true
		}
		return true
	}
	if k == "Escape" {
		g.Screen = "launch"
		return true
	}
	switch g.Screen {
	case "market":
		switch k {
		case "ArrowUp":
			S.Market = (S.Market + 16) % 17
		case "ArrowDown":
			S.Market = (S.Market + 1) % 17
		case "b":
			g.Buy(S.Market, 1)
		case "B":
			g.Buy(S.Market, 10)
		case "s":
			g.Sell(S.Market, 1)
		case "S":
			g.Sell(S.Market, 10)
		}
		return true
	case "equip":
		items := galaxy.AvailableEquipment(g.World.Sys)
		if len(items) == 0 {
			return true
		}
		switch k {
		case "ArrowUp":
			S.Equip = (S.Equip + len(items) - 1) % len(items)
		case "ArrowDown":
			S.Equip = (S.Equip + 1) % len(items)
		case "Enter":
			if S.Equip < len(items) {
				g.BuyEquipment(items[S.Equip].ID)
			}
		}
		return true
	case "chart":
		gal := galaxy.Get(g.Cmd.Galaxy)
		switch k {
		case "ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight":
			S.ChartCursor = nearestInDirection(gal, chartCursorID(), k)
		case "h":
			if g.Hyperspace(chartCursorID()) {
				g.Overlay = ""
			}
		}
		return true
	}
	return true
}

func chartCursorID() int {
	if S.ChartCursor < 0 {
		return game.G.Cmd.System
	}
	return S.ChartCursor
}

func flightKey(k string) bool {
	g := game.G
	if g.Overlay != "" {
		return overlayKey(k)
	}
	switch k {
	case "c":
		g.Overlay = "chart"
		S.ChartCursor = g.Cmd.System
		return true
	case "v":
		g.Overlay = "status"
		return true
	case "Escape":
		g.Overlay = "pause"
		return true
	case "m":
		g.FireMissile()
		return true
	case "e":
		g.FireECM()
		return true
	case "Tab":
		g.FireBomb()
		return true
	case "l":
		g.ToggleAutopilot()
		return true
	case "j":
		g.JumpInSystem()
		return true
	case "g":
		g.GalacticJump()
		return true
	case "k":
		g.Jettison()
		return true
	}
	return false
}

func overlayKey(k string) bool {
	g := game.G
	if k == "Escape" {
		g.Overlay = ""
		return true
	}
	if k == "q" && g.Overlay == "pause" {
		QuitToTitle()
		return true
	}
	if g.Overlay == "chart" {
		gal := galaxy.Get(g.Cmd.Galaxy)
		switch k {
		case "ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight":
			S.ChartCursor = nearestInDirection(gal, chartCursorID(), k)
			return true
		case "h":
			if g.Hyperspace(chartCursorID()) {
				g.Overlay = ""
			}
			return true
		case "c":
			g.Overlay = ""
			return true
		}
		return true
	}
	if k == "v" {
		g.Overlay = ""
	}
	return true
}

func nearestInDirection(gal *galaxy.Galaxy, fromID int, key string) int {
	if fromID < 0 || fromID >= len(gal.Systems) {
		fromID = 0
	}
	from := gal.Systems[fromID]
	dx, dy := 0, 0
	switch key {
	case "ArrowUp":
		dy = 1
	case "ArrowDown":
		dy = -1
	case "ArrowLeft":
		dx = -1
	case "ArrowRight":
		dx = 1
	default:
		return fromID
	}
	best := -1
	bestScore := math.Inf(1)
	for _, s := range gal.Systems {
		if s.ID == fromID {
			continue
		}
		ddx := float64(s.X - from.X)
		ddy := float64(s.Y - from.Y)
		if ddx > 128 {
			ddx -= 256
		}
		if ddx < -128 {
			ddx += 256
		}
		if ddy > 128 {
			ddy -= 256
		}
		if ddy < -128 {
			ddy += 256
		}
		along := ddx*float64(dx) + ddy*float64(dy)
		if along <= 2 {
			continue
		}
		perp := math.Abs(ddx*float64(dy) - ddy*float64(dx))
		score := along + perp*2.5
		if score < bestScore {
			bestScore = score
			best = s.ID
		}
	}
	if best < 0 {
		return fromID
	}
	return best
}

var _ = models.StationSize

// ---------------------------------------------------------------- settings

// OpenSettings enters the settings screen from wherever the player is.
func OpenSettings() {
	g := game.G
	S.PrevState = g.State
	S.PrevScreen = g.Screen
	S.Settings = 0
	g.Screen = "settings"
}

func settingsScreen(fb *render.FB) {
	plate(fb, 40, 20, 240, 158, "")
	rows := []string{"LANGUAGE", "SOUND", "VOLUME", "CONTROLS", "BACK"}
	for i, name := range rows {
		y := 42 + i*28
		sel := S.Settings == i
		if sel {
			fb.FillRect(48, y-8, 272, y+11, render.RGB(0x12384f))
			fb.FillRect(48, y-8, 51, y+11, ColYellow)
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.Text(name, 60, y, col, 1)
		switch i {
		case 0:
			val := "ENGLISH"
			if lang.IsRU() {
				val = "RUSSIAN"
			}
			fb.TextR(val, 264, y, ColGreen, 1)
		case 1:
			val := "ON"
			if lang.Cfg.Muted {
				val = "OFF"
			}
			fb.TextR(val, 264, y, ColGreen, 1)
		case 2:
			pips := int(lang.Cfg.Volume*5 + 0.5)
			for k := 0; k < 5; k++ {
				c := render.RGB(0x25405a)
				if k < pips {
					c = ColGreen
				}
				fb.FillRect(224+k*9, y-1, 230+k*9, y+6, c)
			}
		case 3:
			fb.TextR(">", 264, y, ColGreen, 1)
		}
	}
	fb.TextC("SHIFT+M - SOUND ON/OFF", 160, 168, ColDim, 1)
	Hint = "UP/DOWN  LEFT/RIGHT  ENTER  ESC"
}

func settingsKey(k string) bool {
	switch k {
	case "ArrowUp":
		S.Settings = (S.Settings + 4) % 5
	case "ArrowDown":
		S.Settings = (S.Settings + 1) % 5
	case "Escape":
		closeSettings()
	case "Enter", "ArrowLeft", "ArrowRight":
		switch S.Settings {
		case 0:
			if lang.IsRU() {
				lang.Set("en")
			} else {
				lang.Set("ru")
			}
			lang.Save()
		case 1:
			lang.Cfg.Muted = !lang.Cfg.Muted
			lang.Save()
			state := "ON"
			if lang.Cfg.Muted {
				state = "OFF"
			}
			game.G.Message(lang.T("SOUND")+" "+lang.T(state), "info")
		case 2:
			steps := []float64{0.2, 0.4, 0.6, 0.8, 1.0}
			idx := 3
			for i, v := range steps {
				if math.Abs(v-lang.Cfg.Volume) < 0.01 {
					idx = i
				}
			}
			if k == "ArrowLeft" {
				idx = (idx + len(steps) - 1) % len(steps)
			} else {
				idx = (idx + 1) % len(steps)
			}
			lang.Cfg.Volume = steps[idx]
			lang.Save()
		case 3:
			openControls()
		case 4:
			closeSettings()
		}
	}
	return true
}

func closeSettings() {
	g := game.G
	g.State = S.PrevState
	g.Screen = S.PrevScreen
}

// ------------------------------------------------------------------ contracts

func openContracts() {
	g := game.G
	if len(g.Offers) == 0 && len(g.Cmd.Contracts) == 0 {
		g.Message("NO CONTRACTS ON THE BOARD", "bad")
		return
	}
	g.Screen = "contracts"
	S.Contract = 0
}

func contractsScreen(fb *render.FB) {
	g := game.G
	cmd := &g.Cmd
	plate(fb, 8, 20, 304, 176, "CONTRACTS")
	y := 36
	fb.Text("AVAILABLE", 16, y, ColDim, 1)
	y += 12
	if len(g.Offers) == 0 {
		fb.Text("NONE RIGHT NOW", 16, y, render.RGB(0x2d4a5c), 1)
		y += 12
	}
	for i, c := range g.Offers {
		sel := S.Contract == i
		if sel {
			fb.FillRect(12, y-2, 300, y+8, render.RGB(0x12384f))
			fb.FillRect(12, y-2, 15, y+8, ColYellow)
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.Text(clip(g.ContractLabel(c), 47), 18, y, col, 1)
		y += 11
	}
	y += 4
	fb.HLine(16, 300, y, render.RGB(0x1b3348))
	y += 8
	fb.Text("ACCEPTED", 16, y, ColDim, 1)
	y += 12
	if len(cmd.Contracts) == 0 {
		fb.Text("NONE", 16, y, render.RGB(0x2d4a5c), 1)
	} else {
		for _, c := range cmd.Contracts {
			fb.Text(clip(g.ContractLabel(c), 47), 16, y, ColGreen, 1)
			y += 10
		}
	}
	if len(g.Offers) > 0 {
		Hint = "UP/DOWN SELECT   ENTER ACCEPT   ESC BACK"
	} else {
		Hint = "ESC BACK"
	}
}

func contractsKey(k string) bool {
	g := game.G
	n := len(g.Offers)
	switch k {
	case "Escape":
		g.Screen = "launch"
		return true
	case "ArrowUp":
		if n > 0 {
			S.Contract = (S.Contract + n - 1) % n
		}
		return true
	case "ArrowDown":
		if n > 0 {
			S.Contract = (S.Contract + 1) % n
		}
		return true
	case "Enter":
		g.AcceptOffer(S.Contract)
		if len(g.Offers) == 0 {
			S.Contract = 0
		}
		return true
	}
	return true
}

// ------------------------------------------------------------------- shipyard

func shipyardScreen(fb *render.FB) {
	g := game.G
	cur := g.Ship()
	plate(fb, 6, 16, 308, 180, "SHIPYARD")
	fb.Text("FLYING "+cur.Name, 14, 30, ColYellow, 1)
	fb.TextR("TRADE IN "+itoa(g.TradeInValue())+" CR", 306, 30, ColGreen, 1)
	for i, sc := range game.ShipYard {
		y := 46 + i*28
		sel := S.Ship == i
		if sel {
			fb.FillRect(10, y-4, 310, y+18, render.RGB(0x12384f))
			fb.FillRect(10, y-4, 13, y+18, ColYellow)
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.Text(sc.Name, 18, y, col, 1)
		fb.Text("HOLD "+itoa(sc.Cargo)+"T   SPD "+itoa(int(sc.Speed))+"   EN "+itoa(int(sc.Energy)), 18, y+9, ColDim, 1)
		fb.Text("SHIELD "+itoa(int(sc.Shield))+"   TURN "+two(sc.Pitch)+"   LASER "+laserCap(sc.MaxLaser), 18, y+17, ColDim, 1)
		if sc.ID == cur.ID {
			fb.TextR("FITTED", 302, y, ColGreen, 1)
		} else {
			net := sc.Price - g.TradeInValue()
			c := ColGreen
			if float64(net) > g.Cmd.Cash {
				c = ColRed
			}
			fb.TextR(itoa(net)+" CR", 302, y, c, 1)
		}
	}
	Hint = "UP/DOWN SELECT   ENTER BUY   ESC BACK"
}

// two formats a multiplier with two decimals, e.g. 1.45.
func two(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

// laserCap names the best gun a hull accepts.
func laserCap(level int) string {
	switch {
	case level >= 2:
		return "MILITARY"
	case level == 1:
		return "BEAM"
	}
	return "PULSE"
}

func shipyardKey(k string) bool {
	g := game.G
	n := len(game.ShipYard)
	switch k {
	case "Escape":
		g.Screen = "launch"
		return true
	case "ArrowUp":
		S.Ship = (S.Ship + n - 1) % n
		return true
	case "ArrowDown":
		S.Ship = (S.Ship + 1) % n
		return true
	case "Enter":
		if S.Ship >= 0 && S.Ship < n {
			g.BuyShip(game.ShipYard[S.Ship].ID)
		}
		return true
	}
	return true
}

// --------------------------------------------------------------------- slots

// openSlots shows the three save files, in save or load mode.
// OpenSlotsForShots exposes the slot screen to the screenshot pass.
func OpenSlotsForShots(save bool) { openSlots(save) }

func openSlots(save bool) {
	S.SlotSave = save
	S.Slot = 0
	game.G.Screen = "slots"
}

func slotsScreen(fb *render.FB) {
	g := game.G
	title := "LOAD COMMANDER"
	if S.SlotSave {
		title = "SAVE COMMANDER"
	}
	plate(fb, 34, 40, 252, 120, title)
	for i := 0; i < game.Slots; i++ {
		y := 62 + i*26
		sel := S.Slot == i
		if sel {
			fb.FillRect(40, y-5, 280, y+14, render.RGB(0x12384f))
			fb.FillRect(40, y-5, 43, y+14, ColYellow)
		}
		name, cash, where, ok := game.SlotInfo(i)
		col := ColDim
		label := itoa(i+1) + "  EMPTY"
		if ok {
			col = ColText
			label = itoa(i+1) + "  " + name + "  " + money(cash) + " CR  " + where
		}
		if sel {
			col = ColYellow
		}
		fb.Text(clip(label, 38), 46, y, col, 1)
	}
	if S.SlotSave {
		Hint = "UP/DOWN   ENTER SAVE   ESC BACK"
	} else {
		Hint = "UP/DOWN   ENTER LOAD   ESC BACK"
		if g.State == "title" {
			Hint = "UP/DOWN   ENTER LOAD   ESC BACK"
		}
	}
}

func slotsKey(k string) bool {
	g := game.G
	switch k {
	case "Escape":
		if g.State == "title" {
			g.Screen = "launch"
			S.Title = 0
		} else {
			g.Screen = "launch"
		}
		return true
	case "ArrowUp":
		S.Slot = (S.Slot + game.Slots - 1) % game.Slots
		return true
	case "ArrowDown":
		S.Slot = (S.Slot + 1) % game.Slots
		return true
	case "Enter":
		if S.SlotSave {
			if g.SaveSlot(S.Slot) {
				g.Message("COMMANDER SAVED", "good")
				g.Screen = "launch"
			} else {
				g.Message("SAVE FAILED", "bad")
			}
			return true
		}
		if !game.SlotExists(S.Slot) {
			g.Message("SLOT EMPTY", "bad")
			return true
		}
		if !g.LoadSlot(S.Slot) {
			g.Message("LOAD FAILED", "bad")
		} else {
			DropSession()
		}
		return true
	}
	return true
}

// ------------------------------------------------------------------- controls

var (
	capturing   bool
	captureWait float64
	captureHold float64
	captureSlot int // 0 primary key, 1 second key
)

// openControls shows the rebinding table.
func openControls() {
	S.Control = 0
	capturing = false
	game.G.Screen = "controls"
}

// IsCapturing reports whether a key/button is being captured.
func IsCapturing() bool { return capturing }

// CancelCapture abandons the rebinding prompt.
func CancelCapture() { capturing = false }

// PollCapture binds whatever key or pad button is pressed this frame.
func PollCapture(dt float64) {
	if !capturing {
		return
	}
	if captureWait > 0 {
		captureWait -= dt
		return
	}
	kind, code, dir, ok := input.Capture()
	if !ok {
		// a key may already be down when the prompt opens; give the fresh-press
		// path a moment, then accept a held key as well
		if captureHold > 0 {
			captureHold -= dt
			return
		}
		kind, code, ok = input.CaptureDown()
		if !ok {
			return
		}
	}
	acts := input.Actions()
	if S.Control < 0 || S.Control >= len(acts) {
		capturing = false
		return
	}
	var name string
	// only keys have a second slot; a button or an axis replaces its own slot
	if captureSlot == 1 && kind != "key" {
		capturing = false
		game.G.Message("SECOND SLOT TAKES A KEY", "bad")
		return
	}
	switch kind {
	case "key":
		if captureSlot == 1 {
			input.SetSecondKey(acts[S.Control], code)
		} else {
			input.SetKey(acts[S.Control], code)
		}
		name = input.KeyName(code)
	case "axis":
		input.SetAxis(acts[S.Control], code, dir)
		name = input.AxisName(code, dir)
	default:
		input.SetPad(acts[S.Control], code)
		name = input.PadName(code)
	}
	input.Save()
	capturing = false
	if captureSlot == 1 {
		game.G.Message("SECOND KEY: "+name+" FOR "+lang.T(input.Name(acts[S.Control])), "good")
	} else {
		game.G.Message("BOUND: "+name+" FOR "+lang.T(input.Name(acts[S.Control])), "good")
	}
}

// galSystems is a short alias used by the chart markers.
func galSystems(g int) []*galaxy.System { return galaxy.Get(g).Systems }

func controlsScreen(fb *render.FB) {
	acts := input.Actions()
	rows := len(acts) + 2
	const visible = 10
	plate(fb, 4, 16, 312, 184, "CONTROLS")
	fb.Text("ACTION", 12, 28, ColDim, 1)
	fb.TextR("KEY", 234, 28, ColDim, 1)
	fb.TextR("PAD", 306, 28, ColDim, 1)
	off := 0
	if S.Control >= visible {
		off = S.Control - visible + 1
	}
	for i := 0; i < visible && off+i < rows; i++ {
		idx := off + i
		y := 42 + i*11
		sel := S.Control == idx
		if sel {
			fb.FillRect(8, y-2, 308, y+8, render.RGB(0x12384f))
			fb.FillRect(8, y-2, 11, y+8, ColYellow)
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		if idx == len(acts) {
			fb.TextC("RESET DEFAULTS", 160, y, ColGreen, 1)
			continue
		}
		if idx == len(acts)+1 {
			fb.TextC("PAD TEST", 160, y, ColGreen, 1)
			continue
		}
		a := acts[idx]
		fb.Text(input.Name(a), 16, y, col, 1)
		if sel && capturing {
			if seen := input.PeekDown(); seen != "" {
				fb.TextR(seen, 234, y, ColGreen, 1)
			} else if ax, dir, ok := input.AxisPick(); ok {
				fb.TextR(input.AxisName(ax, dir), 234, y, ColGreen, 1)
			} else if captureSlot == 1 {
				fb.TextR("PRESS 2ND", 234, y, ColYellow, 1)
			} else {
				fb.TextR("PRESS KEY", 234, y, ColYellow, 1)
			}
		} else {
			fb.TextR(input.KeyDesc(a), 234, y, ColText, 1)
		}
		fb.TextR(input.PadDesc(a), 306, y, ColGreen, 1)
	}
	if off > 0 {
		fb.Text("^", 296, 36, ColDim, 1)
	}
	if off+visible < rows {
		fb.Text("v", 296, 152, ColDim, 1)
	}
	// legend under the list: three plain lines beat one crammed hint
	lx, ly := 16, 162
	for _, l := range []struct{ key, what string }{
		{"ENTER", "SET KEY"},
		{"SHIFT+ENTER", "SET 2ND KEY"},
		{"BACKSPACE", "RESET ROW"},
	} {
		fb.Text(l.key, lx, ly, ColYellow, 1)
		fb.Text(l.what, lx+80, ly, ColText, 1)
		ly += 11
	}
	switch {
	case capturing && captureSlot == 1:
		Hint = "PRESS THE SECOND KEY   ESC CANCEL"
	case capturing:
		Hint = "PRESS KEY, BUTTON OR STICK   ESC CANCEL"
	case input.PeekDown() != "":
		// live read-out: makes it obvious whether the platform sends the key
		Hint = "KEY SEEN: " + input.PeekDown() + "   ESC BACK"
	default:
		Hint = "UP/DOWN SELECT   ESC BACK"
	}
}

func controlsKey(k string) bool {
	acts := input.Actions()
	rows := len(acts) + 2 // actions, RESET DEFAULTS, PAD TEST
	switch k {
	case "Escape":
		capturing = false
		game.G.Screen = "settings"
		S.Settings = 3
		return true
	case "ArrowUp":
		capturing = false
		S.Control = (S.Control + rows - 1) % rows
		return true
	case "ArrowDown":
		capturing = false
		S.Control = (S.Control + 1) % rows
		return true
	case "Enter":
		if S.Control == rows-2 {
			input.Reset()
			input.Save()
			game.G.Message("CONTROLS RESET", "good")
			return true
		}
		if S.Control == rows-1 {
			openPadTest("controls")
			return true
		}
		capturing = true
		captureSlot = 0
		if shiftHeld() {
			captureSlot = 1
		}
		captureWait = 0.05
		captureHold = 0.4
		input.BeginCapture()
		return true
	case "Backspace":
		if S.Control < len(acts) {
			input.ResetAction(acts[S.Control])
			input.Save()
			game.G.Message("RESTORED: "+input.Name(acts[S.Control]), "good")
		}
		return true
	}
	return true
}

// shiftHeld reports whether a shift modifier is down; the controls screen uses
// it to tell "set the primary key" from "set the second key".
func shiftHeld() bool {
	return ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
}

// ---------------------------------------------------------------- pad test

// openPadTest shows, live, what the gamepad actually reports: raw SDL axes,
// the standard stick axes and every button. It exists so a mapping problem can
// be diagnosed by wiggling a stick instead of guessing.
func openPadTest(from string) {
	S.PadTestFrom = from
	game.G.Screen = "padtest"
}

// padState is what the pad test screen draws; kept apart from Ebiten so the
// screen can be rendered with sample values in the reference shots.
type padState struct {
	ok     bool
	name   string
	count  int
	raw    []float64
	std    [4]float64
	btn    [17]bool
	hasStd bool
	active []string
}

func readPadState() padState {
	var p padState
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) == 0 {
		return p
	}
	id := ids[0]
	p.ok = true
	p.name = ebiten.GamepadName(id)
	p.count = len(ids)
	p.hasStd = ebiten.IsStandardGamepadLayoutAvailable(id)
	n := ebiten.GamepadAxisCount(id)
	if n > 6 {
		n = 6
	}
	for ax := 0; ax < n; ax++ {
		p.raw = append(p.raw, ebiten.GamepadAxisValue(id, ebiten.GamepadAxisType(ax)))
	}
	for ax := 0; ax < 4; ax++ {
		p.std[ax] = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxis(ax))
	}
	for b := 0; b < 17; b++ {
		p.btn[b] = ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButton(b))
	}
	for _, a := range input.Actions() {
		if input.Held(a) {
			p.active = append(p.active, lang.T(input.Name(a)))
		}
	}
	return p
}

// samplePad draws a plausible pad so the screen can be screenshotted without
// hardware attached.
func samplePad() padState {
	p := padState{ok: true, name: "Sample Pad", count: 1, hasStd: true}
	p.raw = []float64{0.62, -0.18, 0.00, -0.45, 0.12, 0.00}
	p.std = [4]float64{0.62, -0.18, -0.45, 0.12}
	p.btn[input.PadA] = true
	p.btn[input.PadUp] = true
	p.btn[input.PadRT] = true
	p.active = []string{lang.T("ACCELERATE"), lang.T("FIRE")}
	return p
}

func padTestScreen(fb *render.FB) { drawPadTest(fb, readPadState()) }

// DrawPadTestSample renders the pad test screen with sample values.
func DrawPadTestSample(fb *render.FB) { drawPadTest(fb, samplePad()) }

func drawPadTest(fb *render.FB, p padState) {
	plate(fb, 4, 14, 312, 186, "PAD TEST")
	if !p.ok {
		fb.TextC("NO GAMEPAD CONNECTED", 160, 100, ColRed, 1)
		fb.TextC("CONNECT ONE AND PRESS A BUTTON", 160, 116, ColDim, 1)
		Hint = "ESC BACK"
		return
	}
	fb.Text(strconv.Itoa(p.count)+" "+lang.T("PAD")+"  "+p.name, 12, 24, ColDim, 1)

	fb.Text("RAW AXES", 12, 40, ColYellow, 1)
	for ax, v := range p.raw {
		y := 50 + ax*9
		fb.Text(strconv.Itoa(ax)+" "+strconv.FormatFloat(v, 'f', 2, 64), 12, y, ColText, 1)
		axisBar(fb, 74, y, v)
	}

	fb.Text("STD AXES", 186, 40, ColYellow, 1)
	names := []string{"LX", "LY", "RX", "RY"}
	for ax := 0; ax < 4; ax++ {
		y := 50 + ax*9
		fb.Text(names[ax]+" "+strconv.FormatFloat(p.std[ax], 'f', 2, 64), 186, y, ColText, 1)
	}
	if !p.hasStd {
		fb.Text("NO STANDARD MAPPING", 186, 88, ColRed, 1)
	}

	fb.Text("BUTTONS 0-15", 12, 106, ColYellow, 1)
	for b := 0; b < 16; b++ {
		x := 12 + (b%8)*37
		y := 116 + (b/8)*11
		col := render.RGB(0x24425c)
		if p.btn[b] {
			col = ColGreen
		}
		fb.Text(strconv.Itoa(b), x, y, col, 1)
		if b < 5 || b == 9 {
			fb.Text(input.PadName(int32(b)), x+13, y, col, 1)
		}
	}

	on := func(b int32) string {
		if p.btn[b] {
			return "DOWN"
		}
		return "-"
	}
	fb.Text("D-UP "+on(input.PadUp)+"  DOWN "+on(input.PadDown)+"  LEFT "+on(input.PadLeft)+"  RIGHT "+on(input.PadRight), 12, 142, ColText, 1)
	fb.Text("LT "+on(input.PadLT)+"  RT "+on(input.PadRT)+"  LS "+on(input.PadLS)+"  RS "+on(input.PadRS), 12, 153, ColText, 1)

	// what the game itself thinks is pressed — the end-to-end check that a
	// binding actually responds
	fb.Text("ACTIVE: "+strings.Join(p.active, " "), 12, 168, ColGreen, 1)
	Hint = "WIGGLE STICKS AND PRESS BUTTONS   ESC BACK"
}

// axisBar draws a centred bar so a stick position is readable at a glance.
func axisBar(fb *render.FB, x, y int, v float64) {
	const w = 84
	fb.FillRect(x, y+2, x+w, y+4, render.RGB(0x24425c))
	c := int((v + 1) / 2 * float64(w-2))
	if c < 0 {
		c = 0
	}
	if c > w-2 {
		c = w - 2
	}
	fb.FillRect(x+c, y, x+c+2, y+7, ColGreen)
}

// QuitToTitle abandons the current flight and returns to the title screen.
// Leaving the application is a separate choice, made from the main menu.
func QuitToTitle() {
	g := game.G
	// remember where the pilot was, so CONTINUE FLIGHT can pick it up
	S.HasSession = g.State != "title"
	S.SessionState = g.State
	S.SessionScreen = g.Screen
	g.Overlay = ""
	g.Screen = "launch"
	g.In = game.InputState{}
	g.State = "title"
	S.Title = 0
}

// ContinueFlight returns to the flight that was left from the pause menu.
func ContinueFlight() {
	if !S.HasSession {
		return
	}
	g := game.G
	g.State = S.SessionState
	g.Screen = S.SessionScreen
	g.Overlay = ""
	S.Title = 0
}

// DropSession forgets the interrupted flight: a new commander or a loaded save
// replaces it.
func DropSession() {
	S.HasSession = false
	S.SessionState = ""
	S.SessionScreen = ""
}

// ----------------------------------------------------------------- game mode

// openMode asks for the rule set before a commander is created.
func openMode() {
	S.Mode = 0
	game.G.Screen = "mode"
}

func modeScreen(fb *render.FB) {
	rows := []struct{ name, desc string }{
		{"CLASSIC", "1984 RULES: TWO MISSIONS, NO CONTRACTS"},
		{"ADVANCED", "EVERYTHING: CONTRACTS, SHIPYARD, EXTRAS"},
	}
	plate(fb, 14, 34, 292, 120, "")
	for i, r := range rows {
		y := 52 + i*34
		sel := S.Mode == i
		if sel {
			fb.FillRect(22, y-8, 298, y+22, render.RGB(0x12384f))
			fb.FillRect(22, y-8, 25, y+22, ColYellow)
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.Text(r.name, 34, y, col, 1)
		fb.Text(r.desc, 34, y+11, ColDim, 1)
	}
	fb.TextC("THE MODE IS FIXED FOR THIS COMMANDER", 160, 138, ColRed, 1)
	Hint = "UP/DOWN   ENTER START   ESC BACK"
}

func modeKey(k string) bool {
	g := game.G
	switch k {
	case "Escape":
		g.Screen = "launch"
		S.Title = 0
		return true
	case "ArrowUp", "ArrowDown":
		S.Mode = 1 - S.Mode
		return true
	case "Enter":
		mode := game.ModeAdvanced
		if S.Mode == 0 {
			mode = game.ModeClassic
		}
		g.NewGameMode(mode, "")
		DropSession()
		return true
	}
	return true
}

// ------------------------------------------------------------------- journal

// openJournal shows the story contract and every accepted job in one place.
func openJournal() {
	g := game.G
	if g.Cmd.Mission.Kind == "" && len(g.Cmd.Contracts) == 0 {
		g.Message("NOTHING TO REPORT", "bad")
		return
	}
	g.Screen = "journal"
}

func journalScreen(fb *render.FB) {
	g := game.G
	plate(fb, 8, 16, 304, 182, "JOURNAL")
	y := 32
	if !g.Classic() {
		c := &g.Cmd.Campaign
		fb.Text("CAMPAIGN", 16, y, ColDim, 1)
		y += 11
		if c.Stage >= game.CampaignTotal() {
			fb.Text("COMPLETE", 16, y, ColGreen, 1)
		} else {
			q := g.CampaignQuest()
			fb.Text(clip(q.Title, 46), 16, y, ColYellow, 1)
			y += 10
			if st := g.CampaignStatus(); st != "" {
				fb.Text(clip(st, 46), 16, y, ColGreen, 1)
				y += 10
			}
			fb.Text("STAGE "+itoa(c.Stage+1)+"/"+itoa(game.CampaignTotal()), 16, y, ColDim, 1)
		}
		y += 14
		fb.HLine(16, 300, y, ColFrame)
		y += 8
	} else if ms := g.MissionStatus(); ms != "" {
		fb.Text("STORY CONTRACT", 16, y, ColDim, 1)
		y += 11
		fb.Text(clip(ms, 46), 16, y, ColGreen, 1)
		y += 10
		fb.Text("PAYS "+itoa(g.Cmd.Mission.Pay)+" CR", 16, y, ColDim, 1)
		y += 12
	} else {
		fb.Text("STORY CONTRACT", 16, y, ColDim, 1)
		y += 11
		fb.Text("NONE ACTIVE", 16, y, render.RGB(0x2d4a5c), 1)
		y += 12
	}
	fb.HLine(16, 300, y, ColFrame)
	y += 8
	fb.Text("ACCEPTED JOBS", 16, y, ColDim, 1)
	y += 12
	if len(g.Cmd.Contracts) == 0 {
		fb.Text("NONE", 16, y, render.RGB(0x2d4a5c), 1)
		y += 11
	} else {
		for _, c := range g.Cmd.Contracts {
			fb.Text(clip(g.ContractDetail(c), 44), 16, y, ColText, 1)
			y += 11
		}
	}
	y += 6
	if g.Cmd.Pax > 0 || g.Cmd.Mail > 0 {
		fb.Text("PAX "+itoa(g.Cmd.Pax)+"   MAIL "+itoa(g.Cmd.Mail)+" T", 16, y, ColDim, 1)
	}
	Hint = "ESC BACK"
}
