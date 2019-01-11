package models

type Otp struct {
	Mobile string `json:"mobile"`
	Type   string `json:"type"`
	Otp    string `json:"otp"`
}
