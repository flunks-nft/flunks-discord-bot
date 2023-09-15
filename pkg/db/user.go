package db

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/flunks-nft/discord-bot/pkg/zeero"
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

	CreateOrUpdateFlunks(user, items)

	// Sort the items by TemplateID
	sort.Slice(items, func(i, j int) bool {
		return items[i].TemplateID < items[j].TemplateID
	})

	return items, nil
}

func (user *User) OwnsFlunk(templateID int) bool {
	// Get the flunks owned by the user
	flunks, err := user.GetFlunks()
	if err != nil {
		return false
	}

	// Check if the user owns the flunk with the specified template ID
	for _, flunk := range flunks {
		if flunk.TemplateID == templateID {
			return true
		}
	}

	return false
}

func (user *User) GetNextTokenIndex(totalCount int) uint {
	lastIndex := user.LastFetchedTokenIndex
	// Increment the last fetched token index by 1 and update the user profile
	db.Model(&user).Update("last_fetched_token_index", (lastIndex+1)%uint(totalCount))
	return lastIndex % uint(totalCount)
}

func (user *User) ForceUpdateNextTokenIndex(totalCount int, targetIndex int) {
	// Update the last fetched token index to a provided value
	db.Model(&user).Update("last_fetched_token_index", uint(targetIndex)%uint(totalCount))
}

func CreateNewUser(DiscordID string, FlowWalletAddress string) error {
	user := User{DiscordID: DiscordID, FlowWalletAddress: FlowWalletAddress}
	result := db.Create(&user)
	return result.Error
}

func CreateOrUpdateFlowAddress(DiscordID string, FlowWalletAddress string) error {
	user, err := UserProfile(DiscordID)
	if err != nil {
		// If UserProfile returns an error when a user is not found
		// Create a new user
		err = CreateNewUser(DiscordID, FlowWalletAddress)
		if err != nil {
			return err
		}
	} else {
		// User exists, update the Flow wallet address
		err = user.UpdateFlowAddress(FlowWalletAddress)
		if err != nil {
			return err
		}
	}

	return nil
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

func (user *User) UpdateFlowAddress(FlowWalletAddress string) error {
	result := db.Model(user).Update("FlowWalletAddress", FlowWalletAddress)
	return result.Error
}

func (user *User) GetTokenIds() []string {
	// TODO: implement with correct logic
	return []string{"1", "2", "3"}
}
