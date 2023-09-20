package utils

// TraitToScore is a map that associates trait names with scores.
var TraitToScore = map[string]uint{
	"Superlative": 1,
}

var TraitNameToEmoji = map[string]string{
	"Clique":      "👯",
	"Face":        "😀",
	"Torso":       "👕",
	"Head":        "🧠",
	"Pigment":     "🎨",
	"Backdrop":    "🏞️",
	"Type":        "🎓",
	"Superlative": "🏆",
}

var DiscordEmojis = map[string]string{
	"RAID_WON_EMOJI_ID":  "1144740476117860442",
	"RAID_LOST_EMOJI_ID": "1144752213760163910",
	"RADI_WIP_EMOJI_ID":  "1144752263697530961",
}

var CliqueEmojis = map[string]string{
	"Prep":  "👔",
	"Jock":  "🏈",
	"Freak": "👽",
	"Geek":  "🎮",
}

var BattleBgImages = map[string]string{
	"Prep":  "https://storage.googleapis.com/zeero-public/prep_battle.png",
	"Jock":  "https://storage.googleapis.com/zeero-public/jock_battle.png",
	"Freak": "https://storage.googleapis.com/zeero-public/freak_battle.png",
	"Geek":  "https://storage.googleapis.com/zeero-public/geek_battle.png",
}
