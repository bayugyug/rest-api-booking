package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
)

type Customer User

func (u *Customer) Bind(r *http.Request) error {
	//sanity check
	if u == nil {
		return errors.New("Missing required parameter")
	}
	return nil
}

func (u *Customer) Exists(ctx context.Context, db *sql.DB, mobile string) int {
	r := `SELECT id
		FROM  customers WHERE mobile = ?`

	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -1
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRowContext(ctx, mobile).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		log.Println("SQL_ERR::NO_ROWS", err)
		return -2
	case err != nil:
		log.Println("SQL_ERR::", err)
		return -3
	}
	//sounds good ;-)
	return id
}

func (u *Customer) GetCustomer(ctx context.Context, db *sql.DB, mobile string) (*Customer, error) {
	//fmt
	r := `SELECT 
			ifnull(id,''), 
			ifnull(mobile,''), 
			ifnull(firstname,''), 
			ifnull(lastname,''), 
			ifnull(status,''), 
			ifnull(pass,''), 
			ifnull(otp,''), 
			ifnull(otp_expiry,''), 
			ifnull(latitude,0.0), 
			ifnull(longitude,0.0), 
			ifnull(created_dt,''), 
			ifnull(modified_dt,'')
		FROM  customers WHERE mobile = ?`
	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR", err)
		return nil, err
	}
	defer stmt.Close()
	var data Customer
	err = stmt.QueryRowContext(ctx, mobile).Scan(
		&data.ID,
		&data.Mobile,
		&data.Firstname,
		&data.Lastname,
		&data.Status,
		&data.Pass,
		&data.Otp,
		&data.OtpExpiry,
		&data.Latitude,
		&data.Longitude,
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
	data.Type = UserTypeCustomer
	//sounds good ;-)
	return &data, nil
}

func (u *Customer) CreateCustomer(ctx context.Context, db *sql.DB, data *Customer) (bool, error) {
	//fmt
	r := `INSERT INTO customers (mobile, firstname, lastname, pass, latitude, longitude, status, otp, otp_expiry, created_dt)
	      VALUES (?, ?, ?, ?, ?, ?, 'pending', Now()) 
	      ON DUPLICATE KEY UPDATE
	        firstname =?, 
		lastname  =?, 
		latitude  =?, 
		longitude =?,
		modified_dt = Now() `
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Mobile,
		data.Firstname,
		data.Lastname,
		data.Pass,
		data.Latitude,
		data.Longitude,
		data.Otp,
		data.OtpExpiry,
		data.Firstname, //update starts here
		data.Lastname,
		data.Latitude,
		data.Longitude,
	)
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

func (u *Customer) UpdateCustomer(ctx context.Context, db *sql.DB, data *Customer) (bool, error) {
	//fmt
	r := `UPDATE customers 
		SET 
		firstname =?, 
		lastname  =?, 
		latitude  =?, 
		longitude =?,
		modified_dt = Now() 
	      WHERE  mobile = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Firstname,
		data.Lastname,
		data.Latitude,
		data.Longitude,
		data.Mobile,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	if rows != 1 {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")

	}
	//sounds good ;-)
	return true, nil
}

func (u *Customer) DeleteCustomer(ctx context.Context, db *sql.DB, mobile string) (bool, error) {
	//fmt
	r := `UPDATE customers 
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

func (u *Customer) UpdateCustomerCoords(ctx context.Context, db *sql.DB, data *Customer) (bool, error) {
	//fmt
	r := `UPDATE customers 
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

func (u *Customer) UpdateCustomerPass(ctx context.Context, db *sql.DB, data *Customer) (bool, error) {
	//fmt
	r := `UPDATE customers 
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

func (u *Customer) UpdateCustomerStatus(ctx context.Context, db *sql.DB, status, mobile string) (bool, error) {
	//fmt
	r := `UPDATE customers 
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

func (u *Customer) UpdateCustomerOtp(ctx context.Context, db *sql.DB, mobile, otp, otpexp string) (bool, error) {
	//fmt
	r := `UPDATE customers 
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
