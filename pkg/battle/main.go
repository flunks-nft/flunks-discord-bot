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
				"You're not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I'm not sure how effective a **Protein Shake** is in this situation. Maybe it's psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don't work on their own, Flunk_#1. Hit the gym!",
				"I think Flunk_#2 could use that **Protein Shake** more than you, Flunk_#1. But seriously, what are you planning to do with that right now?",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
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
				// Negative Outcomes
				{
					displayMsg:    "Flunk_#1 downed the lumpy-ass protein shake in one, then flexed his muscles so much they passed out.",
					challengerWon: false,
				},
				{
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he's starting to break down and cry for some weird reason.",
					challengerWon: false,
				},
				{
					displayMsg:    "OH HELL NO! The janitor isn't paid enough for this. Flunk_#1 has thrown up all over themselves before a punch is even thrown, that protein shake smells horrid! EWWW",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Golf Club",
			displayMsgOptions: []string{
				"Wow! {Flunk_#1} is pulling out a **3 wood**! Is that thing signed by Tiger?",
				"What on earth is that commotion? {Flunk_#2} scrambles to the **9 iron**. {Flunk_#1} isn't a happy Gilmore…",
				"FOUUUUUURRRR! For flunk sake {Flunk_#1}! You've got the **golf balls** all over the place.",
				"Is that a Birdie… a plane… oh no its {Flunk_#1} with a trusty **pitching wedge** looking to take some swings at {Flunk_#2}",
				"Someones looking for a hole-in-one today! Didn't know this school could afford golf equipment, must've been stolen by {Flunk_#1} from the crazy golf course…",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
				{
					displayMsg:    "That club was meant for golf balls, but today {Flunk_#1} used it for chopping down {Flunk_#2} at the ankles!",
					challengerWon: true,
				},
				{
					displayMsg:    "Driver, 5 iron, or putter. It doesn't matter what club that was, {Flunk_#1} just CRANKED {Flunk_#2}",
					challengerWon: true,
				},
				{
					displayMsg:    "Oh, that wasn't even a golf club. {Flunk_#1} just whacked {Flunk_#2} across the head with a broomstick!",
					challengerWon: true,
				},
				// Negative Outcomes
				{
					displayMsg:    "Wow. I don't think I've seen a worse swing in my life. {Flunk_#1} just swung themselves out of control and into an emergency.",
					challengerWon: false,
				},
				{
					displayMsg:    "I always thought I was the embarrassment of flunks high but jesus christ that swing was straight out of kindergarten, he's run away out of embarrassment",
					challengerWon: false,
				},
				{
					displayMsg:    "{Flunk_#1} keeps swinging frantically as the golf ball stays put on the grass.",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Baseball Bat",
			displayMsgOptions: []string{
				"Flunk_#1 is swinging for the fences today. They just brought out a **Baseball Bat**.",
				"I don't think Flunk_#1 has ever put a ball in play. Maybe they'll have better luck hitting Flunk_#2 with their **Baseball Bat**!",
				"I thought Flunk_#1 was a football player? Where'd they get a **Baseball Bat** from?",
				"Does Flunk_#1 know this isn't baseball practice, and Flunk_#2 is NOT a baseball?",
				"It's DINGER season! Flunk_#1 is ready to go with a **Baseball Bat**!",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
				{
					displayMsg:    "Down goes Anderson! I mean Flunk_#2. That was quite the swing from Flunk_#1.",
					challengerWon: true,
				},
				{
					displayMsg:    "HOME RUN ALL DAY LONG! Well looks like Flunk_#2 is going to need braces after that hit",
					challengerWon: true,
				},
				{
					displayMsg:    "Sweeter than pumpkin pie, this kid has to have a baseball future - but he might have to stop off at the principal's office after that strike! Flunk_#2 is crying out for their mummy",
					challengerWon: true,
				},
				// Negative Outcomes
				{
					displayMsg:    "Swing ‘n' a miss. Yerrr' out, Flunk_#1!",
					challengerWon: false,
				},
				{
					displayMsg:    "Air ball! Flunk_#1 swings into sweet nothingness as Flunk_#2 dips and delivers a gruesome uppercut to end this battle",
					challengerWon: false,
				},
				{
					displayMsg:    "DISGUSTING! As Flunk_#1 swings they dislocated their shoulder and they are OUTTA THERE",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Boxing Gloves",
			displayMsgOptions: []string{
				"Flunk_#2 better watch out! Flunk_#1 just tightened their **Boxing Gloves.**",
				"Call them Mike! Flunk_#1 has got the **Boxing Gloves**.",
				"Okay, who just brings **Boxing Gloves** with them to school? Come on now Flunk_#1…",
				"Flunk_#1 is looking for a K.O. today with those **Boxing Gloves**! Round one is upon us.",
				"Move like a butterfly sting like a bee! Looks like those **Boxing Gloves** are too heavy for Flunk_#1.",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
				{
					displayMsg:    "One two, one two, behind the jab, THATS A BOXING MASTERCLASS",
					challengerWon: true,
				},
				{
					displayMsg:    "DING DING DING, even though it looks like Flunk_#1 couldn't punch to save his life, Flunk_#2 passed out from anxiety, oh dear…",
					challengerWon: true,
				},
				{
					displayMsg:    "BOMMMMB SQUADDDDD! Was that Wilder?! ONE PUNCH and he's out for the count",
					challengerWon: true,
				},
				// Negative Outcomes
				{
					displayMsg:    "You'd think with **Boxing Gloves**, Flunk_#1 would have an edge, but they just swung at air and missed!",
					challengerWon: false,
				},
				{
					displayMsg:    "Flunk_#1 tried to throw a combo, but Flunk_#2 dodged easily, laughing off the pathetic attempt.",
					challengerWon: false,
				},
				{
					displayMsg:    "Oof, those **Boxing Gloves** just slowed Flunk_#1 down. Flunk_#2 took advantage and retaliated quickly.",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Hockey Stick",
			displayMsgOptions: []string{
				"I didn't think it was cold enough for ice, but hey, what do I know? Flunk_#1 has a **Hockey Stick**.",
				"Hey Flunk_#1! Are you strong enough to be using a **Hockey Stick** with 110 Flex?",
				"Woah! How did you fit a **Hockey Stick** in your backpack, Flunk_#1? Is this a magic trick?",
				"Your wrist shot isn't any good, Flunk_#1. Maybe you'll fight better with that **Hockey Stick**.",
				"Puck or no puck that **hockey stick** Flunk_#1's holding is gonna see some messy action tonight",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
				{
					displayMsg:    "Slap Shot BANG! With a strike like that the Maple Leafs should be scouting this Flunk Talent.",
					challengerWon: true,
				},
				{
					displayMsg:    "Looks like we've found a new grinder for the hockey team! Flunk_#1 beat Flunk_#2 to a pulp.",
					challengerWon: true,
				},
				{
					displayMsg:    "What a shot! Flunk_#2 tried to dodge it, but there was no escaping Flunk_#1's speed with that twig.",
					challengerWon: true,
				},
				// Negative Outcomes
				{
					displayMsg:    "It's like Bambi on Ice out there! Flunk_#1 went to strike Flunk_#2 but missed horribly and slapped themself out cold. EMBARRASSING",
					challengerWon: false,
				},
				{
					displayMsg:    "That hockey stick is done for! Flunk_#2 just took it and broke it over their leg! Run Flunk_#1, run! I guess Flunk_#2 wins thanks to Flunk_#1's cowardness.",
					challengerWon: false,
				},
				{
					displayMsg:    "Oops! Flunk_#1 tried to twirl the **Hockey Stick** for a dramatic move, but ended up getting tangled and fell over! Talk about a rookie move.",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Skipping Rope",
			displayMsgOptions: []string{
				"**Jump ropes** in session! Flunk_#1’s looking more threatening by the second, those speed steps, and side swings are heating up the atmosphere",
				"Is that a **skipping rope** in Flunk_#1’s hands or the most lethal ligature weapon known to flunkkind?",
				"CRISS CROSS… now everybody clap your hands! Seems Flunk_#1’s looking for a fun time with a **skipping rope** in hand, not a fight!",
				"Jump, Jump, Jump, jump around… Wait, are you playing with that rope or battling, Flunk_#1?",
				"What the Flunk is Flunk_#1 doing with that **skipping rope**… ewww",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
				{
					displayMsg:    "Is Flunk_#2 supposed to look that blue?! Someone call the nurse.",
					challengerWon: true,
				},
				{
					displayMsg:    "Crouching tiger, hidden jump rope. Flunk_#1 works that rope like nun-chucks and crashes into Flunk_#2, that looked painful!",
					challengerWon: true,
				},
				{
					displayMsg:    "Flunk_#2 TAPS out, who knew skipping ropes were so efficient at strangulation",
					challengerWon: true,
				},
				// Negative Outcomes
				{
					displayMsg:    "Lord have mercy, Flunk_#1 was too busy skipping to realize that sucker punch was coming his way!",
					challengerWon: false,
				},
				{
					displayMsg:    "Turns out if you're good at skipping that doesn't make you a good fighter, oh dear oh dear",
					challengerWon: false,
				},
				{
					displayMsg:    "Flunk_#1 has tripped over their skipping rope head first into the ground… Flunk_#2 SHOUTS “how was your trip” ouch that's got to burn.",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Football",
			displayMsgOptions: []string{
				"If coach would’ve put Flunk_#1 in 4th quarter man… We would’ve been champions. *Ahem,* anyway, maybe this is vengeance for that **Football** against Flunk_#2",
				"I think Flunk_#1 just picked up a **Football** for the first time… Hey! Do you know how to throw that thing?",
				"Oh no… Flunk_#2 has it coming for them. Flunk_#1 has a CANNON with the **Football**.",
				"Get back to class Flunk_#1! You couldn’t throw a Football 15 yards to save your ass.",
				"I betcha Flunk_#1 could throw a **Football** over them mountains.",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
				{
					displayMsg:    "OOFFT, Flunk_#1 punts the football where the sun dont shine, Flunk_#2’s not going to be walking straight for days",
					challengerWon: true,
				},
				{
					displayMsg:    "WHAT AN ARM, that thing's a lethal bazooka, hits Flunk_#2 right in the face, GOODNIGHT",
					challengerWon: true,
				},
				{
					displayMsg:    "Amazing play! Flunk_#1 fakes a throw to the left and rushes right, smashing into Flunk_#2 like a star linebacker.",
					challengerWon: true,
				},
				// Negative Outcomes
				{
					displayMsg:    "Flunk_#1 attempted a deep pass, but the ball spiraled off course, giving Flunk_#2 the chance to make a sneaky counterattack.",
					challengerWon: false,
				},
				{
					displayMsg:    "Trying to recreate some NFL highlights? Flunk_#1 stumbles and fumbles the ball. Flunk_#2 snatches the opportunity, knocking Flunk_#1 off balance.",
					challengerWon: false,
				},
				{
					displayMsg:    "Epic fail. Flunk_#1 went for a Hail Mary throw, missed Flunk_#2 completely, and tumbled over from their own momentum. Not today, champ.",
					challengerWon: false,
				},
				{
					displayMsg:    "Instead of looking menacing, Flunk_#1's attempt at a quarterback sneak just got them tangled up in a mess of their own legs. Flunk_#2 is laughing too hard to continue.",
					challengerWon: false,
				},
			},
		},
		{
			weapon: "Tennis Racket",
			displayMsgOptions: []string{
				"Well, a backhand could mean two things when you’re facing off against the **Tennis Racket** wielding Flunk_#1",
				"Strawberries and Cream, there’s nothing classy going on at Flunks highschool with this maniac Flunk_#1 pulling out a dirty **tennis racket**",
				"I thought Flunk_#1 was just about to swat some flies with that tennis racket but no I’m pretty sure that's meant to do some serious damage",
				"DUECE, Flunk_#1 has whipped out their brand spanking new tennis racket and is ready for some aggressive net action…",
				"We all know Flunk_#1 has the best serve in class so I can’t wait to see the damage they’ll inflict with that **Tennis racket**",
			},
			outcomeOptions: []WeaponOutcomeOption{
				// Positive Outcomes
				{
					displayMsg:    "Ace! Flunk_#1 used their tennis racket as a catapult and sent Flunk_#2 flying out of the court.",
					challengerWon: true,
				},
				{
					displayMsg:    "Smash! Flunk_#1's powerful overhand slam caught Flunk_#2 off guard, leading to an impressive victory.",
					challengerWon: true,
				},
				{
					displayMsg:    "Using their elite racket skills, Flunk_#1 played mind games, darting left and right, leaving Flunk_#2 dazed and defeated.",
					challengerWon: true,
				},
				// Negative Outcomes
				{
					displayMsg:    "Fault! Flunk_#1 swung hard but only managed to tangle themselves in the tennis net. Advantage, Flunk_#2.",
					challengerWon: false,
				},
				{
					displayMsg:    "Flunk_#1 tried a volley but was too slow. Flunk_#2 took advantage and landed a swift counter blow.",
					challengerWon: false,
				},
				{
					displayMsg:    "Attempting a cheeky drop shot, Flunk_#1 dropped their racket instead. Easy win for Flunk_#2!",
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
				"You're not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I'm not sure how effective a **Protein Shake** is in this situation. Maybe it's psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don't work on their own, Flunk_#1. Hit the gym!",
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
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he's starting to break down and cry for some weird reason.",
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
				"You're not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I'm not sure how effective a **Protein Shake** is in this situation. Maybe it's psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don't work on their own, Flunk_#1. Hit the gym!",
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
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he's starting to break down and cry for some weird reason.",
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
				"You're not using a **Protein Shake** against Flunk_#2 are you, Flunk_#1?",
				"I'm not sure how effective a **Protein Shake** is in this situation. Maybe it's psychological for Flunk_#1, like a power boost.",
				"**Protein Shake**s don't work on their own, Flunk_#1. Hit the gym!",
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
					displayMsg:    "I guess the protein shake is giving Flunk_#1 some mood swings, he's starting to break down and cry for some weird reason.",
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
