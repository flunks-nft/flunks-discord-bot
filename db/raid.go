package db

import (
	"errors"
	"fmt"
	"time"

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

	ChallengeID uint
	Challenge   Challenge `gorm:"foreignKey:ChallengeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserID uint // Foreign key referencing User's primary key
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Reference to the User struct

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

func GetNft(TemplateID uint) (Nft, error) {
	var nft Nft

	result := db.Where("template_id = ?", TemplateID).First(&nft)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nft, fmt.Errorf("NFT not found")
		}
		return nft, result.Error
	}
	return nft, nil
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

type Trait struct {
	ID uint

	NftID uint

	Name  string
	Value string
	Score uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Challenge struct {
	ID uint

	// Challenge is mapped to Traitm only the "Clique" trait should be used
	TraitID uint
	Trait   Trait `gorm:"foreignKey:TraitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
