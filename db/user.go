package db

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/flunks-nft/discord-bot/zeero"
	"gorm.io/gorm"
)

type User struct {
	ID uint

	DiscordID             string
	FlowWalletAddress     string
	LastFetchedTokenIndex uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate will set a default value before every creation
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.LastFetchedTokenIndex = 0
	return
}

func (user *User) GetFlunks() ([]zeero.NftDtoWithActivity, error) {
	// Sort the items by TemplateID
	items, err := zeero.GetFlunks(user.FlowWalletAddress)
	if err != nil {
		return nil, err
	}

	// Sort the items by TemplateID
	sort.Slice(items, func(i, j int) bool {
		return items[i].TemplateID < items[j].TemplateID
	})

	return items, nil
}

func (user *User) GetNextTokenIndex(totalCount int) uint {
	lastIndex := user.LastFetchedTokenIndex
	// Increment the last fetched token index by 1 and update the user profile
	db.Model(&user).Update("last_fetched_token_index", (lastIndex+1)%uint(totalCount))
	return lastIndex
}

func CreateNewUser(DiscordID string, FlowWalletAddress string) {
	user := User{DiscordID: DiscordID, FlowWalletAddress: FlowWalletAddress}
	db.Create(&user)
}

func CreateOrUpdateFlowAddress(DiscordID string, FlowWalletAddress string) {
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
