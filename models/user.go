package models

const (
	UserStatusPending      = "pending"
	UserStatusActive       = "active"
	UserStatusDeleted      = "deleted"
	VehicleStatusOpen      = "open"
	VehicleStatusBooked    = "booked"
	VehicleStatusCanceled  = "canceled"
	VehicleStatusTripStart = "trip-start"
	VehicleStatusTripEnd   = "trip-end"
	VehicleStatusGasUp     = "gas-up"
	VehicleStatusPanic     = "panic"
)

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
	Otp           string  `json:"otp"`
	OtpExpiry     string  `json:"otp_expiry"`
	Logged        int     `json:"logged"`
}
