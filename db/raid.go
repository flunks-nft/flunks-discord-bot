package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Raid struct {
	ID uint

	FromTokenID uint
	FromNftID   uint
	FromNft     Nft `gorm:"foreignKey:FromNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ToTokenID uint
	ToNftID   uint
	ToNft     Nft `gorm:"foreignKey:ToNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ChallengeID uint
	Challenge   Challenge `gorm:"foreignKey:ChallengeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserID uint // Foreign key referencing User's primary key
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Reference to the User struct

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Nft struct {
	ID uint

	TokenID uint
	Traits  []Trait `gorm:"many2many:nft_traits;"`
	Uri     string
	Points  uint

	CreatedAt time.Time
	UpdatedAt time.Time

	lastRaidFinishedAt time.Time
	isRaiding          bool
	isQueuedForRaiding bool
}

func GetNft(tokenID uint) (Nft, error) {
	var nft Nft

	result := db.Where("token_id = ?", tokenID).First(&nft)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nft, fmt.Errorf("NFT not found")
		}
		return nft, result.Error
	}
	return nft, nil
}

func (nft *Nft) IsInRaidQueue() bool {
	return nft.isQueuedForRaiding
}

func (nft *Nft) IsRaiding() bool {
	return nft.isRaiding
}

func (nft *Nft) IsReadyForRaidQueue() (bool, time.Duration) {
	// Get current time
	now := time.Now()

	nextValidRaidTime := nft.lastRaidFinishedAt.Add(24 * time.Hour)

	// If the current time is less than the next valid raid time, return false
	// and the hours remaining until the next valid raid time
	if now.Before(nextValidRaidTime) {
		diff := nextValidRaidTime.Sub(now)
		return false, diff
	}

	return true, 0
}

func (nft *Nft) QueueForRaid() error {
	// Update NFT isQueuedForRaiding to true & isRaiding to false in database
	result := db.Model(&nft).Updates(Nft{isQueuedForRaiding: true, isRaiding: false})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("NFT not found")
		}
		return result.Error
	}
	return nil
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
