package models

type User struct {
	ID            int     `json:"id"`
	Mobile        string  `json:"mobile"`
	Firstname     string  `json:"firstname"`
	Lastname      string  `json:"lastname"`
	Status        string  `json:"status"`
	Pass          string  `json:"-"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Created       string  `json:"created"`
	Modified      string  `json:"modified"`
	Type          string  `json:"type"`
	VehicleStatus string  `json:"vehiclestatus"`
	Logged        int     `json:"logged"`
}
