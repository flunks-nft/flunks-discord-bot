package db

import "time"

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

	LastRaidFinishedAt time.Time
	isRaiding          bool
	isQueuedForRaiding bool
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
