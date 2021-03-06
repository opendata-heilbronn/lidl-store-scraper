package client

import (
	"encoding/json"
	"net/http"
)

type StoreClient struct {
	Url string
}

type StoreResult struct {
	Store   string `json:"Store"`
	Country string `json:"Country"`
	// Bad practise warning: zip codes are numeric, not numbers, use string if possible
	ZipCode      int64   `json:"ZIP"`
	City         string  `json:"City"`
	Street       string  `json:"Street"`
	OpeningHours string  `json:"Opening Hours"`
	Longitude    float64 `json:"X Coordinate"`
	Latitude     float64 `json:"Y Coordinate"`
	ObjectType   string  `json:"Object Type"`
}

func NewStoreClient(url string) *StoreClient {
	return &StoreClient{
		Url: url,
	}
}

func (c *StoreClient) Scrape() ([]StoreResult, error) {
	resp, err := http.Get(c.Url)
	if err != nil {
		return nil, err
	}

	var jsonResult []StoreResult
	err = json.NewDecoder(resp.Body).Decode(&jsonResult)
	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}
