package db

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
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

	LastRaidFinishedAt time.Time
	Raiding            bool
	QueuedForRaiding   bool
}

func (nft Nft) GetTraits() []Trait {
	return nft.Traits
}

func GetNftByTemplateID(templateID uint) (Nft, error) {
	var nft Nft

	result := db.Where("template_id = ?", templateID).First(&nft)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nft, fmt.Errorf("NFT not found")
		}
		return nft, result.Error
	}
	return nft, nil
}

func CreateNft(tokenID uint, templateID uint, uri string, traits []Trait) (Nft, error) {
	nft := Nft{
		TokenID:    tokenID,
		TemplateID: templateID,
		Uri:        uri,
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
func CreateOrUpdateFlunks(flunks []zeero.NftDtoWithActivity) error {
	for _, flunk := range flunks {
		// Create Traits
		traits, err := CreateTraitsForFlunks(flunk.Metadata)
		if err != nil {
			return err
		}

		// Get NFT instance from database
		_, err = GetNftByTemplateID(uint(flunk.TemplateID))
		if err != nil {
			CreateNft(uint(flunk.TokenID), uint(flunk.TemplateID), flunk.Metadata.URI, traits)
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

func (nft *Nft) RaidMatch() error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the next available token
	availableToken, err := GetNextQueuedToken(tx)
	if err != nil {
		tx.Rollback()
		return err
	} else if availableToken == nil {
		// Update NFT QueuedForRaiding to true & Raiding to false in database
		result := tx.Model(&nft).Updates(Nft{QueuedForRaiding: true, Raiding: false})
		if result.Error != nil {
			tx.Rollback()
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("NFT not found")
			}
			return result.Error
		}
	} else {
		// Create a raid
		raid := &Raid{
			FromTemplateID: nft.TemplateID,
			FromNftID:      nft.ID,
			ToTemplateID:   availableToken.TemplateID,
			ToNftID:        availableToken.ID,
		}
		result := tx.Create(&raid)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		// Update both NFTs Raiding to true & QueuedForRaiding to false in database
		result = tx.Model(&nft).Updates(Nft{QueuedForRaiding: false, Raiding: true})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	tx.Commit()
	return nil
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

func GetNextQueuedTokenPair(tx *gorm.DB) ([]Nft, error) {
	// !!TODO: return NFT traits as well
	database := tx
	if database == nil {
		database = db
	}

	var nfts []Nft
	err := database.Where("queued_for_raiding = ?", true).Order("RANDOM()").Limit(2).Find(&nfts).Error
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

func (trait Trait) CreateTrait() error {
	result := db.Create(&trait)
	if result.Error != nil {
		// silently fail
	}
	return nil
}

func CreateTraitsForFlunks(metadata zeero.NftMetadataDto) ([]Trait, error) {
	// Use reflection to iterate over the fields of the struct
	types := reflect.TypeOf(metadata)
	values := reflect.ValueOf(metadata)

	traits := make([]Trait, 0)

	for i := 0; i < types.NumField(); i++ {
		field := types.Field(i)
		value := values.Field(i)

		traitName := field.Name
		traitValue := value.Interface().(string)

		if traitName == "URI" {
			continue
		}

		trait := Trait{
			Name:  traitName,
			Value: traitValue,
		}

		if score, found := utils.TraitToScore[traitValue]; found {
			trait.Score = score
		} else {
			trait.Score = 0 // Set a default value (0) if the traitValue is not found in the map.
		}

		err := trait.CreateTrait()
		if err != nil {
			return traits, err
		}

		traits = append(traits, trait)
	}

	return traits, nil
}

// type Challenge struct {
// 	ID uint

// 	// Challenge is mapped to Traitm only the "Clique" trait should be used
// 	TraitID uint
// 	Trait   Trait `gorm:"foreignKey:TraitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// }
