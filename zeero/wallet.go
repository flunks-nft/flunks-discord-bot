package zeero

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

type NftMetadataDto struct {
	URI string `json:"uri"`
}

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
