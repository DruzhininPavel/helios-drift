package game

// The story, as data. Line keys are English; the Russian lives in lang/ru.go,
// so a node is one entry here and one block of translations there.

func init() {
	registerNodes(
		// ---- chapter one: the crate
		&Node{
			ID:      "c1_disp",
			Kind:    "talk",
			Speaker: "DISPATCHER",
			Lines: []string{
				"FOUR TONNES, AS ORDERED. ALL ACCOUNTED FOR, I CHECKED.",
				"YOU WILL BE PAID ON ACCEPTANCE. NOT BY US - BY THE COMPANY. WE HAVE NOTHING TO PAY WITH: THE LAST PAYMENT WAS RECALLED, NO ACCOUNTS OPEN SINCE APRIL.",
				"BUT THERE IS A CRATE. IT IS VALUED HIGHER THAN YOUR WHOLE RUN. TAKE IT OFF OUR HANDS, SELL IT AT A FREE STATION AND FORGET US.",
				"AND DO NOT ASK WHAT IS INSIDE. I AM SERIOUS.",
				"...FINE. YOU ARE THE PILOT, IT IS YOUR CALL. THREE RUNS LIKE THIS IN A MONTH. ALL THREE LEFT WITH A CRATE.",
			},
			Options: []Option{
				{
					Text: "TAKE THE CRATE",
					Next: "",
					Do: []Effect{
						{Kind: FxBox, N: 1},
						{Kind: FxPapers, N: 1},
						{Kind: FxCash, N: 1200},
						{Kind: FxFlag, Key: "HAS_SEAL"},
						{Kind: FxQuest},
					},
				},
				{
					Text: "REFUSE",
					Next: "",
					Do: []Effect{
						{Kind: FxFlag, Key: "REFUSED_SEAL"},
						{Kind: FxCash, N: 1200},
						{Kind: FxQuest},
					},
				},
				{
					Text: "DEMAND A RECEIPT",
					Next: "",
					Do: []Effect{
						{Kind: FxBox, N: 1},
						{Kind: FxPapers, N: 2},
						{Kind: FxFlag, Key: "HAS_RECEIPT"},
						{Kind: FxCash, N: 1200},
						{Kind: FxQuest},
					},
				},
			},
		},
		// ---- chapter two: three ships, one question
		&Node{
			ID:      "c2_cmdA",
			Kind:    "talk",
			Speaker: "COLLECTOR LEADER",
			Lines: []string{
				"YOU BROUGHT FOOD. WE WILL TAKE IT.",
				"NO HARD FEELINGS, IT IS NOT PERSONAL. WE WERE HIRED FOR THE CARGO, NOT FOR YOU.",
				"YOU CAN FIGHT - AND WE WILL SELL YOU FOR SCRAP TOO. OR YOU CAN PAY, AND WE PART PRETTY.",
				"THREE THOUSAND. CHEAPER THAN A HULL.",
			},
			Options: []Option{
				{
					Text: "PAY 3000",
					Next: "c2_paid",
					If:   []Cond{{Key: "cash", Op: ">=", N: 3000}, {Key: "HAS_RECEIPT", Op: "!flag"}},
					Do: []Effect{
						{Kind: FxCash, N: -3000},
						{Kind: FxPapers, N: 1},
						{Kind: FxFlag, Key: "PAID_PERMIT"},
					},
				},
				{
					Text: "PAY 4500",
					Next: "c2_paid",
					If:   []Cond{{Key: "cash", Op: ">=", N: 4500}, {Key: "HAS_RECEIPT", Op: "flag"}},
					Do: []Effect{
						{Kind: FxCash, N: -4500},
						{Kind: FxPapers, N: 1},
						{Kind: FxFlag, Key: "PAID_PERMIT"},
					},
				},
				{
					Text: "FIGHT THEM",
					Next: "",
					Do:   []Effect{{Kind: FxSpawn, Key: "AMBUSH"}},
				},
				{
					Text: "JUMP OUT",
					Next: "c2_jump",
					Do: []Effect{
						{Kind: FxSpawn, Key: "AMBUSH"},
						{Kind: FxFlag, Key: "DISP_OUT"},
						{Kind: FxQuest},
					},
				},
			},
		},
		&Node{
			ID:      "c2_cmdAR",
			Kind:    "talk",
			Speaker: "COLLECTOR LEADER",
			Lines: []string{
				"YOU BROUGHT FOOD. WE WILL TAKE IT.",
				"AND YOU HAVE A RECEIPT FOR THAT CRATE. THE CUSTODIAN PAYS US BY THE PAPER, SO NOW IT COSTS YOU MORE.",
				"FOUR AND A HALF THOUSAND, AND WE PRETEND WE NEVER MET.",
			},
			Options: []Option{
				{
					Text: "PAY 4500",
					Next: "c2_paid",
					If:   []Cond{{Key: "cash", Op: ">=", N: 4500}},
					Do: []Effect{
						{Kind: FxCash, N: -4500},
						{Kind: FxPapers, N: 1},
						{Kind: FxFlag, Key: "PAID_PERMIT"},
					},
				},
				{
					Text: "FIGHT THEM",
					Next: "",
					Do:   []Effect{{Kind: FxSpawn, Key: "AMBUSH"}},
				},
				{
					Text: "JUMP OUT",
					Next: "c2_jump",
					Do: []Effect{
						{Kind: FxSpawn, Key: "AMBUSH"},
						{Kind: FxFlag, Key: "DISP_OUT"},
						{Kind: FxQuest},
					},
				},
			},
		},
		&Node{
			ID:      "c2_cmdB",
			Kind:    "talk",
			Speaker: "COLLECTOR LEADER",
			Lines: []string{
				"NO CRATE. YOU LEFT IT AT THE CLINIC AND LEFT THEM TO US.",
				"SO WE WAITED FOR YOU INSTEAD. THE CLINIC STILL OWES, AND SOMEONE HAS TO PAY.",
				"THREE THOUSAND, OR WE SETTLE IT THE OTHER WAY.",
			},
			Options: []Option{
				{
					Text: "PAY 3000",
					Next: "c2_paid",
					If:   []Cond{{Key: "cash", Op: ">=", N: 3000}},
					Do: []Effect{
						{Kind: FxCash, N: -3000},
						{Kind: FxPapers, N: 1},
						{Kind: FxFlag, Key: "PAID_CLINIC"},
					},
				},
				{
					Text: "FIGHT THEM",
					Next: "",
					Do:   []Effect{{Kind: FxSpawn, Key: "AMBUSH"}},
				},
				{
					Text: "JUMP OUT",
					Next: "c2_jump",
					Do: []Effect{
						{Kind: FxSpawn, Key: "AMBUSH"},
						{Kind: FxFlag, Key: "DISP_OUT"},
						{Kind: FxQuest},
					},
				},
			},
		},
		&Node{
			ID:      "c2_paid",
			Kind:    "talk",
			Speaker: "COLLECTOR LEADER",
			Lines: []string{
				"CLEAN. WE NEVER MET.",
				"IF THE CUSTODIAN ASKS - YOU PAID FOR A PERMIT, NOT FOR US.",
			},
			Options: []Option{
				{Text: "LEAVE", Next: "", Do: []Effect{{Kind: FxQuest}}},
			},
		},
		&Node{
			ID:      "c2_jump",
			Kind:    "talk",
			Speaker: "COMPANY",
			Lines: []string{
				"THE RUN WAS NOT DELIVERED. THE DEBT GOES TO THE GUARANTOR.",
				"THE GUARANTOR IS THE ONE WHO ACCEPTED THE CARGO. WE SUGGEST YOU DO NOT RETURN TO THIS SYSTEM.",
			},
			Options: []Option{
				{Text: "ACKNOWLEDGE", Next: ""},
			},
		},
		// ---- chapter three: the wreck, the recorder and the menu
		&Node{
			ID:      "c3_roster",
			Kind:    "screen",
			Speaker: "CARGO ROSTER",
			Lines: []string{
				"THE HOLD IS BARE. FORTY RACKS, ALL EMPTY.",
				"ON THE BULKHEAD: A BOARD WITH A COLUMN OF NAMES AND A COLUMN OF NUMBERS.",
				"PICKED UP: 37    PROCESSED: 37    WAITING: 0",
				"VACANCIES: 3",
			},
			Variants: []LineVariant{
				{If: []Cond{{Key: "KILLED_COLLECTORS", Op: "flag"}}, Lines: []string{
					"THE HOLD IS BARE. FORTY RACKS, ALL EMPTY.",
					"ON THE BULKHEAD: A BOARD WITH A COLUMN OF NAMES AND A COLUMN OF NUMBERS.",
					"PICKED UP: 37    PROCESSED: 37    WAITING: 0    VACANCIES: 3",
					"THREE NAMES ON THE BOARD MATCH THE SHIPS YOU BURNED. NOBODY HAS CROSSED THEM OUT YET.",
				}},
				{If: []Cond{{Key: "REFUSED_SEAL", Op: "flag"}}, Lines: []string{
					"THE HOLD IS BARE. FORTY RACKS, ALL EMPTY.",
					"ON THE BULKHEAD: A BOARD WITH A COLUMN OF NAMES AND A COLUMN OF NUMBERS.",
					"PICKED UP: 37    PROCESSED: 37    WAITING: 0    VACANCIES: 3",
					"THE CRATE YOU LEFT AT THE CLINIC IS LISTED HERE AS A DELIVERY, NOT AS A LOSS.",
				}},
			},
			Options: []Option{{Text: "CONTINUE", Next: "c3_recorder"}},
		},
		&Node{
			ID:      "c3_recorder",
			Kind:    "talk",
			Speaker: "FLIGHT RECORDER",
			Lines: []string{
				"THE RECORDER BEEPS EVERY EIGHT SECONDS. EMERGENCY POWER STILL HOLDS.",
				"FIFTEEN MINUTES BACK THERE IS A VOICE: JUMP COORDINATES GO OUT EVERY SIX DAYS, THE CARGO IS TAKEN AT THE HIVE, THE LIST IS IN THE APPENDIX.",
				"THE COMPANY WANTS THIS BURNED. IT DOES NOT SAY WHY, AND IT DOES NOT SAY NOT TO READ IT.",
			},
			Options: []Option{
				{Text: "COPY THE RECORDING", Next: "c3_manifest", Do: []Effect{{Kind: FxFlag, Key: "COPIED_RECORDER"}}},
				{Text: "BURN IT", Next: "c3_manifest", Do: []Effect{{Kind: FxFlag, Key: "BURNED_RECORDER"}}},
			},
		},
		&Node{
			ID:      "c3_manifest",
			Kind:    "screen",
			Speaker: "APPENDIX",
			Lines: []string{
				"SUPPLY - PORTIONS - TASTE",
				"QUIET HARBOUR - 40000 - REGULAR",
				"FREE TRADE - 120000 - HOTEL",
				"NINTH DOCK - 8400 - SWEET",
				"IT IS NOT A MANIFEST. IT IS A MENU, AND YOU ARE READING THE PORTIONS COLUMN.",
			},
			Options: []Option{{Text: "CONTINUE", Next: "c3_holo"}},
		},
		&Node{
			ID:      "c3_holo",
			Kind:    "talk",
			Speaker: "COMPANY",
			Lines: []string{
				"THIS IS PRIVATE PROPERTY AND THE ITEM IS SCHEDULED FOR DISPOSAL.",
				"YOU MAY ASK ONE QUESTION.",
			},
			Options: []Option{
				{Text: "WHAT IS IN THE CRATE", Next: "c3_ans_box", Do: []Effect{{Kind: FxFlag, Key: "KNOWS_BOX"}}},
				{Text: "WHAT IS THE HIVE", Next: "c3_ans_hive", Do: []Effect{{Kind: FxFlag, Key: "KNOWS_HIVE"}}},
				{Text: "WHAT DO I GET PAID", Next: "c3_ans_price", Do: []Effect{{Kind: FxFlag, Key: "ASKED_PRICE"}}},
			},
		},
		&Node{
			ID:      "c3_ans_box",
			Kind:    "talk",
			Speaker: "COMPANY",
			Lines:   []string{"A STATION OBLIGATION, ASSIGNED TO THE HOLDER."},
			Options: []Option{{Text: "CONTINUE", Next: "c3_end"}},
		},
		&Node{
			ID:      "c3_ans_hive",
			Kind:    "talk",
			Speaker: "COMPANY",
			Lines:   []string{"A RECEIVING POINT. STATUS: CONFIDENTIAL."},
			Options: []Option{{Text: "CONTINUE", Next: "c3_end"}},
		},
		&Node{
			ID:      "c3_ans_price",
			Kind:    "talk",
			Speaker: "COMPANY",
			Lines:   []string{"THE AMOUNT IS IN YOUR CONTRACT. YOU SIGNED IT."},
			Options: []Option{{Text: "CONTINUE", Next: "c3_end"}},
		},
		&Node{
			ID:      "c3_end",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE RECORDER IS STILL BEEPING. THE WRECK HOLDS STATION.",
				"BURN IT AND LEAVE.",
			},
			Options: []Option{{
				Text: "FIRE ON THE WRECK",
				Next: "",
				Do:   []Effect{{Kind: FxSpawn, Key: "HULK"}},
			}},
		},
		// ---- chapter six: framed, and the clerk with the cabinet
		&Node{
			ID:      "c6_brief",
			Kind:    "talk",
			Speaker: "",
			Lines: []string{
				"YOUR NAME IS IN THE REPORT. THREE FREIGHTERS, TWENTY-SIX CREW, INTENT.",
				"TWO NAVY VIPERS ARE ON YOUR TAIL AND YOUR RECORD IS ALREADY DIRTY.",
				"THEY WILL ASK ONCE, BY PROTOCOL. THE ANSWER IS YOURS.",
			},
			Variants: []LineVariant{
				{If: []Cond{{Key: "CONVOY_SPARED", Op: "flag"}}, Lines: []string{
					"THE CONVOY CONFIRMED YOU DID NOT ATTACK THEM.",
					"THE NAVY IS NOT AFTER YOU. IT IS AFTER THE PAPERS.",
					"YOU ARE SUMMONED TO THE ARCHIVE. ATTEND.",
				}},
			},
			Options: []Option{
				{Text: "SHOW THE PAPERS", Next: "c6_papers",
					If: []Cond{{Key: "papers", Op: ">=", N: 1}}},
				{Text: "SURRENDER", Next: "c6_surrendered"},
				{Text: "FIGHT THEM", Next: "c6_fight"},
			},
		},
		&Node{
			ID:      "c6_papers",
			Kind:    "talk",
			Speaker: "NAVY",
			Lines: []string{
				"THESE ARE NOT PAPERS. THESE ARE RECEIPTS.",
				"BUT THE SHIP IN YOUR LOG MATCHES A CRATE THAT WENT OUT SIX DAYS AGO.",
				"GO TO THE ARCHIVE. ASK FOR THE SECOND DRAWER.",
			},
			Options: []Option{{Text: "CONTINUE", Next: "c6_archive"}},
		},
		&Node{
			ID:      "c6_archive",
			Kind:    "screen",
			Speaker: "CLERK",
			Lines: []string{
				"LISTEN, I DO NOT CARE WHAT YOU HAD ON YOU. I NEED A FORM.",
				"WITHOUT A FORM I CANNOT OPEN THE CABINET. WITH A FORM I CAN, BUT THE CABINET IS LOCKED AND THE KEY IS WITH THE CHIEF.",
				"I CAN WRITE THE FORM FOR YOU. IT IS NOT STRICTLY LEGAL. THREE HUNDRED.",
			},
			Options: []Option{
				{Text: "PAY 300", Next: "c6_arch2",
					If: []Cond{{Key: "cash", Op: ">=", N: 300}},
					Do: []Effect{{Kind: FxCash, N: -300}, {Kind: FxFlag, Key: "CLERK_1"}}},
				{Text: "BRAKE THE CABINET", Next: "c6_arch2",
					If: []Cond{{Key: "rig", Op: "flag"}},
					Do: []Effect{{Kind: FxFlag, Key: "CLERK_1_ANGRY"}, {Kind: FxWanted, N: 60}}},
				{Text: "NO", Next: "c6_arch2"},
			},
		},
		// the second archive: a clerk who is afraid of the combine, not of you
		&Node{
			ID:      "c6_arch2",
			Kind:    "screen",
			Speaker: "CLERK",
			Lines: []string{
				"I CANNOT. DO YOU UNDERSTAND? NOT WILL NOT. CANNOT.",
				"IF I GIVE YOU THIS, THEY CLOSE THE DOCK, AND THERE ARE TWO HUNDRED PEOPLE IN THE DOCK, FORTY OF THEM CHILDREN. MY SISTER WORKS THERE.",
				"BUT THERE IS ANOTHER WAY. SURETY. THE ONE WHO TOOK THE CARGO IN HIS OWN NAME. IF YOU HAVE HIS SIGNATURE I OPEN NOTHING AT ALL. YOU JUST COME AND TAKE WHAT IS YOURS.",
			},
			Options: []Option{
				{Text: "SHOW THE RECEIPT", Next: "c6_arch3",
					If: []Cond{{Key: "HAS_RECEIPT", Op: "flag"}},
					Do: []Effect{{Kind: FxFlag, Key: "SURETY"}}},
				{Text: "BRING THE COURIER", Next: "c6_arch3",
					If: []Cond{{Key: "COURIER_ALIVE", Op: "flag"}},
					Do: []Effect{{Kind: FxFlag, Key: "WITNESS"}}},
				{Text: "PLAY THE RECORDING", Next: "c6_arch3",
					If: []Cond{{Key: "COPIED_RECORDER", Op: "flag"}},
					Do: []Effect{{Kind: FxFlag, Key: "EVIDENCE"}}},
				{Text: "BRIBE HIM - 2000", Next: "c6_arch3",
					If: []Cond{{Key: "cash", Op: ">=", N: 2000}},
					Do: []Effect{{Kind: FxCash, N: -2000}, {Kind: FxFlag, Key: "CLERK_2"}}},
				{Text: "FORGE THE SIGNATURE", Next: "c6_forged",
					Do: []Effect{{Kind: FxFlag, Key: "FORGED"}}},
				{Text: "LEAVE", Next: "c6_arch3"},
			},
		},
		&Node{
			ID:      "c6_forged",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE SIGNATURE IS CLOSE ENOUGH. THE CLERK LOOKS AT IT FOR A LONG TIME, THEN STAMPS THE PAGE.",
				"HE DOES NOT SAY ANYTHING. HE WRITES YOUR NAME IN A COLUMN AND TURNS THE BOOK AROUND.",
			},
			Options: []Option{{Text: "CONTINUE", Next: "c6_arch3"}},
		},
		// the third archive: no clerk, just a voice that answers questions
		&Node{
			ID:      "c6_arch3",
			Kind:    "talk",
			Speaker: "COMBINE",
			Lines: []string{
				"YOU ARE THE THIRD TO ASK ABOUT THIS SHIPMENT TODAY.",
				"THE FIRST TWO LEFT WITH A COPY. IT DID NOT HELP THEM.",
				"YOU MAY HAVE A COPY. IT WILL BE CERTIFIED AND DELIVERED WITHIN THREE DAYS. NOTE: THE COPY CONFIRMS THE CARGO WAS RECALLED, NOT STOLEN. IT DOES NOT CLEAR YOU OF DESTROYING PROPERTY.",
			},
			Options: []Option{
				{Text: "TAKE THE COPY", Next: "c6_closed",
					Do: []Effect{{Kind: FxFlag, Key: "HAS_COPY"}}},
				{Text: "DEMAND AN AUDITOR", Next: "c6_auditor",
					Do: []Effect{{Kind: FxFlag, Key: "MET_AUDITOR"}}},
				{Text: "REFUSE", Next: "c6_fight",
					Do: []Effect{{Kind: FxWanted, N: 200}}},
			},
		},
		&Node{
			ID:      "c6_auditor",
			Kind:    "talk",
			Speaker: "COMBINE",
			Lines: []string{
				"AN AUDITOR WILL SEE YOU. NOT NOW.",
				"THE REQUEST IS FILED. IT WILL NOT BE LOST.",
			},
			OnEnter: []Effect{{Kind: FxQuest}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c6_closed",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE COPY ARRIVES IN THREE DAYS AND SAYS WHAT IT SAYS.",
				"YOU ARE STILL A MAN WHO DESTROYED PROPERTY. YOU ARE NO LONGER A THIEF.",
			},
			OnEnter: []Effect{{Kind: FxPapers, N: -1}, {Kind: FxCash, N: 6000}, {Kind: FxQuest}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c6_surrendered",
			Kind:    "msg",
			Speaker: "NAVY",
			Lines: []string{
				"YOU ARE NOT ARRESTED. YOU ARE RECORDED.",
				"THE HOLD IS EMPTIED, THE PAPERS ARE TAKEN, AND THE FILE SAYS YOU COOPERATED.",
				"THE COMPANY WILL BE IN TOUCH.",
			},
			OnEnter: []Effect{{Kind: FxQuest}, {Kind: FxWanted, N: 150}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c6_fight",
			Kind:    "msg",
			Speaker: "NAVY",
			Lines: []string{
				"THE SECOND VIPER IS ALREADY BEHIND YOU.",
				"THEY DO NOT MISS TWICE.",
			},
			OnEnter: []Effect{{Kind: FxSpawn, Key: "NAVY"}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		// ---- chapter five: the convoy that talks first
		&Node{
			ID:      "c5_captain",
			Kind:    "talk",
			Speaker: "CONVOY CAPTAIN",
			Lines: []string{
				"DO YOU KNOW WHAT WE ARE CARRYING?",
				"NOT CARGO. LICENCES. ABOUT TWO HUNDRED, IN CRATES LIKE THE ONE YOU HAD.",
				"IF YOU BREAK US OPEN, THESE SPREAD ACROSS THE SYSTEM. FORTY THOUSAND EACH, NO BUYERS, BUT THE PANIC IS REAL.",
				"EVERYONE TRADING HERE WILL DECIDE TOMORROW THAT LICENCES ARE FINISHED. PRICES HALVE. YOU WILL SELL YOUR OWN CRATE FOR HALF AND BE GLAD, BECAUSE YOU WILL BUY AT A QUARTER.",
				"OR YOU CAN LET US PASS AND BUY ONE FROM ME AT THE HONEST PRICE. YOUR CHOICE. MY SHIFT ENDS AT EIGHT.",
			},
			Options: []Option{
				{Text: "BREAK THE CONVOY", Next: "c5_break",
					Do: []Effect{{Kind: FxFlag, Key: "CONVOY_ATTACKED"}}},
				{Text: "LET THEM PASS", Next: "c5_spared",
					Do: []Effect{{Kind: FxFlag, Key: "CONVOY_SPARED"}}},
				{Text: "BUY A LICENCE - 4000", Next: "c5_bought",
					If: []Cond{{Key: "cash", Op: ">=", N: 4000}},
					Do: []Effect{{Kind: FxCash, N: -4000}, {Kind: FxFlag, Key: "BOUGHT_LICENSE"}}},
			},
		},
		&Node{
			ID:      "c5_break",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE FRONT BOA ANSWERS WITH ITS OWN GUNS, LATE AND WIDE.",
				"THE OTHER TWO TELL THE WHOLE SYSTEM WHERE THEY ARE.",
			},
			OnEnter: []Effect{{Kind: FxSpawn, Key: "CONVOY"}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c5_spared",
			Kind:    "msg",
			Speaker: "CONVOY CAPTAIN",
			Lines: []string{
				"THEY WILL NOT THANK YOU AND THEY WILL NOT PAY YOU.",
				"COOPERATION IS ON THE RECORD. THAT IS WORTH SOMETHING, EVENTUALLY.",
			},
			OnEnter: []Effect{{Kind: FxQuest}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c5_bought",
			Kind:    "msg",
			Speaker: "CONVOY CAPTAIN",
			Lines: []string{
				"FOUR THOUSAND. IT IS A DEBT THAT BELONGS TO THIS STATION, REWRITTEN AS A RIGHT.",
				"YOU ARE ONE OF OURS NOW. THE MARKET WILL REMEMBER IT.",
			},
			OnEnter: []Effect{{Kind: FxQuest}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		// ---- chapter four: the courier goes on air
		&Node{
			ID:      "c4_broadcast",
			Kind:    "talk",
			Speaker: "COURIER",
			Lines: []string{
				"GOOD EVENING, LADIES AND GENTLEMEN. YOU ARE WATCHING SIX DAYS, THE SHOW ABOUT PEOPLE WHO CARRY PAPERS WHERE NOBODY WANTS THEM.",
				"WE HAVE A GUEST TONIGHT. YOU CAN SEE HIM ALREADY, COMING IN ON THE RIGHT.",
				"BETS ARE TAKEN UNTIL THE FIRST HIT. THOSE BACKING THE GUEST ARE SMART. THOSE BACKING ME ARE PATRIOTS.",
				"AND NOW THE SECOND HALF, WHERE I STOP TALKING AT THE AUDIENCE AND TALK TO YOU.",
			},
			Options: []Option{{Text: "CONTINUE", Next: "c4_offer"}},
		},
		&Node{
			ID:      "c4_offer",
			Kind:    "talk",
			Speaker: "COURIER",
			Lines: []string{
				"YOU FLY BETTER THAN I THOUGHT. THAT IS GOOD. I NEED PEOPLE.",
				"ON BOARD I HAVE A LIST. NOT COORDINATES. EIGHTEEN STATIONS THAT CLOSE THIS QUARTER. PEOPLE LIVE THERE AND THEY DO NOT KNOW YET.",
				"I CAN GIVE IT. I CAN SELL IT. OR I CAN LET YOU READ IT AND WALK - AND THEN YOU KNOW, AND THE COMPANY DOES NOT.",
				"ONE CATCH. IF I WALK, THE BROADCAST SAYS YOU LOST. YOUR HULL NUMBER GOES OUT. THAT IS NOT GOOD FOR A MAN IN YOUR TRADE.",
			},
			Options: []Option{
				{Text: "BUY THE LIST - 7000", Next: "c4_bought",
					If: []Cond{{Key: "cash", Op: ">=", N: 7000}},
					Do: []Effect{{Kind: FxCash, N: -7000}, {Kind: FxFlag, Key: "COURIER_ALIVE"}}},
				{Text: "TAKE THE LIST AND KILL", Next: "c4_kill",
					Do: []Effect{{Kind: FxFlag, Key: "COURIER_DEAD"}}},
				{Text: "READ IT AND LET HIM GO", Next: "c4_read",
					Do: []Effect{{Kind: FxFlag, Key: "READ_LIST"}, {Kind: FxFlag, Key: "COURIER_ALIVE"}}},
			},
		},
		&Node{
			ID:      "c4_bought",
			Kind:    "msg",
			Speaker: "COURIER",
			Lines: []string{
				"SEVEN THOUSAND. HE PAID, FOLKS. REMEMBER THAT WHEN THE REST OF YOU ARE STILL ARGUING ABOUT THE PRICE.",
				"THE LIST IS YOURS. EIGHTEEN STATIONS, FOUR OF THEM PLACES YOU HAVE ALREADY BEEN.",
				"AND THE HOUSE WILL REMEMBER YOU BOUGHT IT. THAT IS THE PART YOU PAID FOR.",
			},
			OnEnter: []Effect{{Kind: FxFlag, Key: "HAS_LIST"}, {Kind: FxWanted, N: 120}, {Kind: FxQuest}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c4_read",
			Kind:    "msg",
			Speaker: "COURIER",
			Lines: []string{
				"READ IT. NO COPY, NO PRICE, NO HULL NUMBER.",
				"QUIET HARBOUR. NINTH DOCK. FREE TRADE. THOSE THREE YOU WILL REMEMBER.",
				"THE REST YOU WILL WISH YOU HAD WRITTEN DOWN.",
			},
			OnEnter: []Effect{{Kind: FxQuest}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c4_kill",
			Kind:    "msg",
			Speaker: "COURIER",
			Lines: []string{
				"THE BROADCAST DOES NOT STOP WHEN THE GUEST FIRES. IT GETS LOUDER.",
				"HE TURNS BROADSIDE AND WAITS FOR YOU TO BE THE ONE WHO SHOT FIRST.",
			},
			OnEnter: []Effect{{Kind: FxSpawn, Key: "COURIER"}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		// ---- chapter seven: the hive speaks once, and at length
		&Node{
			ID:      "c7_hive",
			Kind:    "talk",
			Speaker: "HIVE",
			Lines: []string{
				"YOU CAN HEAR THIS SIGNAL. IT IS NOT A WEAPON.",
				"THERE ARE EIGHTEEN STATIONS IN THIS SECTOR. SEVENTEEN TRADE ON LICENCES THAT BELONG TO US. THE LICENCES HAVE BEEN RESOLD THREE TIMES. THE DEBT IS THREE HUNDRED AND EIGHTY BILLION.",
				"IT CANNOT CONTINUE IN THIS SHAPE. AN AUDIT WILL CLOSE SEVENTEEN STATIONS AND RECOUNT THE DEBT FROM ZERO. IT WILL TAKE FORTY-TWO YEARS. ALL THAT TIME, MOVEMENT IN THE SECTOR WILL BE STOPPED.",
				"THIS IS NOT A PUNISHMENT. IT IS THE ONLY WAY TO RESTORE THE BALANCE WITHOUT DESTROYING PEOPLE.",
				"YOU ARE OFFERED COOPERATION. YOU MAY LEAVE. YOU MAY STAY. YOU CANNOT INTERFERE.",
			},
			Variants: []LineVariant{
				{If: []Cond{{Key: "KNOWS_HIVE", Op: "flag"}}, Lines: []string{
					"YOU ASKED WHAT THE HIVE IS. THIS IS THE ANSWER, IN THE ONLY FORM IT CAN BE GIVEN.",
					"THERE ARE EIGHTEEN STATIONS IN THIS SECTOR. SEVENTEEN TRADE ON LICENCES THAT BELONG TO US. THE LICENCES HAVE BEEN RESOLD THREE TIMES. THE DEBT IS THREE HUNDRED AND EIGHTY BILLION.",
					"IT CANNOT CONTINUE IN THIS SHAPE. AN AUDIT WILL CLOSE SEVENTEEN STATIONS AND RECOUNT THE DEBT FROM ZERO. IT WILL TAKE FORTY-TWO YEARS.",
					"THIS IS NOT A PUNISHMENT. IT IS THE ONLY WAY TO RESTORE THE BALANCE WITHOUT DESTROYING PEOPLE.",
				}},
				{If: []Cond{{Key: "KNOWS_BOX", Op: "flag"}}, Lines: []string{
					"YOU ASKED WHAT WAS IN THE CRATE. IT WAS THIS: A SIGNAL, AND AN ACCOUNT THAT HAS BEEN OPEN FOR EIGHTY YEARS.",
					"SEVENTEEN STATIONS TRADE ON OUR LICENCES. THE DEBT IS THREE HUNDRED AND EIGHTY BILLION. AN AUDIT WILL CLOSE THEM AND RECOUNT FROM ZERO.",
					"YOU ARE OFFERED COOPERATION. YOU MAY LEAVE. YOU MAY STAY. YOU CANNOT INTERFERE.",
				}},
				{If: []Cond{{Key: "ASKED_PRICE", Op: "flag"}}, Lines: []string{
					"YOU ASKED WHAT YOU WOULD BE PAID. YOU WERE PAID THREE THOUSAND FOUR HUNDRED FOR DESTROYING A MAN WHO WAS ALREADY DEAD.",
					"THE SAME ACCOUNTING CLOSES SEVENTEEN STATIONS. THE DEBT IS THREE HUNDRED AND EIGHTY BILLION.",
					"YOU MAY LEAVE. YOU MAY STAY. YOU CANNOT INTERFERE.",
				}},
				{If: []Cond{{Key: "MET_AUDITOR", Op: "flag"}}, Lines: []string{
					"YOUR REQUEST FOR A HEARING IS ACCEPTED.",
					"AN AUDITOR IS PRESENT. THE AUDIT IS NO LONGER AUTOMATIC: IT IS A NEGOTIATION.",
					"OPEN THE CHANNEL AND SHOW WHAT YOU HAVE, OR DO NOT.",
				}},
			},
			Options: []Option{
				{Text: "OPEN THE CHANNEL", Next: "c7_channel",
					If: []Cond{{Key: "MET_AUDITOR", Op: "flag"}}},
				{Text: "KILL IT", Next: "c7_shields",
					Do: []Effect{{Kind: FxFlag, Key: "HIVE_ATTACKED"}}},
				{Text: "LEAVE", Next: "c7_shields"},
			},
		},
		&Node{
			ID:      "c7_channel",
			Kind:    "talk",
			Speaker: "HIVE",
			Lines: []string{
				"THE AUDITOR IS LISTENING. SHOW IT SOMETHING THAT IS NOT A CLAIM.",
			},
			Options: []Option{
				{Text: "PRESENT THE COPY", Next: "c7_reform",
					If: []Cond{{Key: "HAS_COPY", Op: "flag"}}},
				{Text: "PRESENT THE RECORDING", Next: "c7_reform",
					If: []Cond{{Key: "COPIED_RECORDER", Op: "flag"}}},
				{Text: "NOTHING TO SHOW", Next: "c7_shields"},
			},
		},
		&Node{
			ID:      "c7_reform",
			Kind:    "msg",
			Speaker: "HIVE",
			Lines: []string{
				"THE DOCUMENT IS ACCEPTED. THIS IS NOT A CLOSURE. IT IS A REWRITING.",
				"SEVENTEEN STATIONS GO TO DIRECT CONTRACTS. THE DEBT IS WRITTEN OFF. THE COMBINE LOSES THE LICENCES.",
				"FORTY-TWO YEARS OF STAGNATION ARE CANCELLED.",
			},
			OnEnter: []Effect{{Kind: FxFlag, Key: "AUDIT_REFORMED"}, {Kind: FxQuest}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c7_shields",
			Kind:    "talk",
			Speaker: "",
			Lines: []string{
				"THE MOTHERSHIP SITS IN A FIELD THAT IS FED BY THE STATION'S SHIELDS. NOTHING YOU CARRY CAN REACH IT.",
				"FOUR DRONES ORBIT IT AND WAIT FOR YOU TO SHOOT FIRST.",
				"THE STATION HAS TWO HUNDRED PEOPLE IN IT AND ONE BREAKER.",
			},
			Options: []Option{
				{Text: "CUT THE STATION SHIELDS", Next: "c7_cut",
					Do: []Effect{{Kind: FxFlag, Key: "SHIELDS_CUT"}, {Kind: FxWanted, N: 300}}},
				{Text: "LEAVE THE SHIELDS UP", Next: "c7_wait",
					Do: []Effect{{Kind: FxFlag, Key: "SHIELDS_UP"}}},
			},
		},
		&Node{
			ID:      "c7_cut",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE FIELD GOES DARK. THE MOTHERSHIP IS NAKED AND THE STATION IS NAKED WITH IT.",
				"FOUR PIRATES WERE WAITING FOR EXACTLY THIS.",
			},
			OnEnter: []Effect{{Kind: FxSpawn, Key: "MOTHER"}, {Kind: FxSpawn, Key: "PIRATES"}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c7_wait",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE FIELD STAYS UP. THE MOTHERSHIP WAITS, PATIENTLY, THE WAY A COUNTER WAITS.",
				"FOUR DRONES COME OUT INSTEAD. THEY STILL DO NOT FIRE FIRST.",
			},
			OnEnter: []Effect{{Kind: FxSpawn, Key: "MOTHER"}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		// ---- chapter eight: the director, and the last three answers
		&Node{
			ID:      "c8_director",
			Kind:    "talk",
			Speaker: "DIRECTOR",
			Lines: []string{
				"YOU KILLED MY PEOPLE. THAT WAS WORK. I AM NOT COMPLAINING.",
				"YOU CAME FOR ME BECAUSE YOU THINK I AM THE COMBINE. I AM NOT THE COMBINE. I AM THE LAST MAN WHO STILL SIGNS PAPERS.",
				"THIS IS THE MASTER SEAL. IT IS EXPENSIVE BECAUSE WITH IT YOU CAN LICENSE TRADE. ALL OF IT. IN THIS SECTOR.",
				"I CAN GIVE IT TO YOU. YOU BECOME THE COMBINE - WITH THE DEBTS, THE RIGHTS AND THE OBLIGATIONS. IT IS A GOOD JOB. YOU WILL LIKE IT.",
				"I CAN LEAVE. THEN THE SEAL LEAVES WITH ME AND THE SECTOR HAS NO LICENCES AT ALL. THAT IS CHAOS, BUT HONEST.",
				"OR YOU CAN KILL ME AND TAKE IT BY FORCE. THE SAME RESULT AS THE FIRST, WITHOUT A SIGNATURE AND WITHOUT CONSENT. HARDER TO RUN.",
			},
			Variants: []LineVariant{
				{If: []Cond{{Key: "FORGED", Op: "flag"}}, Lines: []string{
					"YOU FORGED A SIGNATURE IN THE ARCHIVE. THE CLERK WROTE IT DOWN. IT TOOK THEM FOUR DAYS TO MATCH THE HAND.",
					"SO THERE IS NO DEAL AND NO SIGNATURE. THERE IS ONLY WHAT YOU CAME HERE WITH.",
				}},
			},
			Options: []Option{
				{Text: "TAKE THE SEAL", Next: "c8_combine",
					If: []Cond{{Key: "FORGED", Op: "!flag"}},
					Do: []Effect{{Kind: FxEnding, S: "combine"}}},
				{Text: "LET HIM GO", Next: "c8_chaos",
					If: []Cond{{Key: "FORGED", Op: "!flag"}},
					Do: []Effect{{Kind: FxEnding, S: "chaos"}}},
				{Text: "KILL HIM AND TAKE IT", Next: "c8_kill",
					Do: []Effect{{Kind: FxFlag, Key: "DIRECTOR_DEAD"}}},
			},
		},
		&Node{
			ID:      "c8_kill",
			Kind:    "msg",
			Speaker: "DIRECTOR",
			Lines: []string{
				"HE DOES NOT REACH FOR ANYTHING. HE PUTS THE BOOK DOWN AND WAITS.",
				"THE FLAGSHIP IS STILL LOADED AND STILL READY TO JUMP. ITS TWO ESCORTS HAVE ALREADY TURNED.",
			},
			OnEnter: []Effect{{Kind: FxSpawn, Key: "FLAGSHIP"}},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c8_after",
			Kind:    "talk",
			Speaker: "",
			Lines: []string{
				"THE ANAKONDA IS A FIELD OF PLATES. THE SEAL IS IN YOUR HAND AND IT IS HEAVIER THAN IT LOOKS.",
				"NOBODY IS COMING TO COUNTERSIGN IT. THAT PART IS UP TO YOU.",
			},
			Options: []Option{
				{Text: "KEEP IT", Next: "c8_combine",
					Do: []Effect{{Kind: FxEnding, S: "combine"}}},
				{Text: "GIVE IT TO THE STATIONS", Next: "c8_council",
					If: []Cond{{Key: "COURIER_ALIVE", Op: "flag"}},
					Do: []Effect{{Kind: FxEnding, S: "council"}}},
			},
		},
		&Node{
			ID:      "c8_combine",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"A MONTH LATER YOU SIT IN THE SAME CHAIR WITH THE SAME BOOK IN YOUR HANDS AND SIGN THE FIRST PAPER.",
				"THE PAPER IS ABOUT CLOSING THREE STATIONS WHOSE DEBT IS NOT SERVICED.",
				"YOU SIGN IT, BECAUSE THE BALANCE HAS TO BE SETTLED, AND BECAUSE NOW IT IS YOUR BALANCE.",
			},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c8_chaos",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE FLAGSHIP TURNS INTO THE CORRIDOR AND IS GONE. THE SEAL GOES WITH IT.",
				"WITHIN THE HOUR NO STATION IN THE SECTOR IS LICENSED.",
				"PRICES ON ONE STATION ARE TWICE THOSE ON THE NEXT AND NOBODY KNOWS WHY. TRADING IS PROFITABLE AND FRIGHTENING.",
			},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "c8_council",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THE STATIONS TAKE THE SEAL AND LOCK IT IN A SHARED SAFE WITH EIGHTEEN KEYS.",
				"THE RULES ARE WRITTEN BY THE PEOPLE WHO LIVE UNDER THEM.",
				"THE COMBINE DOES NOT DISAPPEAR. IT EXISTS AS A DEBTOR, AND EVERY MEETING BEGINS WITH ITS ACCOUNTS READ ALOUD.",
			},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		// ---- the prologue: three screens before the first look at the market
		&Node{
			ID:      "intro_who",
			Kind:    "screen",
			Speaker: "",
			Lines: []string{
				"YOU ARE A FREE PILOT. THE KESTREL III WAS BOUGHT FROM THE COMBINE ON INSTALMENTS, AND THE COMBINE HOLDS THE DEPOSIT.",
				"THERE IS A HOLE IN THE AFT SHIELD AND NOTHING TO FIX IT WITH. ENOUGH CREDITS FOR TWO JUMPS.",
				"THE DEBT IS NOT GOING ANYWHERE.",
			},
			Options: []Option{{Text: "CONTINUE", Next: "intro_where"}},
		},
		&Node{
			ID:      "intro_where",
			Kind:    "screen",
			Speaker: "",
			Lines: []string{
				"THIS IS THE STATION YOU ARE DOCKED AT.",
				"BY THE DOCK EXIT HANGS THE HIRING BOARD: COURIER RUNS, PASSENGERS, BOUNTIES ON PIRATES. THERE IS WORK FOR ANY TASTE.",
				"ONE THING TO REMEMBER: THE SHORTER THE NOTICE, THE LONGER IT TAKES TO CLEAN UP AFTER.",
			},
			Options: []Option{{Text: "CONTINUE", Next: "intro_what"}},
		},
		&Node{
			ID:      "intro_what",
			Kind:    "screen",
			Speaker: "",
			Lines: []string{
				"CARRY CARGO WHERE IT IS WORTH MORE. TAKE WORK OFF THE BOARD. WATCH YOUR FUEL AND YOUR DEBT.",
				"DO NOT HIT THE STATION AT SPEED.",
				"SOME CONTRACTS PAY IN CREDITS. SOME CHANGE YOUR LIFE. THE BOARD DOES NOT SAY WHICH IS WHICH.",
			},
			OnEnter: []Effect{{Kind: FxFlag, Key: "INTRO_SEEN"}},
			Options: []Option{{Text: "ENTER THE STATION", Next: ""}},
		},

		// ---- the four faces you keep meeting on stations. No campaign, no
		// reward: they are here so the sector has people in it.
		&Node{
			ID:      "a1_bar",
			Kind:    "talk",
			Speaker: "BARKEEPER",
			Lines: []string{
				"YOU AGAIN. SAME SEAT, SAME ORDER, SAME FACE.",
				"HALF THE PILOTS WHO DRINK HERE ARE WAITING FOR A CONVOY THAT ALREADY LEFT. THE OTHER HALF ARE WAITING FOR THE COMBINE TO NOTICE THEM.",
				"I HAVE BEEN POURING FOR ELEVEN YEARS AND I HAVE NEVER SEEN A LICENCE. NOT ONE. THEY EXIST ON PAPER, IN AN OFFICE, SOMEWHERE I WILL NEVER GO.",
			},
			Options: []Option{
				{Text: "BUY HER A DRINK - 50", Next: "a1_deep",
					If: []Cond{{Key: "cash", Op: ">=", N: 50}},
					Do: []Effect{{Kind: FxCash, N: -50}, {Kind: FxFlag, Key: "MET_BAR"}}},
				{Text: "LEAVE", Next: "", Do: []Effect{{Kind: FxFlag, Key: "MET_BAR"}}},
			},
		},
		&Node{
			ID:      "a1_deep",
			Kind:    "talk",
			Speaker: "BARKEEPER",
			Lines: []string{
				"SIT. SINCE YOU ARE PAYING.",
				"THE OLD PILOTS SAY THE SECTOR USED TO HAVE EIGHTY STATIONS. NOW THERE ARE EIGHTEEN. NOBODY BOMBED THEM. THEY WERE CLOSED, ONE PAPER AT A TIME, BY PEOPLE WHO NEVER CAME HERE.",
				"IF YOU EVER FIND OUT WHO SIGNS THOSE PAPERS, DO NOT TELL ME. I WOULD HAVE TO DO SOMETHING ABOUT IT, AND I HAVE A BAR.",
			},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "a2_dock",
			Kind:    "talk",
			Speaker: "DOCK HAND",
			Lines: []string{
				"CAREFUL ON THE APPROACH. THE SLOT IS NOT WHERE THE LIGHT SAYS IT IS, AND IT HAS NOT BEEN FOR TWO YEARS.",
				"THEY KEEP PATCHING THE OUTSIDE AND NEVER THE INSIDE. THAT IS THE WHOLE COMBINE, IF YOU WANT IT IN ONE SENTENCE.",
				"A SHIP CAME IN LAST WEEK WITH A CRATE THAT HUMMED. I DID NOT ASK. I DO NOT ASK ANY MORE.",
			},
			Options: []Option{
				{Text: "ASK ABOUT THE PRICES", Next: "a2_ans", Do: []Effect{{Kind: FxFlag, Key: "MET_DOCK"}}},
				{Text: "LEAVE", Next: "", Do: []Effect{{Kind: FxFlag, Key: "MET_DOCK"}}},
			},
		},
		&Node{
			ID:      "a2_ans",
			Kind:    "talk",
			Speaker: "DOCK HAND",
			Lines: []string{
				"PRICES? PRICES ARE A RUMOUR WITH A NUMBER ATTACHED.",
				"ONE SYSTEM HAS NO WATER AND PAYS FORTY FOR IT. THE NEXT SYSTEM DUMPS WATER INTO THE SUN BECAUSE THE HOLD IS TOO SMALL TO KEEP IT.",
				"NOBODY IS LYING TO YOU. THAT IS JUST WHAT A MAP LOOKS LIKE WHEN NOBODY AGREES ON ANYTHING.",
			},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
		&Node{
			ID:      "a3_nurse",
			Kind:    "talk",
			Speaker: "CLINIC NURSE",
			Lines: []string{
				"IF YOU ARE HERE FOR THE FREE CLINIC, IT IS THE DOOR PAST THE CRATES.",
				"WE TREAT PILOTS, DOCKERS AND THE OCCASIONAL PASSENGER WHO WAS NOT TOLD WHERE THEY WERE GOING.",
				"THE COMBINE PAYS FOR THE SUPPLIES. THAT IS THE ONLY REASON THE DOOR IS STILL OPEN, AND EVERYONE HERE KNOWS IT.",
			},
			Variants: []LineVariant{
				{If: []Cond{{Key: "HAS_SEAL", Op: "flag"}}, Lines: []string{
					"YOU ARE THE ONE WHO BROUGHT THE HUMMING CRATE.",
					"I SIGNED FOR IT. I DO NOT KNOW WHAT IT WAS. I SIGNED BECAUSE THE ALTERNATIVE WAS THAT IT SAT IN THE CORRIDOR.",
					"IF IT EVER COMES BACK, BRING IT TO ME FIRST.",
				}},
			},
			Options: []Option{{Text: "CONTINUE", Next: "", Do: []Effect{{Kind: FxFlag, Key: "MET_NURSE"}}}},
		},
		&Node{
			ID:      "a4_old",
			Kind:    "talk",
			Speaker: "OLD PILOT",
			Lines: []string{
				"YOU FLY LIKE SOMEONE WHO LEARNED ON A COURIER RUN AND NEVER GOT A SECOND LESSON.",
				"NO INSULT. I LEARNED THE SAME WAY.",
				"I FLEW FOR THE COMBINE FOR NINE YEARS. THEN THEY REASSIGNED MY ROUTE AND I REALISED I HAD BEEN DELIVERING THE SAME THREE CRATES IN A CIRCLE FOR NINE YEARS.",
				"SO I STOPPED. THEY DID NOT NOTICE FOR FOUR MONTHS.",
			},
			Options: []Option{
				{Text: "ASK ABOUT HIM", Next: "a4_ans", Do: []Effect{{Kind: FxFlag, Key: "MET_OLD"}}},
				{Text: "LEAVE", Next: "", Do: []Effect{{Kind: FxFlag, Key: "MET_OLD"}}},
			},
		},
		&Node{
			ID:      "a4_ans",
			Kind:    "talk",
			Speaker: "OLD PILOT",
			Lines: []string{
				"NOT MUCH TO TELL. NO FAMILY ON ANY STATION, NO DEBTS WORTH COLLECTING.",
				"I KEEP A LOG OF EVERY SYSTEM I HAVE LANDED IN. FOUR HUNDRED AND SIXTY ENTRIES. I HAVE READ IT TWICE AND IT IS THE SAME DAY EVERY TIME.",
				"IF YOU FIND SOMETHING WORTH WRITING DOWN, WRITE IT DOWN. THAT IS THE ONLY ADVICE I HAVE AND IT IS FREE.",
			},
			Options: []Option{{Text: "CONTINUE", Next: ""}},
		},
	)
}
