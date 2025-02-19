package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ApiGoogleMapDriver struct{}

type Location struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type GeolocationApiResponse struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	Accuracy float64 `json:"accuracy"`
}

type PlacesApiResponse struct {
	Places []struct {
		Id          string `json:"id"`
		DisplayName struct {
			Text string `json:"text"`
		} `json:"displayName"`
		RegularOpeningHours struct {
			WeekdayDescriptions []string `json:"weekdayDescriptions"`
		} `json:"regularOpeningHours"`
		PriceLevel string   `json:"priceLevel"`
		Location   Location `json:"location"`
	} `json:"places"`
}

type Store struct {
	Id                  string   `json:"places.id"`
	Name                string   `json:"places.displayName.text"`
	RegularOpeningHours []string `json:"places.regularOpeningHours.weekdayDescriptions"`
	PriceLevel          string   `json:"places.priceLevel"`
	Location            Location `json:"places.location"`
}

func NewGoogleMapDriver() *ApiGoogleMapDriver {
	return &ApiGoogleMapDriver{}
}

func (ap *ApiGoogleMapDriver) GetStores() ([]*Store, error) {
	location, err := getCurrentLocation()
	if err != nil {
		fmt.Println("Error:", err)
		return make([]*Store, 0), ErrFailedToGetLocation
	}
	stores, err := searchStoresNearby(location)
	if err != nil {
		fmt.Println("Error:", err)
		return make([]*Store, 0), ErrFailedToGetNearbyStores
	}
	return stores, nil
}

func getCurrentLocation() (Location, error) {
	resp, err := http.Post(
		"https://www.googleapis.com/geolocation/v1/geolocate?key="+os.Getenv("GOOGLE_MAP_API_KEY"),
		"application/json",
		nil,
	)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}
	var response GeolocationApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Location{}, err
	}
	location := Location{
		Lat: response.Location.Lat,
		Lng: response.Location.Lng,
	}
	return location, nil
}

func searchStoresNearby(location Location) ([]*Store, error) {
	requestBody := fmt.Sprintf(`{
		"includedTypes": ["cafe", "restaurant"],
		"maxResultCount": 10,
		"languageCode": "ja",
		"regionCode": "JP",
		"locationRestriction": {
			"circle": {
				"center": {"latitude": "%f", "longitude": "%f"},
				"radius": 500.0
			}
		}
	}`, location.Lat, location.Lng)
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://places.googleapis.com/v1/places:searchNearby",
		bytes.NewBuffer([]byte(requestBody)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", os.Getenv("GOOGLE_MAP_API_KEY"))
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.regularOpeningHours.weekdayDescriptions,places.priceLevel,places.location")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var placeResponse PlacesApiResponse
	err = json.Unmarshal(body, &placeResponse)
	if err != nil {
		return nil, err
	}
	stores := make([]*Store, 0)
	for _, place := range placeResponse.Places {
		stores = append(stores, &Store{
			Id:                  place.Id,
			Name:                place.DisplayName.Text,
			RegularOpeningHours: place.RegularOpeningHours.WeekdayDescriptions,
			PriceLevel:          place.PriceLevel,
			Location: Location{
				Lat: place.Location.Lat,
				Lng: place.Location.Lng,
			},
		})
	}
	return stores, nil
}
