package db

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/flunks-nft/discord-bot/utils"
	"github.com/flunks-nft/discord-bot/zeero"
	"gorm.io/gorm"
)

type Raid struct {
	ID uint

	FromTemplateID uint
	FromNftID      uint
	FromNft        Nft `gorm:"foreignKey:FromNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ToTemplateID uint
	ToNftID      uint
	ToNft        Nft `gorm:"foreignKey:ToNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// ChallengeID uint
	// Challenge   Challenge `gorm:"foreignKey:ChallengeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ChallengeID uint

	CreatedAt time.Time
	UpdatedAt time.Time
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
}

func (nft Nft) GetTraits() []Trait {
	return nft.Traits
}

func GetNftByTemplateID(templateID uint) (Nft, error) {
	var nft Nft

	result := db.Where("template_id = ?", templateID).Preload("Owner").Preload("Traits").First(&nft)
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

	nextValidRaidTime := nft.LastRaidFinishedAt.Add(24 * time.Hour)

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
		ChallengeID:    uint(rand.Intn(4) + 1),
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
