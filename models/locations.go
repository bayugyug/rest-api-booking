package models

import (
	"errors"
	"net/http"
)

type Location struct {
	Address   string  `json:"address"`
	Mobile    string  `json:"mobile"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (u *Location) Bind(r *http.Request) error {
	//sanity check
	if u == nil {
		return errors.New("Missing required parameter")
	}
	// just a post-process after a decode..
	return nil
}
