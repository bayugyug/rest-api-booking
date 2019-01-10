package config

import (
	"encoding/json"
	"log"

	"github.com/bayugyug/rest-api-booking/driver"
)

//AppConfig optional parameter structure
type AppConfig struct {
	Driver       driver.DbConnectorConfig `json:"driver"`
	GoogleApiKey string                   `json:"google_api_key"`
}

//api global handler
var ApiConfig *AppConfig

//NewAppConfig new AppConfig
func FormatAppConfig(s string) *AppConfig {
	var cfg AppConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		log.Println("FormatAppConfig", err)
		return nil
	}
	return &cfg
}
