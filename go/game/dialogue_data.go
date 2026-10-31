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
	)
}
