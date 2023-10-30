package battle

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/flunks-nft/discord-bot/pkg/gpt"
)

// GenerateBattleLog interacts with OpenAI's GPT-3.5Turbo API to generate a battle log
func GenerateBattleLog(clique string, challenger, defender uint, location string, winner uint) (*BattleLog, error) {
	ctx := context.Background()
	isPositiveOutcome := winner == 0

	prompt := gpt.GenerateBattlePrompt(clique, challenger, defender, location, isPositiveOutcome)

	res, err := gpt.ChapGPTClient.SimpleSend(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var log BattleLog
	err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &log)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

type BattleLog struct {
	Weapon        string `json:"weapon"`
	Action        string `json:"action"`
	ActionOutcome string `json:"actionOutcome"`
	BattleOutcome string `json:"battleOutcome"`
	Winner        int    `json:"winner"`
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

func DrawBattleByClique(clique string, challenger, defender uint, location string, winner uint) BattleLog {
	log, err := GenerateBattleLog(clique, challenger, defender, location, winner)
	if err != nil {
		panic(err)
	}
	return *log
}
