package ui

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

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
	// composed strings bypass the dictionary, so translate the word itself
	fb.TextR(lang.T("GALAXY")+" "+itoa(g.Cmd.Galaxy+1)+"   "+money(g.Cmd.Cash)+" CR", 314, 4, ColText, 1)
}

func footer(fb *render.FB) {
	fb.FillRect(0, 211, 319, 223, render.RGB(0x08131f))
	fb.HLine(0, 319, 211, ColFrame)
	fb.TextC(Hint, 160, 215, ColDim, 1)
	// the build number sits in the free corner: the hint is centred, and the
	// right side of the footer collided with it
	fb.Text(game.Version, 3, 215, render.RGB(0x24414f), 1)
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
	case "log":
		logScreen(fb)
	case "closing":
		closingScreen(fb)
	case "ending":
		endingScreen(fb)
	case "talk":
		talkScreen(fb)
	case "shipyard":
		shipyardScreen(fb)
	case "slots":
		slotsScreen(fb)
	default:
		stationMenu(fb)
	}
	footer(fb)
	// the answer to "why did nothing happen" belongs on every screen, not just
	// in flight - except the contract board, where a running list of toasts
	// under the description reads like part of the contract
	// the market and the board say nothing when a thing works: the log records
	// it. They still say why something did not happen.
	switch g.Screen {
	case "contracts", "market":
		drawMessagesWhere(fb, 192, 2, true)
	default:
		drawMessages(fb, 192, 2)
	}
}

