package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID uint

	DiscordID         string
	FlowWalletAddress string
	Points            uint
}

func CreateNewUser(DiscordID string, FlowWalletAddress string) {
	user := User{DiscordID: DiscordID, FlowWalletAddress: FlowWalletAddress, Points: 0}
	db.Create(&user)
}

func UpdateFlowAddress(DiscordID string, FlowWalletAddress string) {
	user, err := UserProfile(DiscordID)
	// User doesn't exist, create a new one
	if err != nil {
		CreateNewUser(DiscordID, FlowWalletAddress)
	}

	// User exists, update the Flow wallet address
	user.UpdateFlowAddress(FlowWalletAddress)
}

// UserProfile retrieves the user profile based on the Discord ID from the provided database connection.
// It performs a query to find the user with the specified Discord ID and returns the user profile.
// If the user is found, the user profile is returned along with a nil error.
// If the user is not found, it returns an empty User struct and an error with the message "User not found".
// For any other database query errors, it returns the corresponding error.
func UserProfile(DiscordID string) (User, error) {
	var user User

	result := db.Where("discord_id = ?", DiscordID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("User not found")
		}
		return user, result.Error
	}
	return user, nil
}

func (user *User) UpdateFlowAddress(FlowWalletAddress string) {
	db.Model(&user).Update("FlowWalletAddress", FlowWalletAddress)
}

func (user *User) GetTokenIds() []string {
	// TODO: implement with correct logic
	return []string{"1", "2", "3"}
}

type Raid struct {
	ID uint

	FromTokenID     uint
	ToTokenID       uint
	FromTokenPoints uint
	ToTokenPoints   uint

	UserID uint // Foreign key referencing User's primary key
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Reference to the User struct
}
