package models

import (
	"context"
	"database/sql"
	"errors"
)

type Booking struct {
	ID           int     `json:"id"`
	CustomerID   int     `json:"customer_id"`
	DriverID     int     `json:"driver_id"`
	Src          string  `json:"src"`
	SrcLatitude  float64 `json:"src_latitude"`
	SrcLongitude float64 `json:"src_longitude"`
	Dst          string  `json:"dst"`
	DstLatitude  float64 `json:"dst_latitude"`
	DstLongitude float64 `json:"dst_longitude"`
	Status       string  `json:"status"`
	PickupTime   string  `json:"pickup_time"`
	Dropofftime  string  `json:"dropoff_time"`
	Remarks      string  `json:"remarks"`
	RemarksBy    string  `json:"remarks_by"`
	Created      string  `json:"created"`
	Modified     string  `json:"modified"`
}

func (u *Booking) GetBooking(ctx context.Context, db *sql.DB, id int) (*Booking, error) {
	return nil, errors.New("Not implemented")
}

func (u *Booking) CreateBooking(ctx context.Context, db *sql.DB, data *Booking) (int, error) {
	return 0, errors.New("Not implemented")
}

func (u *Booking) UpdateBooking(ctx context.Context, db *sql.DB, data *Booking) (bool, error) {
	return true, errors.New("Not implemented")
}

func (u *Booking) DeleteBooking(ctx context.Context, db *sql.DB, id int) (bool, error) {
	return true, errors.New("Not implemented")
}