func stationMenu(fb *render.FB) {
	g := game.G
	cmd := g.Cmd
	Hint = "UP/DOWN   ENTER CONFIRM   F REFUEL   Q MAIN MENU"

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
	fb.Text("TARGET", 16, 146, ColDim, 1)
	drawDataR(fb, gal.Systems[((cmd.TargetSystem%256)+256)%256].Name, 150, 146, ColText)

	// system block at the bottom of the same plate
	fb.HLine(16, 150, 160, render.RGB(0x1b3348))
	if g.World != nil {
		sys := g.World.Sys
		drawData(fb, sys.Name, 16, 168, ColYellow)
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
	if S.Menu >= 0 && S.Menu < len(items) && items[S.Menu] == "DIPLOMAT" {
		Hint = lang.T("CLEAR RECORD FOR") + " " + itoa(game.G.PardonPrice()) + " CR"
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
	// The sheet first, because that is what a pilot checks on arrival; trade,
	// then work, then the ship. The chart sits second from the bottom - it is a
	// map, not a service - and the way out is always the last row.
	items := []string{"COMMANDER STATUS", "BUY / SELL CARGO"}
	if !game.G.Classic() && (len(game.G.Offers) > 0 || len(game.G.Cmd.Contracts) > 0 ||
		game.G.CampaignOffer() || game.G.Cmd.Campaign.Active) {
		items = append(items, "CONTRACTS")
	}
	items = append(items, "EQUIP SHIP")
	if game.G.Classic() && (game.G.MissionOffer() != "" || game.G.Cmd.Mission.Kind != "") {
		// the 1984 build has no contract board: its two story jobs get a row
		items = append(items, "MISSION")
	}
	if !game.G.Classic() {
		items = append(items, "SHIPYARD", "LOG")
	}
	if !game.G.Classic() {
		// always, not only when there is work: the row is how a pilot asks
		// "where am I in all this", and a missing row reads as a missing feature
		items = append(items, "JOURNAL")
	}
	if game.G.PardonAvailable() {
		items = append(items, "DIPLOMAT")
	}
	// the way out is always the last row
	items = append(items, "SAVE COMMANDER", "GALACTIC CHART", "LEAVE STATION")
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
	fb.TextR(lang.T("SELL"), 150, 30, ColDim, 1)
	fb.TextR(lang.T("BUY"), 196, 30, ColDim, 1)
	fb.TextR(lang.T("QTY"), 246, 30, ColDim, 1)
	fb.TextR(lang.T("HOLD"), 302, 30, ColDim, 1)
	fb.HLine(10, 310, 37, ColFrame)
	for i, m := range g.Market {
		y := 42 + i*8
		sel := S.Market == i
		if sel {
			fb.FillRect(8, y-1, 311, y+7, render.RGB(0x12384f))
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		drawIcon(fb, 10, y, i, sel)
		fb.Text(m.Name, 20, y, col, 1)
		if m.Unit == "KG" {
			fb.Text("KG", 104, y, ColDim, 1)
		}
		sc := ColDim
		if m.Sell > m.Price {
			sc = ColGreen
		}
		fb.TextR(itoa(m.Sell), 150, y, sc, 1)
		pc := ColGreen
		if sel {
			pc = ColWhite
		}
		fb.TextR(itoa(m.Price), 196, y, pc, 1)
		fb.TextR(itoa(m.Qty), 246, y, ColDim, 1)
		have := cmd.Cargo[i]
		hc := ColDim
		if have > 0 {
			hc = ColYellow
		}
		fb.TextR(itoa(have), 302, y, hc, 1)
	}
	// amount prompt: how many tonnes, before anything is bought or sold
	if S.MarketMode != "" && S.Market < len(g.Market) {
		m := g.Market[S.Market]
		selling := S.MarketMode == "sell"
		maxQty := marketMax(S.Market, selling)
		if maxQty < 1 {
			// nothing to sell, nothing to buy: do not open the prompt at all
			S.MarketMode = ""
			if selling {
				g.Message("NOTHING TO SELL", "bad")
			} else {
				g.Message("CANNOT BUY", "bad")
			}
		}
		qty := S.MarketQty
		if qty > maxQty {
			qty = maxQty
			S.MarketQty = qty
		}
		if qty < 1 {
			qty = 1
		}
		unit, total := m.Price, qty*m.Price
		label := lang.T("BUY")
		if selling {
			unit, total = m.Sell, qty*m.Sell
			label = lang.T("SELL")
		}
		fb.FillRect(52, 128, 268, 180, render.RGB(0x08131f))
		fb.Rect(52, 128, 268, 180, ColYellow)
		fb.TextC(label+" "+itoa(qty)+" X "+m.Name, 160, 134, ColYellow, 1)
		fb.TextC(itoa(unit)+" CR X "+itoa(qty)+" = "+itoa(total)+" CR", 160, 146, ColGreen, 1)
		fb.TextC(lang.T("MAX")+" "+itoa(maxQty), 160, 158, ColDim, 1)
		fb.TextC(lang.T("UP/DOWN  ENTER  ESC"), 160, 168, ColDim, 1)
		Hint = "UP/DOWN  AMOUNT   ENTER  ESC"
		return
	}

	// the crate rides in the hold and the market will not touch it: show it
	// where a player looks for their cargo, and let him try
	if cmd.Boxes > 0 {
		y := 42 + len(g.Market)*8
		sel := S.Market == len(g.Market)
		if sel {
			fb.FillRect(8, y-1, 311, y+7, render.RGB(0x12384f))
		}
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.Text("CRATE x"+itoa(cmd.Boxes), 20, y, col, 1)
		fb.TextR(itoa(40000*cmd.Boxes)+" NOM", 302, y, ColDim, 1)
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
	items := galaxy.AvailableEquipment(g.World.Sys, g.Classic())

	// two columns: the catalogue outgrew a single list
	const rows = 9
	for i, it := range items {
		col := i / rows
		row := i % rows
		x := 12 + col*150
		y := 32 + row*17
		sel := S.Equip == i
		if sel {
			fb.FillRect(x-4, y-3, x+144, y+8, render.RGB(0x12384f))
		}
		// the shop says why before it charges: a rig with no laser, a laser
		// the hull cannot carry, something already bolted on
		blocked := g.EquipBlocked(it.ID)
		owned := blocked == "ALREADY FITTED" && g.EquipFitted(it.ID)
		c := ColText
		switch {
		case owned:
			c = ColDim
		case blocked != "":
			c = render.RGB(0x6a7f8c) // available, but not to you yet
		case S.Equip == i:
			c = ColYellow
		}
		name := it.Name
		if len(name) > 18 {
			name = name[:18]
		}
		fb.Text(name, x, y, c, 1)
		switch {
		case owned:
			fb.TextR(lang.T("FITTED"), x+144, y, ColGreen, 1)
		case it.ID == "missile":
			fb.TextR(itoa(it.Price)+" x"+itoa(cmd.Missiles), x+144, y, ColGreen, 1)
		default:
			pc := ColRed
			if cmd.Cash >= float64(it.Price) {
				pc = ColGreen
			}
			fb.TextR(itoa(it.Price), x+144, y, pc, 1)
		}
	}
	// the reason a row cannot be bought goes in the footer, where the player
	// is already looking for the key hints
	if S.Equip >= 0 && S.Equip < len(items) {
		if why := g.EquipBlocked(items[S.Equip].ID); why != "" && why != "ALREADY FITTED" {
			Hint = lang.T(why)
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

	plate(fb, 4, 18, 196, 190, lang.T("GALACTIC CHART")+" "+itoa(cmd.Galaxy+1))
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
	if cmd.Boxes > 0 {
		any = true
		fb.Text("CRATE x"+itoa(cmd.Boxes), 168, yy, ColYellow, 1)
		fb.TextR(itoa(40000*cmd.Boxes)+" NOM", 308, yy, ColDim, 1)
		yy += 9
	}
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
	plate(fb, 8, 18, 312, 190, "FLIGHT CONTROLS")
	steer := "PITCH AND ROLL"
	if !game.G.Classic() {
		steer = "PITCH / ROLL / YAW"
	}
	rows := [][2]string{
		{"ARROWS / MOUSE", steer},
		{"Q / R", "YAW LEFT / RIGHT"},
		{"X / Z", "ACCELERATE / BRAKE"},
		{"SPACE / LMB", "FIRE LASERS"},
		{"M", "LAUNCH MISSILE"},
		{"E", "E.C.M."},
		{"TAB", "ENERGY BOMB"},
		{"L", "DOCKING COMPUTER"},
		{"J", "IN-SYSTEM JUMP"},
		{"C", "GALACTIC CHART"},
		{"V", "COMMANDER STATUS"},
		{"ESC", "RESUME"},
		{"H", "HANDBOOK"},
		{"Q", "MAIN MENU"},
	}
	// one column, tight rows: fourteen lines have to fit the plate
	for i, r := range rows {
		y := 30 + i*12
		fb.Text(r[0], 16, y, ColYellow, 1)
		fb.TextR(r[1], 300, y, ColText, 1)
	}
}

// Key handles a key press; returns true when consumed.
func Key(k string) bool {
	g := game.G
	// a conversation owns every key, in flight as much as at a counter: it used
	// to be checked after the flight branch, which ate the digits
	if g.Screen == "talk" {
		return talkKey(k)
	}
	// the settings screen sits on top of whatever state opened it, so it gets
	// first refusal on every key
	if g.Screen == "settings" {
		return settingsKey(k)
	}
	if g.Screen == "controls" {
		return controlsKey(k)
	}
	if g.Screen == "help" {
		return helpKey(k)
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
			case "HOW TO PLAY":
				S.Help = 0
				game.G.Screen = "help"
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
	// The station has no pause, so Q is the way out to the title. It has to be
	// checked before the screens, because each of them answers every key it is
	// given and would swallow the Q.
	if k == "q" || k == "Q" {
		QuitToTitle()
		return true
	}
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
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			// a hotkey means "this row, confirmed": the digit on screen is the
			// digit that works, whatever order the services are in
			n := 10
			if k != "0" {
				n = int(k[0] - '0')
			}
			if items := menuItems(); n <= len(items) {
				S.Menu = n - 1
				return dockedKey("Enter")
			}
			return true
		case "Enter":
			it := menuItems()
			if S.Menu < 0 || S.Menu >= len(it) {
				return true
			}
			switch it[S.Menu] {
			case "BUY / SELL CARGO":
				g.Screen = "market"
				S.MarketMode = ""
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
			case "LEAVE STATION":
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
			case "LOG":
				g.Screen = "log"
				S.Log = 0
			case "DIPLOMAT":
				g.BuyPardon()
			}
			return true
		}
		return true
	}
	if g.Screen == "journal" {
		journalKey(k)
		return true
	}
	if g.Screen == "ending" {
		return endingKey(k)
	}
	if g.Screen == "closing" {
		return closingKey(k)
	}
	if g.Screen == "log" {
		return logKey(k)
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
		// the amount prompt is a layer of its own: Escape closes it and stays
		// in the market, so a purchase can be called off
		if g.Screen == "market" && S.MarketMode != "" {
			S.MarketMode = ""
			return true
		}
		g.Screen = "launch"
		return true
	}
	switch g.Screen {
	case "market":
		if S.MarketMode != "" {
			// the amount prompt owns the keys while it is up
			switch k {
			case "Escape":
				S.MarketMode = ""
			case "ArrowUp":
				S.MarketQty++
			case "ArrowDown":
				S.MarketQty--
			case "Enter":
				qty := S.MarketQty
				if qty > 0 && S.Market < len(g.Market) {
					if S.MarketMode == "sell" {
						g.Sell(S.Market, qty)
					} else {
						g.Buy(S.Market, qty)
					}
				}
				S.MarketMode = ""
			}
			if max := marketMax(S.Market, S.MarketMode == "sell"); max > 0 {
				if S.MarketQty > max {
					S.MarketQty = max
				}
			}
			if S.MarketQty < 1 {
				S.MarketQty = 1
			}
			return true
		}
		switch k {
		case "ArrowUp":
			S.Market = (S.Market + marketRows() - 1) % marketRows()
		case "ArrowDown":
			S.Market = (S.Market + 1) % marketRows()
		case "b", "B", "s", "S":
			if S.Market >= len(g.Market) {
				// the crate is owned by whoever holds it
				g.Message("OWNER IS THE HOLDER", "bad")
				return true
			}
			selling := !(k == "b" || k == "B")
			if max := marketMax(S.Market, selling); max < 1 {
				if selling {
					g.Message("NOTHING TO SELL", "bad")
				} else {
					g.Message("CANNOT BUY", "bad")
				}
				return true
			}
			if selling {
				S.MarketMode, S.MarketQty = "sell", 1
			} else {
				S.MarketMode, S.MarketQty = "buy", 1
			}
		}
		return true
	case "equip":
		items := galaxy.AvailableEquipment(g.World.Sys, g.Classic())
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
	if k == "h" && g.Overlay == "pause" {
		// the handbook from inside a flight: leave the overlay open behind it
		S.PrevScreen = g.Screen
		g.Screen = "help"
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

// settingsRows is built each time: the tester switch only exists once the
// phrase has been typed, and rows are addressed by name so adding one cannot
// silently shift the key handling.
func settingsRows() []string {
	rows := []string{"LANGUAGE", "SOUND", "VOLUME", "COCKPIT"}
	if lang.Cfg.TesterUnlocked {
		rows = append(rows, "TESTER MODE")
	}
	return append(rows, "CONTROLS", "BACK")
}

func settingsScreen(fb *render.FB) {
	plate(fb, 40, 18, 280, 186, "")
	rows := settingsRows()
	if S.Settings >= len(rows) {
		S.Settings = len(rows) - 1
	}
	for i, name := range rows {
		y := 34 + i*20
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
		switch name {
		case "LANGUAGE":
			val := "ENGLISH"
			if lang.IsRU() {
				val = "RUSSIAN"
			}
			fb.TextR(val, 264, y, ColGreen, 1)
		case "SOUND":
			val := "ON"
			if lang.Cfg.Muted {
				val = "OFF"
			}
			fb.TextR(val, 264, y, ColGreen, 1)
		case "VOLUME":
			pips := int(lang.Cfg.Volume*5 + 0.5)
			for k := 0; k < 5; k++ {
				c := render.RGB(0x25405a)
				if k < pips {
					c = ColGreen
				}
				fb.FillRect(224+k*9, y-1, 230+k*9, y+6, c)
			}
		case "COCKPIT":
			val := "ON"
			if lang.Cfg.NoCockpit {
				val = "OFF"
			}
			fb.TextR(val, 264, y, ColGreen, 1)
		case "TESTER MODE":
			val := "OFF"
			if lang.Cfg.TesterMode {
				val = "ON"
			}
			fb.TextR(val, 264, y, ColGreen, 1)
		case "CONTROLS":
			fb.TextR(">", 264, y, ColGreen, 1)
		}
	}
	fb.TextC("SHIFT+M - SOUND ON/OFF", 160, 176, ColDim, 1)
	Hint = "UP/DOWN  LEFT/RIGHT  ENTER  ESC"
}

func settingsKey(k string) bool {
	rows := settingsRows()
	n := len(rows)
	switch k {
	case "ArrowUp":
		S.Settings = (S.Settings + n - 1) % n
	case "ArrowDown":
		S.Settings = (S.Settings + 1) % n
	case "Escape":
		closeSettings()
	case "Enter", "ArrowLeft", "ArrowRight":
		if S.Settings >= n {
			return true
		}
		switch rows[S.Settings] {
		case "LANGUAGE":
			if lang.IsRU() {
				lang.Set("en")
			} else {
				lang.Set("ru")
			}
			lang.Save()
		case "SOUND":
			lang.Cfg.Muted = !lang.Cfg.Muted
			lang.Save()
			state := "ON"
			if lang.Cfg.Muted {
				state = "OFF"
			}
			game.G.Message(lang.T("SOUND")+" "+lang.T(state), "info")
		case "VOLUME":
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
		case "COCKPIT":
			lang.Cfg.NoCockpit = !lang.Cfg.NoCockpit
			lang.Save()
			state := "ON"
			if lang.Cfg.NoCockpit {
				state = "OFF"
			}
			game.G.Message(lang.T("COCKPIT")+" "+lang.T(state), "info")
		case "TESTER MODE":
			lang.Cfg.TesterMode = !lang.Cfg.TesterMode
			lang.Save()
			state := "OFF"
			if lang.Cfg.TesterMode {
				state = "ON"
			}
			game.G.Message(lang.T("TESTER MODE")+" "+lang.T(state), "info")
		case "CONTROLS":
			openControls()
		case "BACK":
			closeSettings()
		}
	default:
		// the phrase that reveals the tester switch: typed, never listed
		if len(k) == 1 {
			S.Typed = (S.Typed + k)
			if len(S.Typed) > 8 {
				S.Typed = S.Typed[len(S.Typed)-8:]
			}
			if !lang.Cfg.TesterUnlocked && strings.HasSuffix(strings.ToUpper(S.Typed), "IDDQD") {
				lang.Cfg.TesterUnlocked = true
				lang.Save()
				S.Typed = ""
				game.G.Message("TESTER MODE UNLOCKED", "good")
			}
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

// boardRow is one line of the contract board: an offer taken or not, the
// story job when there is one, or a contract still running from another
// station.
type boardRow struct {
	kind string // offer | story | active
	idx  int
}

func sameContract(a, b game.Contract) bool {
	return a.Kind == b.Kind && a.Good == b.Good && a.System == b.System &&
		a.Galaxy == b.Galaxy && a.Pay == b.Pay && a.Qty == b.Qty
}

func boardRows() []boardRow {
	g := game.G
	var rows []boardRow
	// the story leads: it is the one job that is not a number
	if g.CampaignOffer() || g.Cmd.Campaign.Active {
		rows = append(rows, boardRow{"story", 0})
	}
	for i := range g.Offers {
		rows = append(rows, boardRow{"offer", i})
	}
	for i, c := range g.Cmd.Contracts {
		dup := false
		for _, o := range g.Offers {
			if o.Taken && sameContract(o, c) {
				dup = true
			}
		}
		if !dup {
			rows = append(rows, boardRow{"active", i})
		}
	}
	return rows
}

// contractWhat is the description template: what the job actually asks for.
func contractWhat(c game.Contract) []string {
	sys := galaxy.Get(c.Galaxy).Systems[c.System].Name
	good := galaxy.Commodities[c.Good].Name
	days := itoa(maxInt(0, c.Due-game.G.Cmd.Day))
	switch c.Kind {
	case game.ContractCourier:
		return []string{
			lang.T("TAKE THIS CARGO TO ANOTHER STATION."),
			lang.T("CARGO") + ": " + lang.T(good) + " X" + itoa(c.Qty) + " T",
			lang.T("DESTINATION") + ": " + sys,
			lang.T("THE CARGO IS ALREADY ABOARD."),
			lang.T("DAYS LEFT") + ": " + days,
		}
	case game.ContractSmuggle:
		return []string{
			lang.T("MOVE THIS CARGO WITHOUT BEING SEARCHED."),
			lang.T("CARGO") + ": " + lang.T(good) + " X" + itoa(c.Qty) + " T",
			lang.T("DESTINATION") + ": " + sys,
			lang.T("POLICE MAY STOP YOU ON THE WAY."),
			lang.T("DAYS LEFT") + ": " + days,
		}
	case game.ContractPassenger:
		return []string{
			lang.T("TAKE PASSENGERS TO ANOTHER STATION."),
			lang.T("PASSENGERS") + ": " + itoa(c.Qty),
			lang.T("DESTINATION") + ": " + sys,
			lang.T("THEY TAKE HOLD SPACE AND THEY COMPLAIN."),
			lang.T("DAYS LEFT") + ": " + days,
		}
	case game.ContractMail:
		return []string{
			lang.T("CARRY SEALED POST TO ANOTHER STATION."),
			lang.T("MAIL") + ": " + itoa(c.Qty) + " T",
			lang.T("DESTINATION") + ": " + sys,
			lang.T("DAYS LEFT") + ": " + days,
		}
	}
	return []string{
		lang.T("CLEAR HOSTILES IN ONE SYSTEM."),
		lang.T("TARGET") + ": " + itoa(c.Qty) + " " + lang.T("SHIPS"),
		lang.T("SYSTEM") + ": " + sys,
		lang.T("PROGRESS") + ": " + itoa(c.Done) + "/" + itoa(c.Qty),
		lang.T("DAYS LEFT") + ": " + days,
	}
}

// storyWhat is the story job in the same shape, so one board reads one way.
func storyWhat() []string {
	g := game.G
	q := g.CampaignQuest()
	para := ""
	for _, ln := range q.Brief {
		if ln != "" {
			para += lang.T(ln) + " "
		}
	}
	if g.Cmd.Campaign.Active {
		out := []string{lang.T("PRIVATE CONTRACT.")}
		out = append(out, wrapText(para, 48)...)
		out = append(out, g.CampaignStatus())
		return out
	}
	out := []string{lang.T("PRIVATE CONTRACT. THE COMPANY PAYS.")}
	out = append(out, wrapText(para, 48)...)
	if q.Pay > 0 {
		out = append(out, lang.T("PAYS")+" "+itoa(q.Pay)+" CR")
	}
	return out
}

func contractsScreen(fb *render.FB) {
	g := game.G
	plate(fb, 6, 14, 310, 172, "")
	rows := boardRows()
	if len(rows) == 0 {
		fb.Text(lang.T("NO CONTRACTS ON THE BOARD"), 14, 30, ColDim, 1)
		Hint = "ESC BACK"
		return
	}
	if S.Contract >= len(rows) {
		S.Contract = len(rows) - 1
	}
	// the list never reorders itself: an accepted offer stays where it was and
	// only gains a marker
	y := 26
	for i, r := range rows {
		sel := S.Contract == i
		if sel {
			fb.FillRect(10, y-2, 306, y+8, render.RGB(0x12384f))
		}
		mark := " "
		col := ColText
		var label string
		switch r.kind {
		case "offer":
			c := g.Offers[r.idx]
			label = g.ContractLabel(c)
			if c.Taken {
				mark = "*"
				col = ColGreen
			}
		case "story":
			label = lang.T("PRIVATE") + ": " + lang.T(g.CampaignQuest().Title)
			if g.Cmd.Campaign.Active {
				mark = "*"
				col = ColGreen
			} else {
				// an open private order: the only row on the board that asks to
				// be taken
				mark = "!"
				col = ColYellow
			}
		case "active":
			c := g.Cmd.Contracts[r.idx]
			label = g.ContractLabel(c)
			mark = "*"
			col = ColGreen
		}
		if sel {
			col = ColYellow
		}
		fb.Text(mark, 12, y, ColGreen, 1)
		fb.TextRaw(clip(label, 45), 20, y, col, 1)
		y += 10
	}
	// description of whatever is under the cursor: it takes the room the list
	// leaves, so a short board shows a long briefing
	descTop := y + 8
	if descTop < 70 {
		descTop = 70
	}
	fb.HLine(10, 306, descTop-6, render.RGB(0x1b3348))
	var lines []string
	switch rows[S.Contract].kind {
	case "offer":
		lines = contractWhat(g.Offers[rows[S.Contract].idx])
	case "active":
		lines = contractWhat(g.Cmd.Contracts[rows[S.Contract].idx])
	case "story":
		lines = storyWhat()
	}
	// the last row belongs to the status line
	const statusY = 180
	const lineH = 9
	maxLines := (statusY - 4 - descTop) / lineH
	if maxLines < 1 {
		maxLines = 1
	}
	if S.ContractScroll > len(lines)-maxLines {
		S.ContractScroll = len(lines) - maxLines
	}
	if S.ContractScroll < 0 {
		S.ContractScroll = 0
	}
	dy := descTop
	for i := S.ContractScroll; i < len(lines) && i < S.ContractScroll+maxLines; i++ {
		fb.TextRaw(lines[i], 12, dy, ColText, 1)
		dy += lineH
	}
	if len(lines) > maxLines {
		// say there is more, and which part of it you are looking at
		fb.TextR(itoa(S.ContractScroll+1)+"/"+itoa(len(lines)-maxLines+1), 300, descTop, ColDim, 1)
	}
	// and its status, right under the description
	var status string
	switch rows[S.Contract].kind {
	case "offer":
		if g.Offers[rows[S.Contract].idx].Taken {
			status = lang.T("CONTRACT ACCEPTED")
		} else {
			status = lang.T("ENTER ACCEPT")
		}
	case "story":
		if g.Cmd.Campaign.Active {
			status = lang.T("CONTRACT ACCEPTED")
		} else {
			status = lang.T("ENTER ACCEPT")
		}
	case "active":
		status = lang.T("CONTRACT ACCEPTED")
	}
	fb.TextRaw(status, 12, statusY, ColYellow, 1)
	Hint = "UP/DOWN SELECT   ENTER ACCEPT   ESC BACK"
}

func contractsKey(k string) bool {
	g := game.G
	rows := boardRows()
	n := len(rows)
	switch k {
	case "Escape":
		g.Screen = "launch"
		return true
	case "ArrowDown":
		if n > 0 {
			S.Contract = (S.Contract + 1) % n
			S.ContractScroll = 0
		}
		return true
	case "ArrowUp":
		if n > 0 {
			S.Contract = (S.Contract + n - 1) % n
			S.ContractScroll = 0
		}
		return true
	case "PageDown":
		S.ContractScroll++
		return true
	case "PageUp":
		S.ContractScroll--
		return true
	case "Enter":
		if n == 0 || S.Contract >= n {
			return true
		}
		switch rows[S.Contract].kind {
		case "offer":
			// the row index is not the offer index: the story row sits first
			if idx := rows[S.Contract].idx; !g.Offers[idx].Taken {
				g.AcceptOffer(idx)
			}
		case "story":
			if !g.Cmd.Campaign.Active {
				g.AcceptCampaign()
			}
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
	drawData(fb, "FLYING "+cur.Name, 14, 30, ColYellow)
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
		drawData(fb, sc.Name, 18, y, col)
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
			if S.SlotConfirm && ok {
				label = itoa(i+1) + "  " + lang.T("ERASE?") + "  " + name
				col = ColRed
			}
		}
		fb.Text(clip(label, 38), 46, y, col, 1)
	}
	if S.SlotSave {
		Hint = "UP/DOWN   ENTER SAVE   DEL ERASE   ESC BACK"
	} else {
		Hint = "UP/DOWN   ENTER LOAD   DEL ERASE   ESC BACK"
	}
}

func slotsKey(k string) bool {
	g := game.G
	switch k {
	case "Escape":
		if S.SlotConfirm {
			S.SlotConfirm = false
			return true
		}
		if g.State == "title" {
			g.Screen = "launch"
			S.Title = 0
		} else {
			g.Screen = "launch"
		}
		return true
	case "ArrowUp":
		S.Slot = (S.Slot + game.Slots - 1) % game.Slots
		S.SlotConfirm = false
		return true
	case "ArrowDown":
		S.Slot = (S.Slot + 1) % game.Slots
		S.SlotConfirm = false
		return true
	case "Delete", "Backspace":
		// the same two steps: ask, then erase
		if !game.SlotExists(S.Slot) {
			g.Message("SLOT EMPTY", "bad")
			return true
		}
		if !S.SlotConfirm {
			S.SlotConfirm = true
			g.Message("PRESS DEL AGAIN TO ERASE THIS COMMANDER", "bad")
			return true
		}
		S.SlotConfirm = false
		if game.DeleteSlot(S.Slot) {
			g.Message("COMMANDER ERASED", "good")
		} else {
			g.Message("ERASE FAILED", "bad")
		}
		return true
	case "Enter":
		if S.SlotSave {
			// writing over somebody's commander is a decision, not a keystroke:
			// the first Enter asks, the second does it
			if game.SlotExists(S.Slot) && !S.SlotConfirm {
				S.SlotConfirm = true
				g.Message("PRESS ENTER AGAIN TO OVERWRITE", "bad")
				return true
			}
			S.SlotConfirm = false
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
	if g.State == "docked" {
		// a station sub-screen is not a place to come back to: the journal, the
		// market and the rest are pages, and coming back into a page looks
		// exactly like the game hanging
		S.SessionScreen = "launch"
	}
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
		// the first thing a new commander sees is who he is
		g.StartIntro()
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

// journalKey leaves the journal with Escape, Enter or the back button, and
// scrolls it if the page is longer than the panel.
func journalKey(k string) bool {
	g := game.G
	switch k {
	case "Escape", "Enter", "Backspace":
		g.Screen = "launch"
	case "ArrowUp":
		S.JournalScroll--
	case "ArrowDown":
		S.JournalScroll++
	case "PageUp":
		S.JournalScroll -= 8
	case "PageDown":
		S.JournalScroll += 8
	}
	if S.JournalScroll < 0 {
		S.JournalScroll = 0
	}
	return true
}

func journalScreen(fb *render.FB) {
	g := game.G
	plate(fb, 8, 16, 304, 182, "JOURNAL")
	// the page scrolls instead of overflowing the panel
	off := S.JournalScroll * 10
	y := 32 - off
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
		fb.Text("PRIVATE CONTRACT", 16, y, ColDim, 1)
		y += 11
		fb.Text(clip(ms, 46), 16, y, ColGreen, 1)
		y += 10
		fb.Text("PAYS "+itoa(g.Cmd.Mission.Pay)+" CR", 16, y, ColDim, 1)
		y += 12
	} else {
		fb.Text("PRIVATE CONTRACT", 16, y, ColDim, 1)
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

// ------------------------------------------------------------------- help

// helpPages is the handbook: four topics, one screen each. Everything a first
// flight needs and nothing else.
var helpPages = []struct {
	Title string
	Lines []string
}{
	{"DOCKING", []string{
		"THE STATION HAS A SLOT IN ONE FACE.",
		"FLY AT IT WITH THE NOSE POINTED IN.",
		"SPEED MUST BE UNDER 85% OR YOU HIT",
		"THE HULL INSTEAD OF THE MOUTH.",
		"",
		"FIT A DOCKING COMPUTER (L) AND THE",
		"SHIP FLIES THE LAST STRETCH ITSELF.",
		"",
		"THE STATION TURNS: THE SLOT MOVES.",
	}},
	{"FIGHTING", []string{
		"SPACE FIRES THE LASERS.",
		"M LAUNCHES A MISSILE AT THE CROSSHAIR.",
		"E IS E.C.M. AND KILLS INCOMING MISSILES.",
		"TAB IS THE ENERGY BOMB: EVERYTHING",
		"CLOSE TO YOU DIES, INCLUDING YOUR SHIELDS.",
		"",
		"WHAT A KILL LEAVES FLOATS: CANISTERS",
		"OF CARGO, MARKED YELLOW ON THE RADAR.",
		"A FUEL SCOOP (S) SWALLOWS THEM. WITHOUT",
		"ONE THEY DAMAGE THE HULL ON IMPACT.",
		"",
		"WATCH THE CABIN TEMPERATURE NEAR THE STAR.",
		"SHIELDS RECHARGE, THE HULL DOES NOT.",
	}},
	{"TRAVEL", []string{
		"C OPENS THE GALACTIC CHART.",
		"PICK A SYSTEM WITHIN 7 LIGHT YEARS",
		"AND PRESS H. FUEL BURNS WITH THE JUMP.",
		"",
		"J IS A SHORT JUMP INSIDE THE SYSTEM.",
		"G CROSSES TO THE NEXT GALAXY: ONE WAY,",
		"AND YOU NEED A GALACTIC HYPERDRIVE.",
		"",
		"F BUYS FUEL AT ANY STATION.",
	}},
	{"TRADING", []string{
		"BUY WHERE IT IS CHEAP, SELL WHERE IT IS DEAR.",
		"B BUYS A TONNE, SHIFT+B BUYS TEN.",
		"S SELLS A TONNE, SHIFT+S SELLS TEN.",
		"",
		"THE BOARD PAYS FOR COURIER RUNS,",
		"PASSENGERS, MAIL AND BOUNTIES.",
		"JOBS HAVE DEADLINES: A MISSED ONE COSTS.",
		"",
		"TO WIN A FIGHT, BUY A BIGGER GUN FIRST.",
	}},
}

func helpScreen(fb *render.FB) {
	page := S.Help
	if page < 0 || page >= len(helpPages) {
		page = 0
	}
	p := helpPages[page]
	plate(fb, 12, 20, 296, 180, p.Title)
	for i, line := range p.Lines {
		col := ColText
		if line == "" {
			continue
		}
		fb.Text(line, 24, 42+i*13, col, 1)
	}
	// page marks
	for i := range helpPages {
		x := 140 + i*12
		c := render.RGB(0x25405a)
		if i == page {
			c = ColYellow
		}
		fb.FillRect(x, 186, x+6, 191, c)
	}
	Hint = "LEFT/RIGHT PAGE   ESC BACK"
}

func helpKey(k string) bool {
	switch k {
	case "Escape":
		if game.G.State == "title" {
			game.G.Screen = ""
		} else if S.PrevScreen != "" {
			game.G.Screen = S.PrevScreen
			S.PrevScreen = ""
		} else {
			game.G.Screen = "launch"
		}
		return true
	case "ArrowLeft":
		S.Help = (S.Help + len(helpPages) - 1) % len(helpPages)
		return true
	case "ArrowRight":
		S.Help = (S.Help + 1) % len(helpPages)
		return true
	}
	return true
}

// ---------------------------------------------------------------- log

// logFilters are the tabs of the log screen, in order.
var logFilters = []string{"ALL", "TRADE", "FIGHT", "JOBS", "TRAVEL"}

// logLine turns an entry into the sentence the player reads. The words are
// translated one by one and the data (names, numbers) is left alone, so the
// log reads in the player's language without a second copy of every line.
func logLine(e game.LogEntry) string {
	what := lang.T(e.A)
	switch e.Kind {
	case game.LogKill:
		return lang.T("KILLED") + " " + e.A
	case game.LogBuy:
		return lang.T("BOUGHT") + " " + what + " X" + itoa(e.N) + " @" + itoa(perUnit(e))
	case game.LogSell:
		return lang.T("SOLD") + " " + what + " X" + itoa(e.N) + " @" + itoa(perUnit(e))
	case game.LogFuel:
		return lang.T("REFUELLED") + " " + itoa(e.N) + " LY"
	case game.LogEquip:
		return lang.T("FITTED") + " " + what
	case game.LogShip:
		return lang.T("NEW SHIP") + " " + e.A
	case game.LogJob:
		return lang.T("CONTRACT") + " " + e.A
	case game.LogDone:
		return lang.T("PAID") + " " + e.A
	case game.LogFail:
		return lang.T("FAILED") + " " + e.A
	case game.LogJump:
		return lang.T("JUMPED TO") + " " + e.A
	case game.LogMine:
		return lang.T("MINED") + " " + what
	case game.LogDeath:
		return lang.T("LOST SHIP") + " " + lang.T(e.A)
	case game.LogVisit:
		return lang.T("FIRST VISIT") + " " + e.A
	case game.LogPirate:
		return lang.T("ATTACKED BY PIRATES")
	case game.LogPolice:
		return lang.T("STOPPED BY POLICE")
	case game.LogSwarm:
		return lang.T("ALIEN SWARM")
	}
	return e.A
}

// perUnit is the price a single tonne went for: the total is already on the
// right of the row, and the price is the number a trader actually compares.
func perUnit(e game.LogEntry) int {
	if e.N <= 0 || e.CR == 0 {
		return 0
	}
	v := e.CR / float64(e.N)
	if v < 0 {
		v = -v
	}
	return int(v + 0.5)
}

// logVisible returns the entries that match the current filter, newest first.
func logVisible() []game.LogEntry {
	g := game.G
	out := make([]game.LogEntry, 0, len(g.Cmd.Log))
	for i := len(g.Cmd.Log) - 1; i >= 0; i-- {
		e := g.Cmd.Log[i]
		if S.LogFilter > 0 && game.LogGroup(e.Kind) != logFilters[S.LogFilter] {
			continue
		}
		out = append(out, e)
	}
	return out
}

func logScreen(fb *render.FB) {
	g := game.G
	plate(fb, 6, 14, 310, 190, "LOG")
	// filter tabs
	x := 12
	for i, name := range logFilters {
		col := ColDim
		if i == S.LogFilter {
			col = ColYellow
			fb.FillRect(x-2, 24, x+render.TextWidth(name, 1)+2, 31, render.RGB(0x12384f))
		}
		fb.Text(name, x, 25, col, 1)
		x += render.TextWidth(name, 1) + 8
	}
	ents := logVisible()
	if len(ents) == 0 {
		fb.Text(lang.T("THE LOG IS EMPTY"), 14, 44, ColText, 1)
		fb.Text(lang.T("EVERY TRADE, KILL AND JUMP IS WRITTEN HERE"), 14, 56, ColDim, 1)
		Hint = "LEFT/RIGHT FILTER   ESC BACK"
		if g.Cmd.Flags["HAS_LIST"] || g.Cmd.Flags["READ_LIST"] {
			Hint = "LEFT/RIGHT FILTER   L CLOSING NOTICE   ESC BACK"
		}
		return
	}
	if S.Log < 0 {
		S.Log = 0
	}
	if S.Log > len(ents)-1 {
		S.Log = len(ents) - 1
	}
	const perPage = 13
	if S.Log > len(ents)-1 {
		S.Log = len(ents) - 1
	}
	// keep the cursor on screen: the list scrolls under a fixed window
	top := S.Log - perPage/2
	if top < 0 {
		top = 0
	}
	if top > len(ents)-perPage {
		top = len(ents) - perPage
	}
	if top < 0 {
		top = 0
	}
	// column headers, so a row reads as day / what / where / money
	fb.Text("DAY", 8, 34, render.RGB(0x2d5f7a), 1)
	fb.Text("WHAT", 32, 34, render.RGB(0x2d5f7a), 1)
	fb.Text("WHERE", 196, 34, render.RGB(0x2d5f7a), 1)
	fb.TextR("CR", 306, 34, render.RGB(0x2d5f7a), 1)
	fb.HLine(8, 306, 41, render.RGB(0x1b3348))

	y := 44
	for i := top; i < len(ents) && i < top+perPage; i++ {
		e := ents[i]
		sel := i == S.Log
		if sel {
			fb.FillRect(10, y-1, 306, y+7, render.RGB(0x12384f))
		}
		day := "D" + itoa(e.Day)
		fb.Text(day, 8, y, ColDim, 1)
		line := clip(logLine(e), 26)
		col := ColText
		if sel {
			col = ColYellow
		}
		fb.TextRaw(line, 32, y, col, 1)
		// where it happened: the whole point of a log is knowing the station
		where := e.Sys
		if where == "" {
			where = "-"
		}
		fb.TextRaw(clip(where, 9), 196, y, ColDim, 1)
		if e.CR != 0 {
			money := "+" + itoa(int(e.CR))
			c := ColGreen
			if e.CR < 0 {
				money = itoa(int(e.CR))
				c = ColRed
			}
			fb.TextR(money, 306, y, c, 1)
		}
		y += 11
	}
	Hint = "UP/DOWN SCROLL   LEFT/RIGHT FILTER   ESC BACK"
}

// closingScreen is the notice the courier sold or gave: eighteen stations and
// nothing to do about it. No buttons, no actions - that is deliberate.
func closingScreen(fb *render.FB) {
	g := game.G
	plate(fb, 6, 14, 310, 190, lang.T("STATIONS TO BE CLOSED"))
	list := g.ClosingStations()
	if len(list) == 0 {
		fb.Text(lang.T("NOTHING"), 14, 30, ColDim, 1)
		Hint = "ESC BACK"
		return
	}
	const perPage = 9
	if S.ClosingScroll < 0 {
		S.ClosingScroll = 0
	}
	if S.ClosingScroll > len(list)-1 {
		S.ClosingScroll = len(list) - 1
	}
	y := 32
	for i := S.ClosingScroll; i < len(list) && i < S.ClosingScroll+perPage; i++ {
		fb.Text(itoa(i+1), 14, y, ColDim, 1)
		fb.TextRaw(list[i], 32, y, ColText, 1)
		y += 12
	}
	fb.TextRaw(lang.T("CLOSING NOTICE. NO ACTION CAN BE TAKEN."), 12, 176, ColDim, 1)
	Hint = "UP/DOWN   ESC BACK"
}

func closingKey(k string) bool {
	switch k {
	case "Escape":
		game.G.Screen = "log"
	case "ArrowUp":
		S.ClosingScroll--
	case "ArrowDown":
		S.ClosingScroll++
	case "PageUp":
		S.ClosingScroll -= 9
	case "PageDown":
		S.ClosingScroll += 9
	}
	return true
}

// endingScreen is the last page: what the sector looks like once the pilot holds
// the seal, or does not.
func endingScreen(fb *render.FB) {
	g := game.G
	plate(fb, 6, 14, 310, 190, "")
	title, body, last := endingText(g.Cmd.Ending)
	y := 34
	for _, ln := range wrapText(title, 46) {
		fb.TextRaw(ln, 16, y, ColYellow, 1)
		y += 12
	}
	y += 8
	for _, ln := range wrapText(body, 48) {
		fb.TextRaw(ln, 16, y, ColText, 1)
		y += 10
	}
	for _, ln := range wrapText(last, 48) {
		fb.TextRaw(ln, 16, 168, ColGreen, 1)
		y += 10
	}
	Hint = "ENTER MAIN MENU"
}

func endingText(kind string) (title, body, last string) {
	switch kind {
	case "chaos":
		return lang.T("CHAOS"),
			lang.T("THERE ARE NO LICENCES IN THE SECTOR, NO COMBINE AND NO ORDER. PRICES ON ONE STATION ARE TWICE THOSE ON THE NEXT AND NOBODY KNOWS WHY. CARRYING CARGO PAYS AND IS FRIGHTENING. WITHIN A YEAR SOMEONE PROPOSES THAT EVERYONE AGREE ON THE RULES - AND IT IS NOT YOU."),
			lang.T("TRADE RESUMED. NO LICENCES REQUIRED.")
	case "council":
		return lang.T("COUNCIL"),
			lang.T("THE STATIONS TAKE THE SEAL AND LOCK IT IN A SHARED SAFE WITH EIGHTEEN KEYS. RULES ARE WRITTEN BY THE PEOPLE WHO LIVE UNDER THEM. THE COMBINE EXISTS AS A DEBTOR, AND EVERY MEETING OPENS WITH ITS ACCOUNTS READ ALOUD. IT IS DULL, AND IT WORKS."),
			lang.T("QUORUM PRESENT. THE FIRST QUESTION IS OPEN.")
	}
	return lang.T("COMBINE"),
		lang.T("A MONTH LATER YOU SIT IN THE SAME CHAIR WITH THE SAME BOOK AND SIGN THE FIRST PAPER. IT CLOSES THREE STATIONS WHOSE DEBT IS NOT SERVICED. YOU SIGN IT BECAUSE THE BALANCE HAS TO BE SETTLED, AND BECAUSE NOW IT IS YOUR BALANCE."),
		lang.T("YOUR LICENCE IS VALID. PLEASE CONTINUE.")
}

func endingKey(k string) bool {
	switch k {
	case "Enter", "Escape", " ":
		QuitToTitle()
	}
	return true
}

func logKey(k string) bool {
	g := game.G
	n := len(logVisible())
	switch k {
	case "Escape":
		g.Screen = "launch"
	case "l", "L":
		// the list is a page of the log, not a menu item: it is not a thing to do
		if g.Cmd.Flags["HAS_LIST"] || g.Cmd.Flags["READ_LIST"] {
			g.Screen = "closing"
			S.ClosingScroll = 0
		}
	case "ArrowUp":
		S.Log--
		if S.Log < 0 {
			S.Log = 0
		}
	case "ArrowDown":
		S.Log++
		if n > 0 && S.Log > n-1 {
			S.Log = n - 1
		}
	case "PageUp", "PageDown":
		if k == "PageUp" {
			S.Log -= 13
		} else {
			S.Log += 13
		}
		if S.Log < 0 {
			S.Log = 0
		}
		if n > 0 && S.Log > n-1 {
			S.Log = n - 1
		}
	case "ArrowLeft", "ArrowRight":
		S.LogFilter = (S.LogFilter + len(logFilters) + 1) % len(logFilters)
		if k == "ArrowLeft" {
			S.LogFilter = (S.LogFilter + len(logFilters) - 2) % len(logFilters)
		}
		S.Log = 0
	}
	return true
}

// ---------------------------------------------------------------- talk

// wrapText breaks a line into 46-character rows, by words. The dialog window
// shows six of them at a time and scrolls, which is what lets a speech be a
// paragraph instead of a telegram.
func wrapText(src string, width int) []string {
	var out []string
	for _, para := range strings.Split(src, "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}
		line := ""
		for _, w := range words {
			if line == "" {
				line = w
				continue
			}
			// width is in characters: Russian is two bytes each, and counting
			// bytes wrapped every line at half the panel
			if utf8.RuneCountInString(line)+1+utf8.RuneCountInString(w) <= width {
				line += " " + w
			} else {
				out = append(out, line)
				line = w
			}
		}
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// talkScreen is the conversation: six lines of window, a scroll indicator and
// the options pinned at the bottom.
func talkScreen(fb *render.FB) {
	g := game.G
	n := game.DialogueNode(g.Talk)
	if n == nil {
		g.EndTalk()
		return
	}
	plate(fb, 6, 14, 310, 190, "")
	if n.Speaker != "" {
		fb.FillRect(8, 16, 308, 26, render.RGB(0x08131f))
		fb.Text(lang.T(n.Speaker), 12, 19, ColYellow, 1)
	}
	// the answers are pinned to the bottom and the text gets whatever is left:
	// a fixed six-line window left half the panel empty and cut the last line
	opts := g.VisibleOptions(n)
	const lineH = 10
	optsTop := 186 - len(opts)*11
	textTop := 32
	visible := (optsTop - 6 - textTop) / lineH
	if visible < 3 {
		visible = 3
	}

	var rows []string
	for _, ln := range g.TalkText(n) {
		rows = append(rows, wrapText(ln, 48)...)
		rows = append(rows, "") // a gap between paragraphs
	}
	maxTop := len(rows) - visible
	if maxTop < 0 {
		maxTop = 0
	}
	if S.TalkScroll > maxTop {
		S.TalkScroll = maxTop
	}
	if S.TalkScroll < 0 {
		S.TalkScroll = 0
	}
	y := textTop
	for i := S.TalkScroll; i < len(rows) && i < S.TalkScroll+visible; i++ {
		fb.TextRaw(rows[i], 12, y, ColText, 1)
		y += lineH
	}
	scrolling := len(rows) > visible
	if scrolling {
		fb.TextR(itoa(S.TalkScroll+1)+"/"+itoa(maxTop+1), 304, optsTop, ColDim, 1)
	}
	for i, idx := range opts {
		oy := optsTop + i*11
		col := ColYellow
		if i == S.TalkOpt {
			fb.FillRect(10, oy-2, 306, oy+8, render.RGB(0x12384f))
		}
		fb.Text(itoa(i+1)+"  "+lang.T(n.Options[idx].Text), 14, oy, col, 1)
	}
	Hint = "UP/DOWN SELECT   ENTER CHOOSE   ESC BACK"
	if scrolling {
		Hint = "PGUP/PGDN SCROLL   " + Hint
	}
}

func talkKey(k string) bool {
	g := game.G
	n := game.DialogueNode(g.Talk)
	if n == nil {
		g.EndTalk()
		return true
	}
	opts := g.VisibleOptions(n)
	if S.TalkOpt >= len(opts) {
		S.TalkOpt = 0
	}
	// the same maths as the screen: answers pinned low, text above them
	const lineH = 10
	visible := (186 - len(opts)*11 - 6 - 32) / lineH
	if visible < 3 {
		visible = 3
	}
	rows := 0
	for _, ln := range g.TalkText(n) {
		rows += len(wrapText(ln, 48)) + 1
	}
	maxTop := rows - visible
	if maxTop < 0 {
		maxTop = 0
	}
	choose := func(i int) {
		if i < 0 || i >= len(opts) {
			return
		}
		// TalkChoose counts visible answers, not the raw option index
		g.TalkChoose(i)
		S.TalkScroll = 0
		S.TalkOpt = 0
	}
	switch k {
	case "Escape":
		g.EndTalk()
	case "ArrowUp":
		// arrows walk the answers; the text scrolls with the page keys
		if len(opts) > 0 {
			S.TalkOpt = (S.TalkOpt + len(opts) - 1) % len(opts)
		}
	case "ArrowDown":
		if len(opts) > 0 {
			S.TalkOpt = (S.TalkOpt + 1) % len(opts)
		}
	case "Enter", " ":
		choose(S.TalkOpt)
	case "PageUp":
		S.TalkScroll -= visible
		if S.TalkScroll < 0 {
			S.TalkScroll = 0
		}
	case "PageDown":
		S.TalkScroll += visible
		if S.TalkScroll > maxTop {
			S.TalkScroll = maxTop
		}
	default:
		if len(k) == 1 && k[0] >= '1' && k[0] <= '9' {
			choose(int(k[0] - '1'))
		}
	}
	return true
}

// marketRows is how many lines the market list has: the commodities plus the
// crate when one is aboard.

// marketMax is the most the pilot can buy or sell right now: the stock for a
// purchase (also limited by the hold and the wallet), what is aboard for a sale.
func marketMax(idx int, sell bool) int {
	g := game.G
	if idx < 0 || idx >= len(g.Market) {
		return 0
	}
	m := g.Market[idx]
	if sell {
		return g.Cmd.Cargo[idx]
	}
	max := m.Qty
	if free := g.Cmd.CargoCap - game.CargoUsed(g.Cmd); free < max {
		max = free
	}
	price := m.Price
	if g.LicensedHere() {
		price = int(float64(price) * 0.95)
	}
	if price > 0 {
		if afford := int(g.Cmd.Cash) / price; afford < max {
			max = afford
		}
	}
	if max < 0 {
		max = 0
	}
	return max
}

func marketRows() int {
	n := len(game.G.Market)
	if game.G.Cmd.Boxes > 0 {
		n++
	}
	return n
}
