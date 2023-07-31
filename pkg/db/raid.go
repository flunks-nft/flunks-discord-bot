package db

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

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
	RAID_CONCLUDE_TIME = 10 * time.Second

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

// ConcludeRaid completes a raid and updates the NFT scores
// win: 4; draw: 2; lose: 1
func ConcludeOneRaid() (raid Raid, err error) {
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// Find the 1 raid where createdAt is more than <RAID_CONCLUDE_TIME> hours ago and is the oldest, and IsConcluded is false
	result := tx.Where("is_concluded = ? AND created_at < ?", false, time.Now().UTC().Add(-RAID_CONCLUDE_TIME)).Order("created_at ASC").First(&raid)
	if result.Error != nil {
		return raid, result.Error
	}

	// Roll the dice to determine the winner (0 or 1)
	winner := rand.Intn(2) // Randomly generates 0 or 1

	// Update IsConcluded to true
	if err := tx.Model(&raid).Update("is_concluded", true).Error; err != nil {
		return raid, err
	}

	// Determine the winner NFT and update WinnerTemplateID, WinnerNftID
	var winnerNFT, loserNFT Nft
	if winner == 0 {
		// FromNft wins (dice result is 0)
		if err := tx.Model(&raid).Updates(map[string]interface{}{"winner_template_id": raid.FromTemplateID, "winner_nft_id": raid.FromNftID}).Error; err != nil {
			return raid, err
		}
		winnerNFT = raid.FromNft
		loserNFT = raid.ToNft
	} else {
		// ToNft wins (dice result is 1)
		if err := tx.Model(&raid).Updates(map[string]interface{}{"winner_template_id": raid.ToTemplateID, "winner_nft_id": raid.ToNftID}).Error; err != nil {
			return raid, err
		}
		winnerNFT = raid.ToNft
		loserNFT = raid.FromNft
	}

	// Update the scores of the winner and loser NFTs
	updateScores(tx, winnerNFT, loserNFT)

	// Preload the FromNft and ToNft associations for the final raid object
	if err := tx.Preload("FromNft").Preload("ToNft").Preload("WinnerNft").Preload("FromNft.Owner").Preload("ToNft.Owner").Preload("WinnerNft.Owner").First(&raid).Error; err != nil {
		return raid, err
	}

	// If everything went well, return the updated raid with nil error
	return raid, nil
}

// updateScores updates the scores of the winner and loser NFTs.
func updateScores(tx *gorm.DB, winner, loser Nft) {
	// Update the scores based on the outcome (win: 4; draw: 2; lose: 1)
	winner.Points += 4
	loser.Points += 1

	// Save the updated scores in the database
	tx.Model(&winner).Update("points", winner.Points)
	tx.Model(&loser).Update("points", loser.Points)
}

func GetRaidHistoryByTemplateID(tokenID uint) []string {
	var raids []Raid
	result := db.Where("from_template_id = ? OR to_template_id = ?", tokenID, tokenID).Preload("FromNft").Preload("ToNft").Find(&raids).Order("created_at DESC").Limit(30)
	if result.Error != nil {
		log.Println(result.Error)
	}

	records := make([]string, 0)

	// Loop through the raids and create a record for each raid.
	// TODO: update the emoji based on raid result
	for _, raid := range raids {
		emoji := fmt.Sprintf("<:emoji:%s>", RAID_WON_EMOJI_ID) // Default to the spark emoji.

		// Check if the raid was concluded before RAID_CONCLUDE_TIME ago compared to the current time.
		concludeTimeAgo := time.Since(raid.CreatedAt)
		if concludeTimeAgo < RAID_CONCLUDE_TIME {
			emoji = fmt.Sprintf("<:emoji:%s>", RADI_WIP_EMOJI_ID) // Set the WIP emoji if the raid is still in progress.
		}

		timestamp := raid.CreatedAt.Format("2006-01-02") // Format the timestamp as needed.
		record := fmt.Sprintf(
			"%s %s game %s Flunk #%d ⚔️ Flunk #%d %s %s \n",
			emoji,
			raid.ChallengeType,
			raid.ChallengeTypeEmoji(),
			raid.FromNft.TemplateID,
			raid.ToNft.TemplateID,
			"⏰",
			timestamp,
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

func (nft Nft) GetTraits() []Trait {
	return nft.Traits
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

	nextValidRaidTime := nft.LastRaidFinishedAt.Add(RAID_CONCLUDE_TIME)

	// If the current time is less than the next valid raid time, return false
	// and the hours remaining until the next valid raid time
	if now.Before(nextValidRaidTime) {
		diff := nextValidRaidTime.Sub(now)
		return false, diff
	}

	return true, 0
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
	err := database.
		// Preload the User and Traits for each Nft
		Preload("Owner").
		Preload("Traits").
		Where("nfts.queued_for_raiding = ?", true).
		Order("RANDOM()").
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

func QueueNextTokenPairForRaiding() (*Raid, []Nft, error) {
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
		return nil, nil, err
	}

	raid := &Raid{
		FromTemplateID: nfts[0].TemplateID,
		FromNftID:      nfts[0].ID,
		ToTemplateID:   nfts[1].TemplateID,
		ToNftID:        nfts[1].ID,
		ChallengeType:  getRandomChallengeType(),
	}

	result := tx.Create(&raid)
	if result.Error != nil {
		tx.Rollback()
		return nil, nil, result.Error
	}

	// Mark both nfts as Raiding and set QueuedForRaiding to false
	nftUpdateData := map[string]interface{}{
		"Raiding":          true,
		"QueuedForRaiding": false,
	}
	result = tx.Model(&nfts).Updates(nftUpdateData)
	if result.Error != nil {
		tx.Rollback()
		return nil, nil, result.Error
	}

	tx.Commit()
	return raid, nfts, nil
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
