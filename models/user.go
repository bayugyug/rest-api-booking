package models

type User struct {
	ID            int    `json:"id"`
	Mobile        string `json:"mobile"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Status        string `json:"status"`
	Pass          string `json:"-"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	Created       string `json:"created"`
	Modified      string `json:"modified"`
	Type          string `json:"type"`
	VehicleStatus string `json:"vehiclestatus"`
}
