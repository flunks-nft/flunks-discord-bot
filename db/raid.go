package db

type Raid struct {
	ID uint

	FromTokenID uint
	FromNftID   uint
	FromNft     Nft `gorm:"foreignKey:FromNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ToTokenID uint
	ToNftID   uint
	ToNft     Nft `gorm:"foreignKey:ToNftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserID uint // Foreign key referencing User's primary key
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Reference to the User struct
}

type Nft struct {
	ID uint

	TokenID uint
	Traits  []Trait `gorm:"many2many:nft_traits;"`
	Uri     string
	Points  uint
}

type Trait struct {
	ID uint

	NftID uint

	Name  string
	Value string
	Score uint
}
