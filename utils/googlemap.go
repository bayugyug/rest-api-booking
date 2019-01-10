package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	GoogleApiUri = "https://maps.googleapis.com/maps/api/geocode/json?key={key}&address="
	GoogleApiKey = "AIzaSyCQWCoQRJqlsMuBGnX3bi3LjUISFY3xG9o"
)

type GoogleApiResponse struct {
	Results Results `json:"results"`
	Status  string  `json:"status"`
}

type Results []Geometry

type Geometry struct {
	Geometry         Location `json:"geometry"`
	FormattedAddress string   `json:"formatted_address"`
}

type Location struct {
	Location Coordinates `json:"location"`
}

type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type GoogleMapGeoCode struct {
	Key     string `json:"key"`
	Address string `json:"address"`
}

//NewGoogleMapGeoCode geo-code
func NewGoogleMapGeoCode(key string) *GoogleMapGeoCode {
	return &GoogleMapGeoCode{
		Key: key,
	}
}

//GetCoordinates get the geo code info as per address
func (g *GoogleMapGeoCode) GetCoordinates(address string) (bool, float64, float64, string, error) {
	gpoint := strings.Replace(GoogleApiUri+url.QueryEscape(address), "{key}", g.Key, -1)
	body, code, err := httpGet(gpoint, map[string]string{})
	if err != nil {
		log.Println("GMAP_API: failed", err)
		return false, 0, 0, "", err
	}
	if code != http.StatusOK {
		log.Println("GMAP_API: failed", code)
		return false, 0, 0, "", err
	}
	var data GoogleApiResponse
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Println("GMAP_API: failed", code)
		return false, 0, 0, "", err
	}
	if len(data.Results) <= 0 || data.Status != "OK" {
		log.Println("GMAP_API: failed", data.Status)
		return false, 0, 0, "", err
	}
	return true, data.Results[0].Geometry.Location.Latitude, data.Results[0].Geometry.Location.Longitude, data.Results[0].FormattedAddress, nil
}
