package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocationAreas(pageUrl *string) (LocationAreasResp, error) {
	endpoint := "/location-area"
	fullURL := baseURL + endpoint

	if pageUrl != nil {
		fullURL = *pageUrl
	}

	//checking cache
	dat, ok := c.cache.Get(fullURL)

	if ok {
		// cache hit
		fmt.Println("Cache hit!")
		locationAreasResp := LocationAreasResp{}
		err := json.Unmarshal(dat, &locationAreasResp)

		if err != nil {
			return LocationAreasResp{}, err
		}
		return locationAreasResp, nil
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreasResp{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreasResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return LocationAreasResp{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	dat, err = io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResp{}, err
	}

	locationAreasResp := LocationAreasResp{}
	err = json.Unmarshal(dat, &locationAreasResp)

	if err != nil {
		return LocationAreasResp{}, err
	}

	fmt.Println("Cache miss!")
	c.cache.Add(fullURL, dat)

	return locationAreasResp, nil

}

func (c *Client) GetLocationArea(LocationAreaName string) (LocationArea, error) {
	endpoint := "/location-area/" + LocationAreaName
	fullURL := baseURL + endpoint

	//checking cache
	dat, ok := c.cache.Get(fullURL)

	if ok {
		// cache hit
		fmt.Println("Cache hit!")
		locationArea := LocationArea{}
		err := json.Unmarshal(dat, &locationArea)

		if err != nil {
			return LocationArea{}, err
		}
		return locationArea, nil
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationArea{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationArea{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return LocationArea{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	dat, err = io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, err
	}

	locationArea := LocationArea{}
	err = json.Unmarshal(dat, &locationArea)

	if err != nil {
		return LocationArea{}, err
	}

	fmt.Println("Cache miss!")
	c.cache.Add(fullURL, dat)

	return locationArea, nil

}
