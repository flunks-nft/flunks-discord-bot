package zeero

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/flunks-nft/discord-bot/utils"
)

var (
	ZEERO_API_URL string
)

func init() {
	utils.LoadEnv()

	ZEERO_API_URL = os.Getenv("ZEERO_API_URL")
}

type GetUserWalletNftsDto struct {
	Flunks          []NftDtoWithActivity `json:"Flunks"`
	Backpack        []NftDtoWithActivity `json:"Backpack"`
	Patch           []NftDtoWithActivity `json:"Patch"`
	InceptionAvatar []NftDtoWithActivity `json:"InceptionAvatar"`
}

type NftDtoWithActivity struct {
	TokenID    int            `json:"tokenId"`
	TemplateID int            `json:"templateId"`
	Metadata   NftMetadataDto `json:"metadata"`
}

type Trait struct {
	Name  string
	Value string
	Score uint
}

func (metadata NftMetadataDto) Traits() []Trait {
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

		// only Graduated Flunks have a Type trait
		// so skip the Type trait if it's empty
		if traitName == "Type" && traitValue == "" {
			continue
		}

		trait := Trait{
			Name:  traitName,
			Value: traitValue,
		}

		traits = append(traits, trait)
	}

	return traits
}

type NftMetadataDto struct {
	URI string `json:"uri"`

	Clique      string `json:"Clique"`
	Face        string `json:"Face"`
	Torso       string `json:"Torso"`
	Head        string `json:"Head"`
	Pigment     string `json:"Pigment"`
	Backdrop    string `json:"Backdrop"`
	Type        string `json:"Type"`
	Superlative string `json:"Superlative"`
}

// GetFlunks gets all the Flunks and their metadata from the Zeero API
// TODO: it's currently 50 Flunks, but we need to paginate
func GetFlunks(address string) ([]NftDtoWithActivity, error) {
	// Send an HTTP GET request
	url := fmt.Sprintf("%s/users/%s", ZEERO_API_URL, address)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Deserialize the JSON payload
	var data GetUserWalletNftsDto
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return data.Flunks, nil
}
