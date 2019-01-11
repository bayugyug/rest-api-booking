package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bayugyug/rest-api-booking/models"
)

const (
	GoogleApiUri   = "https://maps.googleapis.com/maps/api/geocode/json?key={key}&address="
	GoogleApiKey   = "AIzaSyCCLNYBuaKzGbyXHdmx_tgf4648wu_T794"
	MozillaLinuxUA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36"
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
func (g *GoogleMapGeoCode) GetCoordinates(address string) (bool, *models.Location, error) {

	//fmt
	if len(g.Key) <= 0 {
		g.Key = GoogleApiKey
	}
	gpoint := strings.Replace(GoogleApiUri+url.QueryEscape(address), "{key}", g.Key, -1)
	log.Println("GMAP_API: ", gpoint)
	//get from remote
	body, code, err := HttpGet(gpoint, map[string]string{"User-Agent": MozillaLinuxUA})
	if err != nil {
		log.Println("GMAP_API: failed", err)
		return false, nil, err
	}
	if code != http.StatusOK {
		log.Println("GMAP_API: failed", code)
		return false, nil, err
	}
	var data GoogleApiResponse
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Println("GMAP_API: failed", code)
		return false, nil, err
	}
	if len(data.Results) <= 0 || data.Status != "OK" {
		log.Println("GMAP_API: failed", data.Status)
		return false, nil, err
	}
	//good
	return true, &models.Location{
		Address:   data.Results[0].FormattedAddress,
		Latitude:  data.Results[0].Geometry.Location.Latitude,
		Longitude: data.Results[0].Geometry.Location.Longitude,
	}, nil
}
