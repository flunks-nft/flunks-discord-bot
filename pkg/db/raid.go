package db

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/flunks-nft/discord-bot/pkg/battle"
	"github.com/flunks-nft/discord-bot/pkg/utils"
	"github.com/flunks-nft/discord-bot/pkg/zeero"
	"gorm.io/gorm"
)

// Define the constants for the different challenge types.
const (
	ChallengeTypeGeek  ChallengeType = "Geek"
	ChallengeTypePrep  ChallengeType = "Prep"
	ChallengeTypeFreak ChallengeType = "Freak"
	ChallengeTypeJock  ChallengeType = "Jock"
)

// getRandomChallengeType returns a random challenge type.
func getRandomChallengeType() ChallengeType {
	types := [...]ChallengeType{
		ChallengeTypeGeek,
		ChallengeTypePrep,
		ChallengeTypeFreak,
		ChallengeTypeJock,
	}
	rand.Seed(time.Now().UnixNano())

	return types[rand.Intn(len(types))]
}

var (
	// TODO: replace them with the real emojis
	RAID_WON_EMOJI_ID  = utils.DiscordEmojis["RAID_WON_EMOJI_ID"]
	RAID_LOST_EMOJI_ID = utils.DiscordEmojis["RAID_LOST_EMOJI_ID"]
	RADI_WIP_EMOJI_ID  = utils.DiscordEmojis["RADI_WIP_EMOJI_ID"]
)

type Raid struct {
	ID uint

	FromTemplateID uint
	FromNftID      uint
	FromNft        Nft `gorm:"foreignKey:FromNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ToTemplateID uint
	ToNftID      uint
	ToNft        Nft `gorm:"foreignKey:ToNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ChallengeType ChallengeType

	CreatedAt time.Time
	UpdatedAt time.Time

	IsConcluded      bool `gorm:"default:false;"`
	WinnerTemplateID uint
	WinnerNftID      uint
	WinnerNft        Nft `gorm:"foreignKey:FromNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	LoserTemplateID uint
	LoserNftID      uint
	LoserNft        Nft `gorm:"foreignKey:ToNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	BattleLocation         string           `gorm:"default:'';"`
	BattleLog              battle.BattleLog `gorm:"type:jsonb"`
	BattleLogNounce        uint             `gorm:"default:0;"`
	BattleLogLastUpdatedAt time.Time
	BattleLogMessageID     string `gorm:"default:'';"`
}

func (r *Raid) SetMsgID(msgID string) {
	r.BattleLogMessageID = msgID
	// Sync database
	db.Model(&r).Select("battle_log_message_id").Update("battle_log_message_id", msgID)
}

func (r *Raid) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreatedAt = time.Now()
	r.UpdatedAt = r.CreatedAt
	return
}

func (r *Raid) BeforeUpdate(tx *gorm.DB) (err error) {
	r.BattleLogLastUpdatedAt = time.Now()
	return
}

func (raid Raid) ChallengeTypeEmoji() string {
	return utils.CliqueEmojis[raid.ChallengeType.String()]
}

func LeaderBoard() []Nft {
	// TODO: add win #, loss #, draw #, win rate
	var nfts []Nft
	result := db.Order("points DESC").Limit(10).Find(&nfts)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return nfts
}

func GetRaidByID(raidID int) Raid {
	var raid Raid
	result := db.Preload("FromNft").Preload("ToNft").Preload("WinnerNft").Preload("WinnerNft.Owner").Preload("LoserNft").Preload("LoserNft.Owner").Preload("FromNft.Owner").Preload("ToNft.Owner").First(&raid, raidID)
	if result.Error != nil {
		log.Println(result.Error)
	}

	return raid
}

// concludes a raid by
// 1. selecting winners
// 2. generating fight status text
func (raid Raid) getBattleResult() battle.BattleLog {
	return battle.DrawBattleByClique(raid.ChallengeType.String(), raid.FromTemplateID, raid.ToTemplateID, raid.BattleLocation)
}

func NextRaidToUpdateBattleStatus() (*Raid, error) {
	// Find a raid where battle_log_nounce <= 3 and BattleLogLastUpdatedAt is more than 10 seconds ago
	var raid *Raid
	result := db.Where("battle_log_nounce <= 3 AND battle_log_last_updated_at < ?", time.Now().Add(-10*time.Second)).Preload("FromNft").Preload("ToNft").Preload("WinnerNft").Preload("WinnerNft.Owner").Preload("LoserNft").Preload("LoserNft.Owner").Preload("FromNft.Owner").Preload("ToNft.Owner").First(&raid)
	if result.Error != nil {
		return raid, result.Error
	}

	return raid, nil
}

