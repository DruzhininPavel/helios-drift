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
				{Text: "PAY 300", Next: "c6_end",
					If: []Cond{{Key: "cash", Op: ">=", N: 300}},
					Do: []Effect{{Kind: FxCash, N: -300}, {Kind: FxFlag, Key: "GOT_FORM"}}},
				{Text: "NO", Next: "c6_end"},
			},
		},
		&Node{
			ID:      "c6_end",
			Kind:    "msg",
			Speaker: "",
			Lines: []string{
				"THREE FOLDERS IN THE CABINET. YOURS IS THE SECOND FROM THE BOTTOM.",
				"IT IS THINNER THAN THE OTHERS.",
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
	)
}
