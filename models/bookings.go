package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type Booking struct {
	ID             int     `json:"id"`
	MobileCustomer int     `json:"mobile_customer"`
	MobileDriver   int     `json:"mobile_driver"`
	Src            string  `json:"src"`
	SrcLatitude    float64 `json:"src_latitude"`
	SrcLongitude   float64 `json:"src_longitude"`
	Dst            string  `json:"dst"`
	DstLatitude    float64 `json:"dst_latitude"`
	DstLongitude   float64 `json:"dst_longitude"`
	Status         string  `json:"status"`
	PickupTime     string  `json:"pickup_time"`
	Dropofftime    string  `json:"dropoff_time"`
	Remarks        string  `json:"remarks"`
	RemarksBy      string  `json:"remarks_by"`
	Created        string  `json:"created"`
	Modified       string  `json:"modified"`
}

func (u *Booking) Exists(ctx context.Context, db *sql.DB, id int) int {
	r := `SELECT 1
		FROM  bookings WHERE id = ?`

	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -1
	}
	defer stmt.Close()
	var uid int
	err = stmt.QueryRowContext(ctx, id).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		log.Println("SQL_ERR::NO_ROWS", err)
		return -2
	case err != nil:
		log.Println("SQL_ERR::", err)
		return -3
	}
	//sounds good ;-)
	return uid
}

func (u *Booking) GetBooking(ctx context.Context, db *sql.DB, id int) (*Booking, error) {
	//fmt
	r := `SELECT 
			ifnull(id,''), 
			ifnull(mobile_customer,''), 
			ifnull(mobile_driver,''), 
			ifnull(src,''), 
			ifnull(src_latitude,0.0), 
			ifnull(src_longitude,0.0), 
			ifnull(dst,''), 
			ifnull(dst_latitude,0.0), 
			ifnull(dst_longitude,0.0), 
			ifnull(status,''), 
			ifnull(remarks,''), 
			ifnull(remarks_by,''), 
			ifnull(pickup_time,''), 
			ifnull(dropoff_time,''), 
			ifnull(created_dt,''), 
			ifnull(modified_dt,'')
		FROM  bookings WHERE id = ?`
	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR", err)
		return nil, err
	}
	defer stmt.Close()
	var data Booking
	err = stmt.QueryRowContext(ctx, id).Scan(
		&data.ID,
		&data.MobileCustomer,
		&data.MobileDriver,
		&data.Src,
		&data.SrcLatitude,
		&data.SrcLongitude,
		&data.Dst,
		&data.DstLatitude,
		&data.DstLongitude,
		&data.Status,
		&data.Remarks,
		&data.RemarksBy,
		&data.PickupTime,
		&data.Dropofftime,
		&data.Created,
		&data.Modified,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Println("SQL_ERR::NO_ROWS")
		return nil, errors.New("Info not found")
	case err != nil:
		log.Println("SQL_ERR", err)
		return nil, err
	}
	//sounds good ;-)
	return &data, nil
}

func (u *Booking) CreateBooking(ctx context.Context, db *sql.DB, data *Booking) (bool, error) {
	//fmt
	r := `INSERT INTO bookings(
			mobile_driver, 
			mobile_customer, 
			src, 
			src_latitude, 
			src_longitude, 
			dst, 
			dst_latitude, 
			dst_longitude, 
			status, 
			created_dt
	      )
	      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, Now()) `
	//exec
	result, err := db.ExecContext(ctx, r,
		data.MobileCustomer,
		data.MobileDriver,
		data.Src,
		data.SrcLatitude,
		data.SrcLongitude,
		data.Dst,
		data.DstLatitude,
		data.DstLongitude,
		data.Status)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to create")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to create")
	}
	if rows != 1 {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to create")

	}
	//sounds good ;-)
	return true, nil
}

func (u *Booking) UpdateBooking(ctx context.Context, db *sql.DB, data *Booking) (bool, error) {
	//fmt
	r := `UPDATE bookings
		SET 
		src = ?, 
		src_latitude = ?, 
		src_longitude= ?, 
		dst = ?, 
		dst_latitude = ?, 
		dst_longitude= ?, 
		remarks = ?, 
		remarks_by = ?, 
		modified_dt = Now() 
	      WHERE  id = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Src,
		data.SrcLatitude,
		data.SrcLongitude,
		data.Dst,
		data.DstLatitude,
		data.DstLongitude,
		data.Remarks,
		data.RemarksBy,
		data.ID,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *Booking) DeleteBooking(ctx context.Context, db *sql.DB, id int) (bool, error) {
	//fmt
	r := `UPDATE bookings
		SET 
		status      = 'deleted', 
		modified_dt = Now() 
	      WHERE  id = ?`
	//exec
	result, err := db.ExecContext(ctx, r, id)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *Booking) UpdateBookingStatus(ctx context.Context, db *sql.DB, data *Booking) (bool, error) {
	//fmt
	r := `UPDATE bookings
		SET 
		status      = ?,
		remarks     = ?,
		remarks_by  = ?,
		modified_dt = Now() 
	      WHERE  id = ?`
	//exec
	result, err := db.ExecContext(ctx, r, data.Status, data.Remarks, data.RemarksBy, data.ID)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *Booking) UpdateBookingPickupTime(ctx context.Context, db *sql.DB, data *Booking) (bool, error) {
	//fmt
	r := `UPDATE bookings
		SET 
		pickup_time = ?,
		status      = ?,
		modified_dt = Now() 
	      WHERE  id = ?`
	//exec
	result, err := db.ExecContext(ctx, r, data.PickupTime, data.Status, data.ID)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *Booking) UpdateBookingDropoffTime(ctx context.Context, db *sql.DB, data *Booking) (bool, error) {
	//fmt
	r := `UPDATE bookings
		SET 
		dropoff_time= ?,
		status      = ?,
		modified_dt = Now() 
	      WHERE  id = ?`
	//exec
	result, err := db.ExecContext(ctx, r, data.Dropofftime, data.Status, data.ID)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}