func IncrementBattleLogNounce(raid *Raid, msgID string) error {
	// Increment the nounce by 1 in database
	if err := db.Model(&raid).Select("battle_log_nounce").Update("battle_log_nounce", gorm.Expr("battle_log_nounce + ?", 1)).Error; err != nil {
		return err
	}

	// Update the last updated time in database
	if err := db.Model(&raid).Select("battle_log_last_updated_at").Update("battle_log_last_updated_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

// ConcludeRaid completes a raid and updates the NFT scores
// win: 4; draw: 2; lose: 1
func ConcludeOneRaid() (raid *Raid, err error) {
	tx := db.Begin()

	// TODO: Discord raid conclude msg is still sent even if the raid is not concluded
	// i.e. when tx is rolled back
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("Rollback occurred: %v", r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// Find the 1 raid where IsConcluded is false
	result := tx.Preload("FromNft").Preload("ToNft").Where("is_concluded = ?", false).Order("created_at ASC").First(&raid)
	if result.Error != nil {
		return raid, result.Error
	}

	// Get battle result of a raid, get Battle log and store into db
	raid.BattleLog = raid.getBattleResult()
	winner := raid.BattleLog.Winner
	if err := tx.Model(&raid).Select("battle_log").Updates(map[string]interface{}{
		"battle_log": raid.BattleLog,
	}).Error; err != nil {
		return raid, err
	}

	// Update IsConcluded to true
	if err := tx.Model(&raid).Select("is_concluded").Update("is_concluded", true).Error; err != nil {
		return raid, err
	}

	// Update LastRaidFinishedAt to current time
	if err := tx.Model(&raid.FromNft).Select("last_raid_finished_at").Update("last_raid_finished_at", time.Now()).Error; err != nil {
		return raid, err
	}
	if err := tx.Model(&raid.ToNft).Select("last_raid_finished_at").Update("last_raid_finished_at", time.Now()).Error; err != nil {
		return raid, err
	}

	// Update NFT raiding to false
	if err := tx.Model(&raid.FromNft).Select("raiding").Update("raiding", false).Error; err != nil {
		return raid, err
	}

	// Update NFT raiding to false
	if err := tx.Model(&raid.ToNft).Select("raiding").Update("raiding", false).Error; err != nil {
		return raid, err
	}

	// Determine the winner NFT and update WinnerTemplateID, WinnerNftID, LoserTemplateID, LoserNftID
	var winnerNFT, loserNFT Nft
	if winner == 0 {
		// FromNft wins (dice result is 0)
		if err := tx.Model(&raid).Updates(map[string]interface{}{
			"winner_template_id": raid.FromTemplateID,
			"winner_nft_id":      raid.FromNftID,
			"loser_template_id":  raid.ToTemplateID,
			"loser_nft_id":       raid.ToNftID,
		}).Error; err != nil {
			return raid, err
		}
		winnerNFT = raid.FromNft
		loserNFT = raid.ToNft
	} else {
		// ToNft wins (dice result is 1)
		if err := tx.Model(&raid).Updates(map[string]interface{}{
			"winner_template_id": raid.ToTemplateID,
			"winner_nft_id":      raid.ToNftID,
			"loser_template_id":  raid.FromTemplateID,
			"loser_nft_id":       raid.FromNftID,
		}).Error; err != nil {
			return raid, err
		}
		winnerNFT = raid.ToNft
		loserNFT = raid.FromNft
	}

	// Update the points of the winner and loser NFTs
	updatePoints(tx, winnerNFT, loserNFT)

	// Preload the FromNft and ToNft associations for the final raid object
	if err := tx.Preload("FromNft").Preload("ToNft").Preload("WinnerNft").Preload("WinnerNft.Owner").Preload("LoserNft").Preload("LoserNft.Owner").Preload("FromNft.Owner").Preload("ToNft.Owner").First(&raid).Error; err != nil {
		return raid, err
	}

	// If everything went well, return the updated raid with nil error
	return raid, nil
}

// updatePoints updates the points of the winner and loser NFTs.
func updatePoints(tx *gorm.DB, winner, loser Nft) {
	// Update the scores based on the outcome (win: 4; draw: 2; lose: 1)
	winner.Points += 4
	loser.Points += 1

	// Save the updated scores in the database
	tx.Model(&winner).Select("points").Update("points", winner.Points)
	tx.Model(&loser).Select("points").Update("points", loser.Points)
}

func GetRaidHistoryByTemplateID(templateID uint) []string {
	var raids []Raid
	result := db.Where("from_template_id = ? OR to_template_id = ?", templateID, templateID).Preload("FromNft").Preload("ToNft").Preload("WinnerNft").Preload("LoserNft").Find(&raids).Order("created_at DESC").Limit(30)
	if result.Error != nil {
		log.Println(result.Error)
	}

	records := make([]string, 0)

	// Loop through the raids and create a record for each raid.
	// TODO: update the emoji based on raid result
	for _, raid := range raids {
		emoji := fmt.Sprintf("<:emoji:%s>", RAID_WON_EMOJI_ID) // Default to the spark emoji.

		if !raid.IsConcluded {
			emoji = fmt.Sprintf("<:emoji:%s>", RADI_WIP_EMOJI_ID) // Set the WIP emoji if the raid is still in progress.
		}

		if raid.WinnerNft.TemplateID != templateID {
			// Loser emoji granted to the loser
			emoji = fmt.Sprintf("<:emoji:%s>", RAID_LOST_EMOJI_ID)
		}

		timestamp := raid.CreatedAt.Format("2006-01-02") // Format the timestamp as needed.
		record := fmt.Sprintf(
			"%s [%s] %s Flunk #%d ⚔️ Flunk #%d \n",
			emoji,
			timestamp,
			raid.ChallengeTypeEmoji(),
			raid.FromNft.TemplateID,
			raid.ToNft.TemplateID,
		)
		records = append(records, record)
	}

	return records
}

type Nft struct {
	ID uint

	TokenID    uint
	TemplateID uint
	Traits     []Trait `gorm:"many2many:nft_traits;"`
	Uri        string
	Points     uint

	CreatedAt time.Time
	UpdatedAt time.Time

	OwnerUserId uint
	Owner       User `gorm:"foreignKey:OwnerUserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	LastRaidFinishedAt time.Time
	Raiding            bool
	QueuedForRaiding   bool

	// Add a reference to the raids where this NFT is the "from" NFT
	FromRaids []Raid `gorm:"foreignKey:FromNftID"`

	// Add a reference to the raids where this NFT is the "to" NFT
	ToRaids []Raid `gorm:"foreignKey:ToNftID"`
}

// Note: the smaller grade, the stronger the NFT
func (nft Nft) Grade() uint {
	// Load NFT traits if not loaded
	if len(nft.Traits) == 0 {
		if err := db.Model(&nft).Association("Traits").Find(&nft.Traits); err != nil {
			// Handle error?
		}
	}

	var total uint

	// sum up power by traits and return the sum
	for _, trait := range nft.Traits {
		count, ok := utils.TraitValueToCount[trait.Value]
		if !ok || count <= 0 {
			log.Printf("Error: Invalid count for trait value %s", trait.Value)
			continue
		}
		total += count
	}
	return total
}

func (nft Nft) GetTraits() []Trait {
	return nft.Traits
}

func (nft Nft) GetWinsAndLosses() (int64, int64) {
	// Count the number of wins and losses from raid model
	var wins, losses int64

	// Preload FromRaids and ToRaids
	if err := db.Model(nft).Association("FromRaids").Find(&nft.FromRaids); err != nil {
		// Handle error
	}
	if err := db.Model(nft).Association("ToRaids").Find(&nft.ToRaids); err != nil {
		// Handle error
	}

	// Count wins and losses from the NFT's "FromRaids"
	for _, raid := range nft.FromRaids {
		if raid.IsConcluded {
			if raid.WinnerNftID == nft.ID {
				wins++
			} else if raid.LoserNftID == nft.ID {
				losses++
			}
		}
	}

	// Count wins and losses from the NFT's "ToRaids"
	for _, raid := range nft.ToRaids {
		if raid.IsConcluded {
			if raid.WinnerNftID == nft.ID {
				wins++
			} else if raid.LoserNftID == nft.ID {
				losses++
			}
		}
	}

	return wins, losses
}

func GetNftByTemplateID(templateID uint) (Nft, error) {
	var nft Nft

	result := db.Where("template_id = ?", templateID).Preload("Owner").Preload("Traits").Preload("FromRaids").Preload("ToRaids").First(&nft)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nft, fmt.Errorf("NFT not found")
		}
		return nft, result.Error
	}
	return nft, nil
}

func CreateNft(user *User, tokenID uint, templateID uint, uri string, traits []Trait) (Nft, error) {
	nft := Nft{
		TokenID:     tokenID,
		TemplateID:  templateID,
		Uri:         uri,
		OwnerUserId: user.ID,
	}

	db.Create(&nft)
	db.Model(&nft).Association("Traits").Append(traits)

	nft, err := GetNftByTemplateID(templateID)
	if err != nil {
		return nft, err
	}

	return nft, nil
}

// TODO: optimize the logic for CreateOrUpdateFlunks
func CreateOrUpdateFlunks(user *User, flunks []zeero.NftDtoWithActivity) error {
	for _, flunk := range flunks {
		// Create Traits
		traits, err := CreateTraits(flunk.Metadata)
		if err != nil {
			return err
		}

		// Get NFT instance from database
		_, err = GetNftByTemplateID(uint(flunk.TemplateID))
		if err != nil {
			CreateNft(user, uint(flunk.TokenID), uint(flunk.TemplateID), flunk.Metadata.URI, traits)
		}
	}
	return nil
}

func (nft *Nft) IsInRaidQueue() bool {
	return nft.QueuedForRaiding
}

func (nft *Nft) IsRaiding() bool {
	return nft.Raiding
}

func (nft *Nft) IsReadyForRaidQueue() (bool, time.Duration) {
	// Get current time
	now := time.Now()

	// Calculate the time for the next day
	nextDay := nft.LastRaidFinishedAt.AddDate(0, 0, 1)

	// If the current time is greater than or equal to the next day,
	// return true indicating that it's ready for the raid queue
	// and return the time until 12:00 am next day as the second return value
	if now.After(nextDay) || now.Equal(nextDay) {
		return true, 0
	}

	timeUntilNextDay := nextDay.Sub(now)
	return false, timeUntilNextDay
}

func (nft *Nft) AddToRaidQueue() error {
	// Update NFT QueuedForRaiding to true in database
	result := db.Model(&nft).Updates(Nft{QueuedForRaiding: true})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetNextQueuedToken(tx *gorm.DB) (*Nft, error) {
	database := tx
	if database == nil {
		database = db
	}

	nft := &Nft{}
	err := db.Where("queued_for_raiding = ?", true).Order("RANDOM()").First(nft).Error
	if err != nil {
		return nil, err
	}
	return nft, nil
}

// GetNextQueuedTokenPair gets the next 2 available tokens pairs
// that are queued for raiding, including their traits
func GetNextQueuedTokenPair(tx *gorm.DB) ([]Nft, error) {
	database := tx
	if database == nil {
		database = db
	}

	var nfts []Nft
	subQuery := database.
		Select("DISTINCT ON (owner_user_id) *").
		Where("nfts.queued_for_raiding = ?", true).
		Order("owner_user_id, RANDOM()").
		Table("nfts")
	err := database.
		Preload("Owner").
		Preload("Traits").
		Joins("JOIN (?) as distinct_nfts on distinct_nfts.id = nfts.id", subQuery).
		Limit(2).
		Find(&nfts).
		Error
	if err != nil {
		return nil, err
	}

	if len(nfts) != 2 {
		return nil, fmt.Errorf("No available pair - has %v", len(nfts))
	}

	return nfts, nil
}

func QueueNextTokenPairForRaiding() (*Raid, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	nfts, err := GetNextQueuedTokenPair(tx)

	if err != nil {
		tx.Rollback()
		log.Println(err)
		return nil, err
	}

	// Pick a random key from the map as battle location
	challengeType := getRandomChallengeType()
	battleBgImgMap, _ := utils.BattleBgImages[challengeType.String()]
	battleLocation := utils.GetRandomKeyFromMap(battleBgImgMap)

	raid := &Raid{
		FromTemplateID: nfts[0].TemplateID,
		FromNftID:      nfts[0].ID,
		ToTemplateID:   nfts[1].TemplateID,
		ToNftID:        nfts[1].ID,
		ChallengeType:  challengeType,
		BattleLocation: battleLocation,
	}

	result := tx.Create(&raid).Preload("FromNft").Preload("ToNft").Preload("FromNft.Owner").Preload("ToNft.Owner")
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Select raid from the db
	result = tx.Preload("FromNft").Preload("ToNft").Preload("FromNft.Owner").Preload("ToNft.Owner").First(&raid)

	// Mark both nfts as Raiding and set QueuedForRaiding to false
	nftUpdateData := map[string]interface{}{
		"Raiding":          true,
		"QueuedForRaiding": false,
	}
	tx.Model(&nfts[0]).Updates(nftUpdateData)
	tx.Model(&nfts[1]).Updates(nftUpdateData)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	tx.Commit()
	return raid, nil
}

type Trait struct {
	ID uint

	Name  string
	Value string
	Score uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (trait Trait) create() error {
	result := db.Create(&trait)
	if result.Error != nil {
		// silently fail
	}
	return nil
}

func CreateTraits(metadata zeero.NftMetadataDto) ([]Trait, error) {
	traits := metadata.Traits()
	dbTraits := make([]Trait, 0)

	for _, trait := range traits {
		dbTrait := Trait{Name: trait.Name, Value: trait.Value}
		if score, found := utils.TraitToScore[trait.Value]; found {
			dbTrait.Score = score
		} else {
			dbTrait.Score = 0 // Set a default value (0) if the traitValue is not found in the map.
		}
		dbTraits = append(dbTraits, dbTrait)

		// Add trait in the database
		dbTrait.create()
	}

	return dbTraits, nil
}
