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
	)
}
