package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
)

const (
	NearestDistance        = 50 //km distance nearest
	VehicleStatusOpen      = "open"
	VehicleStatusBooked    = "booked"
	VehicleStatusCanceled  = "canceled"
	VehicleStatusTripStart = "trip-start"
	VehicleStatusTripEnd   = "trip-end"
	VehicleStatusCompleted = "completed"
	VehicleStatusGasUp     = "gas-up"
	VehicleStatusPanic     = "panic"
)

type Driver User

type DriveListInfo struct {
	ID            int64   `json:"id"`
	Mobile        string  `json:"mobile"`
	Firstname     string  `json:"firstname"`
	Lastname      string  `json:"lastname"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Status        string  `json:"status"`
	VehicleStatus string  `json:"vehicle_status"`
	Distance      float64 `json:"distance"`
}

func (u *Driver) Bind(r *http.Request) error {
	//sanity check
	if u == nil {
		return errors.New("Missing required parameter")
	}
	u.Type = UserTypeDriver
	return nil
}

func (u *Driver) Exists(ctx context.Context, db *sql.DB, mobile string) int {
	r := `SELECT count(1)
		FROM  drivers WHERE mobile = ?`

	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -1
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRowContext(ctx, mobile).Scan(&id)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -2
	}
	//sounds good ;-)
	return id
}

func (u *Driver) GetDriver(ctx context.Context, db *sql.DB, mobile string) (*Driver, error) {
	//fmt
	r := `SELECT 
			ifnull(id,''), 
			ifnull(mobile,''), 
			ifnull(firstname,''), 
			ifnull(lastname,''), 
			ifnull(status,''), 
			ifnull(vehiclestatus,''), 
			ifnull(pass,''), 
			ifnull(otp,''), 
			ifnull(otp_expiry,''), 
			ifnull(latitude,0.0), 
			ifnull(longitude,0.0), 
			ifnull(created_dt,''), 
			ifnull(modified_dt,''),
			ifnull((otp_expiry <now()),0)
		FROM  drivers WHERE mobile = ?`
	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return nil, err
	}
	defer stmt.Close()
	var data Driver
	err = stmt.QueryRowContext(ctx, mobile).Scan(
		&data.ID,
		&data.Mobile,
		&data.Firstname,
		&data.Lastname,
		&data.Status,
		&data.VehicleStatus,
		&data.Pass,
		&data.Otp,
		&data.OtpExpiry,
		&data.Latitude,
		&data.Longitude,
		&data.Created,
		&data.Modified,
		&data.OtpExpired,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Println("SQL_ERR::NO_ROWS")
		return nil, errors.New("Info not found")
	case err != nil:
		log.Println("SQL_ERR", err)
		return nil, err
	}
	data.Type = UserTypeDriver
	//sounds good ;-)
	return &data, nil
}

func (u *Driver) CreateDriver(ctx context.Context, db *sql.DB, data *Driver) (bool, error) {
	//fmt
	r := `INSERT INTO drivers (
		mobile, 
		firstname, 
		lastname, 
		pass, 
		latitude, 
		longitude, 
		status,
		otp, 
		otp_expiry, 
		created_dt)
	      VALUES (?, ?, ?, ?, ?, ?, 'pending', ?, ?, ?) 
	      ON DUPLICATE KEY UPDATE
	        firstname =?, 
		lastname  =?, 
		latitude  =?, 
		longitude =?,
		modified_dt = Now() `
	//exec
	result, err := db.Exec(r,
		data.Mobile,
		data.Firstname,
		data.Lastname,
		data.Pass,
		data.Latitude,
		data.Longitude,
		data.Otp,
		data.OtpExpiry,
		data.Created,
		data.Firstname, //update starts here
		data.Lastname,
		data.Latitude,
		data.Longitude,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to create")
	}
	id, err := result.LastInsertId()
	log.Println("LAST_INSERT_ID", id)
	if err != nil || id < 1 {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to create")
	}
	data.ID = id
	//sounds good ;-)
	return true, nil
}

func (u *Driver) UpdateDriver(ctx context.Context, db *sql.DB, data *Driver) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		firstname =?, 
		lastname  =?, 
		latitude  =?, 
		longitude =?,
		modified_dt = ?
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Firstname,
		data.Lastname,
		data.Latitude,
		data.Longitude,
		data.Modified,
		data.Mobile,
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

func (u *Driver) DeleteDriver(ctx context.Context, db *sql.DB, mobile string) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		status      = 'deleted', 
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r, mobile)
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

func (u *Driver) UpdateDriverCoords(ctx context.Context, db *sql.DB, data *Driver) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		latitude  =?, 
		longitude =?,
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Latitude,
		data.Longitude,
		data.Mobile,
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

func (u *Driver) UpdateDriverPass(ctx context.Context, db *sql.DB, data *Driver) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		pass = ?, 
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Pass,
		data.Mobile,
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

func (u *Driver) UpdateDriverStatus(ctx context.Context, db *sql.DB, status, mobile string) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		status = ?, 
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		status,
		mobile,
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

func (u *Driver) UpdateDriverOtp(ctx context.Context, db *sql.DB, mobile, otp, otpexp string) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		otp = ?, 
		otp_expiry = ?, 
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		otp,
		otpexp,
		mobile,
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

func (u *Driver) UpdateDriverOtpExpiry(ctx context.Context, db *sql.DB, data *Driver) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		otp_expiry= date_add(now(), interval -1 day), 
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Mobile,
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

func (u *Driver) UpdateDriverVehicleStatus(ctx context.Context, db *sql.DB, status, mobile string) (bool, error) {
	//fmt
	r := `UPDATE drivers 
		SET 
		vehiclestatus = ?, 
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		status,
		mobile,
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

func (u *Driver) GetDriversNearestLocation(ctx context.Context, db *sql.DB, lat, lon float64, distance int) ([]DriveListInfo, error) {

	var all []DriveListInfo

	//fmt earth radius 6371 KM
	r := `SELECT 
			ifnull(id,''), 
			ifnull(mobile,''), 
			ifnull(firstname,''), 
			ifnull(lastname,''), 
			ifnull(status,''), 
			ifnull(vehiclestatus,''), 
			ifnull(latitude,0.0), 
			ifnull(longitude,0.0), 
		   (
			6371 * acos (
			  cos ( radians( ? ) )
			  * cos( radians( ifnull(latitude,0.0 ) ) )
			  * cos( radians( ifnull(longitude,0.0) ) - radians( ? ) )
			  + sin ( radians( ? ) )
			  * sin( radians( ifnull(latitude,0.0) ) )
			)
		    ) AS distance
		FROM  drivers 
		WHERE status = 'active'
		HAVING distance < ?
		ORDER BY distance ASC
		LIMIT 10
		`
	rows, err := db.Query(r, lat, lon, lat, distance)
	if err != nil {
		log.Println("SQL_ERR", err)
		return all, err
	}
	defer rows.Close()

	for rows.Next() {
		var data DriveListInfo
		if err := rows.Scan(
			&data.ID,
			&data.Mobile,
			&data.Firstname,
			&data.Lastname,
			&data.Status,
			&data.VehicleStatus,
			&data.Latitude,
			&data.Longitude,
			&data.Distance,
		); err != nil {
			log.Println("SQL_ERR::", err)
			continue
		}
		//save
		all = append(all, data)
	}
	//sounds good ;-)
	return all, nil
}
