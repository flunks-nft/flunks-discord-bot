package battle

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/flunks-nft/discord-bot/pkg/utils"
)

type WeaponOutcomeOption struct {
	displayMsg    string
	challengerWon bool
}
type WeaponOption struct {
	weapon            string
	displayMsgOptions []string
	outcomeOptions    []WeaponOutcomeOption
}

type Battle struct {
	Weapon  string
	Action  string
	Outcome WeaponOutcomeOption
	Winner  int // 0 for challenger, 1 for defender
}

func (b Battle) Log() BattleLog {
	return BattleLog{
		Action:        b.Action,
		ActionOutcome: b.Outcome.displayMsg,
		BattleOutcome: b.Weapon,
	}
}

type BattleLog struct {
	Action        string `json:"action"`
	ActionOutcome string `json:"actionOutcome"`
	BattleOutcome string `json:"battleOutcome"`
}

// Implement the sql.Scanner interface
func (bt *BattleLog) Scan(value interface{}) error {
	if value == nil {
		*bt = BattleLog{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Failed to scan BattleThread: %#v", value)
	}
	return json.Unmarshal(bytes, bt)
}

// Implement the driver.Valuer interface
func (bt BattleLog) Value() (driver.Value, error) {
	return json.Marshal(bt)
}

func (wo WeaponOption) drawBattle() Battle {
	if len(wo.displayMsgOptions) == 0 {
		panic("drawBattle(): No displayMsgOptions")
	}

	if len(wo.outcomeOptions) == 0 {
		panic("drawBattle(): No outcomeOptions")
	}

	displayMsg := utils.RandomItem(wo.displayMsgOptions).(string)
	weaponOutcomeOption := utils.RandomItem(wo.outcomeOptions).(WeaponOutcomeOption)

	// Get winner
	var winner int
	if weaponOutcomeOption.challengerWon == true {
		winner = 0
	} else {
		winner = 1
	}

	return Battle{
		Weapon:  wo.weapon,
		Action:  displayMsg,
		Outcome: weaponOutcomeOption,
		Winner:  winner,
	}
}

func DrawBattleByClique(clique string) Battle {
	// Pick a random weapon
	weaponOption := utils.RandomItem(BattleWeaponMap[clique]).(WeaponOption)

	// draw battle details from the picked weapon
	return weaponOption.drawBattle()
}

var BattleWeaponMap = map[string][]WeaponOption{
	"Jock": {
		{
			weapon: "Protein Shake",
			displayMsgOptions: []string{
				"Look at those gains! {Flunk_#1} is getting prepped with a **Protein Shake**.",
				"You’re not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I’m not sure how effective a **Protein Shake** is in this situation. Maybe it’s psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don’t work on their own, Flunk_#1. Hit the gym!",
				"I think Flunk_#2 could use that **Protein Shake** more than you, Flunk_#1. But seriously, what are you planning to do with that right now?",
			},
			outcomeOptions: []WeaponOutcomeOption{
				{
					displayMsg:    "BOOM! PROTEIN SHAKE TO THE DOME, Flunk_#1 just lobbed that shake straight into their face.",
					challengerWon: true,
				},
				{
					displayMsg:    "That should not be allowed in school, Flunk number one has turned green after taking the protein shake and SMASHES Flunk_#2 to smithereens.",
					challengerWon: true,
				},
				{
					displayMsg:    "Now that's lucky! Flunk_#1 spilt his protein shake after running away scared but Flunk_#2 slipped on the lumpy mess and knocked themselves out COLD! STONE COLD!",
					challengerWon: true,
				},
				{
					displayMsg:    "Flunk_#1 downed the lumpy-ass protein shake in one, then flexed his muscles so much they passed out.",
					challengerWon: false,
				},
				{
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he’s starting to break down and cry for some weird reason.",
					challengerWon: false,
				},
				{
					displayMsg:    "OH HELL NO! The janitor isn't paid enough for this. Flunk_#1 has thrown up all over themselves before a punch is even thrown, that protein shake smells horrid! EWWW",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Protein Shake 2",
			displayMsgOptions: []string{
				"Look at those gains! {Flunk_#1} is getting prepped with a **Protein Shake**.",
				"You’re not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I’m not sure how effective a **Protein Shake** is in this situation. Maybe it’s psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don’t work on their own, Flunk_#1. Hit the gym!",
				"I think Flunk_#2 could use that **Protein Shake** more than you, Flunk_#1. But seriously, what are you planning to do with that right now?",
			},
			outcomeOptions: []WeaponOutcomeOption{
				{
					displayMsg:    "BOOM! PROTEIN SHAKE TO THE DOME, Flunk_#1 just lobbed that shake straight into their face.",
					challengerWon: true,
				},
				{
					displayMsg:    "That should not be allowed in school, Flunk number one has turned green after taking the protein shake and SMASHES Flunk_#2 to smithereens.",
					challengerWon: true,
				},
				{
					displayMsg:    "Now that's lucky! Flunk_#1 spilt his protein shake after running away scared but Flunk_#2 slipped on the lumpy mess and knocked themselves out COLD! STONE COLD!",
					challengerWon: true,
				},
				{
					displayMsg:    "Flunk_#1 downed the lumpy-ass protein shake in one, then flexed his muscles so much they passed out.",
					challengerWon: false,
				},
				{
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he’s starting to break down and cry for some weird reason.",
					challengerWon: false,
				},
				{
					displayMsg:    "OH HELL NO! The janitor isn't paid enough for this. Flunk_#1 has thrown up all over themselves before a punch is even thrown, that protein shake smells horrid! EWWW",
					challengerWon: false,
				},
			},
		},
	},
	"Prep": {
		WeaponOption{
			weapon: "Protein Shake",
			displayMsgOptions: []string{
				"Look at those gains! {Flunk_#1} is getting prepped with a **Protein Shake**.",
				"You’re not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I’m not sure how effective a **Protein Shake** is in this situation. Maybe it’s psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don’t work on their own, Flunk_#1. Hit the gym!",
				"I think Flunk_#2 could use that **Protein Shake** more than you, Flunk_#1. But seriously, what are you planning to do with that right now?",
			},
			outcomeOptions: []WeaponOutcomeOption{
				{
					displayMsg:    "BOOM! PROTEIN SHAKE TO THE DOME, Flunk_#1 just lobbed that shake straight into their face.",
					challengerWon: true,
				},
				{
					displayMsg:    "That should not be allowed in school, Flunk number one has turned green after taking the protein shake and SMASHES Flunk_#2 to smithereens.",
					challengerWon: true,
				},
				{
					displayMsg:    "Now that's lucky! Flunk_#1 spilt his protein shake after running away scared but Flunk_#2 slipped on the lumpy mess and knocked themselves out COLD! STONE COLD!",
					challengerWon: true,
				},
				{
					displayMsg:    "Flunk_#1 downed the lumpy-ass protein shake in one, then flexed his muscles so much they passed out.",
					challengerWon: false,
				},
				{
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he’s starting to break down and cry for some weird reason.",
					challengerWon: false,
				},
				{
					displayMsg:    "OH HELL NO! The janitor isn't paid enough for this. Flunk_#1 has thrown up all over themselves before a punch is even thrown, that protein shake smells horrid! EWWW",
					challengerWon: false,
				},
			},
		},
	},
	"Freak": {
		WeaponOption{
			weapon: "Protein Shake",
			displayMsgOptions: []string{
				"Look at those gains! {Flunk_#1} is getting prepped with a **Protein Shake**.",
				"You’re not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I’m not sure how effective a **Protein Shake** is in this situation. Maybe it’s psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don’t work on their own, Flunk_#1. Hit the gym!",
				"I think Flunk_#2 could use that **Protein Shake** more than you, Flunk_#1. But seriously, what are you planning to do with that right now?",
			},
			outcomeOptions: []WeaponOutcomeOption{
				{
					displayMsg:    "BOOM! PROTEIN SHAKE TO THE DOME, Flunk_#1 just lobbed that shake straight into their face.",
					challengerWon: true,
				},
				{
					displayMsg:    "That should not be allowed in school, Flunk number one has turned green after taking the protein shake and SMASHES Flunk_#2 to smithereens.",
					challengerWon: true,
				},
				{
					displayMsg:    "Now that's lucky! Flunk_#1 spilt his protein shake after running away scared but Flunk_#2 slipped on the lumpy mess and knocked themselves out COLD! STONE COLD!",
					challengerWon: true,
				},
				{
					displayMsg:    "Flunk_#1 downed the lumpy-ass protein shake in one, then flexed his muscles so much they passed out.",
					challengerWon: false,
				},
				{
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he’s starting to break down and cry for some weird reason.",
					challengerWon: false,
				},
				{
					displayMsg:    "OH HELL NO! The janitor isn't paid enough for this. Flunk_#1 has thrown up all over themselves before a punch is even thrown, that protein shake smells horrid! EWWW",
					challengerWon: false,
				},
			},
		},
	},
	"Geek": {
		WeaponOption{
			weapon: "Protein Shake",
			displayMsgOptions: []string{
				"Look at those gains! {Flunk_#1} is getting prepped with a **Protein Shake**.",
				"You’re not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I’m not sure how effective a **Protein Shake** is in this situation. Maybe it’s psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don’t work on their own, Flunk_#1. Hit the gym!",
				"I think Flunk_#2 could use that **Protein Shake** more than you, Flunk_#1. But seriously, what are you planning to do with that right now?",
			},
			outcomeOptions: []WeaponOutcomeOption{
				{
					displayMsg:    "BOOM! PROTEIN SHAKE TO THE DOME, Flunk_#1 just lobbed that shake straight into their face.",
					challengerWon: true,
				},
				{
					displayMsg:    "That should not be allowed in school, Flunk number one has turned green after taking the protein shake and SMASHES Flunk_#2 to smithereens.",
					challengerWon: true,
				},
				{
					displayMsg:    "Now that's lucky! Flunk_#1 spilt his protein shake after running away scared but Flunk_#2 slipped on the lumpy mess and knocked themselves out COLD! STONE COLD!",
					challengerWon: true,
				},
				{
					displayMsg:    "Flunk_#1 downed the lumpy-ass protein shake in one, then flexed his muscles so much they passed out.",
					challengerWon: false,
				},
				{
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he’s starting to break down and cry for some weird reason.",
					challengerWon: false,
				},
				{
					displayMsg:    "OH HELL NO! The janitor isn't paid enough for this. Flunk_#1 has thrown up all over themselves before a punch is even thrown, that protein shake smells horrid! EWWW",
					challengerWon: false,
				},
			},
		},
	},
}

// var BattleResultMap = map[string][]WeaponOption{
// 	"Prep": {
// 		WeaponOption{
// 			weapon:     "debate",
// 			displayMsg: "Through a series of elegantly polite debates, Flunk#1 diplomatically defeated Flunk#2 in the Verbal Joust!",
// 		},
// 		WeaponOption{
// 			weapon:     "yacht",
// 			displayMsg: "Flunk#1, sailing with his daddy's new yacht, triumphs over Flunk#2 in a surprise yacht race.",
// 		},
// 		WeaponOption{
// 			weapon:     "Bugatti Chiron",
// 			displayMsg: "Flunk#1 left Flunk#2 in the dust as he raced to victory in his Bugatti Chiron that he got for his 16th birthday.",
// 		},
// 		WeaponOption{
// 			weapon:     "rumble",
// 			displayMsg: "Flunk#1 aced the runway rumble, leaving Flunk#2 in a fashion faux-pas frenzy!",
// 		},
// 		WeaponOption{
// 			weapon:     "croquet",
// 			displayMsg: "Flunk#1 outshined Flunk#2 in the croquet clash, sending his balls everywhere but straight!",
// 		},
// 	},
// }
