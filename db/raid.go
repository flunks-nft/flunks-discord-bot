package db

type Property string

// @dev: These are the properties that can be assigned to a trait.
// the properties have special meaning in the game mechanics.
// for instance, the "A" property counters the "B" property, which counters the "C" property, which counters the "A" property.
// TODO: implement proper names for the properties
const (
	PropertyA Property = "A"
	PropertyB Property = "B"
	PropertyC Property = "C"
	PropertyD Property = "E"
)

type Raid struct {
	ID uint

	FromTokenID     uint
	ToTokenID       uint
	FromTokenPoints uint
	ToTokenPoints   uint

	UserID uint // Foreign key referencing User's primary key
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Reference to the User struct
}

type Nft struct {
	TokenID uint
	Trait   Trait      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Reference to the Trait struct
	Display NftDisplay `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Reference to the Trait struct
}

type NftDisplay struct {
	Uri string
}

type Trait struct {
	Name     string
	Value    string
	Property Property
}
