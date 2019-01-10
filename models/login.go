package models

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	UserTypeCustomer = "customer"
	UserTypeDriver   = "driver"
)

type UserLogin struct {
	Mobile string `json:"mobile"`
	Pass   string `json:"pass"`
	Hash   string `json:"hash"`
	Type   string `json:"type"`
}

func (u *UserLogin) Bind(r *http.Request) error {
	//sanity check
	if u == nil {
		return errors.New("Missing required parameter")
	}
	// just a post-process after a decode..
	u.Hash = fmt.Sprintf("%x", md5.Sum([]byte(u.Pass)))
	if u.Type != UserTypeDriver && u.Type != UserTypeCustomer {
		return errors.New("Missing required parameter: " + u.Type)
	}
	return nil
}

func (u *UserLogin) GetUserInfo(ctx context.Context, db *sql.DB, who, mobile string) (*User, error) {
	//check whith type
	var r string

	switch who {
	case UserTypeDriver:
		r = `SELECT 
			ifnull(id,''), 
			ifnull(mobile,''), 
			ifnull(firstname,''), 
			ifnull(lastname,''), 
			ifnull(status,''), 
			ifnull(vehiclestatus,''), 
			ifnull(pass,''), 
			ifnull(latitude,0.0), 
			ifnull(longitude,0.0), 
			ifnull(created_dt,''), 
			ifnull(modified_dt,'')
		FROM  drivers WHERE mobile = ?`
	case UserTypeCustomer:
		r = `SELECT 
			ifnull(id,''), 
			ifnull(mobile,''), 
			ifnull(firstname,''), 
			ifnull(lastname,''), 
			ifnull(status,''), 
			"" as vehiclestatus, 
			ifnull(pass,''), 
			ifnull(latitude,0.0), 
			ifnull(longitude,0.0), 
			ifnull(created_dt,''), 
			ifnull(modified_dt,'')
		FROM  customers WHERE mobile = ?`
	}
	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("GetLoginInfo", err)
		return nil, err
	}
	defer stmt.Close()
	var data User
	err = stmt.QueryRowContext(ctx, mobile).Scan(
		&data.ID,
		&data.Mobile,
		&data.Firstname,
		&data.Lastname,
		&data.Status,
		&data.VehicleStatus,
		&data.Pass,
		&data.Latitude,
		&data.Longitude,
		&data.Created,
		&data.Modified,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Println("GetLoginInfo NOT_FOUND", mobile)
		return nil, errors.New("Info not found")
	case err != nil:
		log.Println("GetLoginInfo", err)
		return nil, err
	}
	data.Type = UserTypeDriver
	//sounds good ;-)
	return &data, nil
}

func (u *UserLogin) SetUserLogStatus(ctx context.Context, db *sql.DB, who, mobile string, flag int) bool {
	//check whith type
	var r string

	switch who {
	case UserTypeDriver:
		r = `UPDATE drivers SET logged=?, modified_dt=Now() WHERE  mobile = ?`
	case UserTypeCustomer:
		r = `UPDATE customers SET logged=?, modified_dt=Now() WHERE  mobile = ?`
	}

	result, err := db.ExecContext(ctx, r, flag, mobile)
	if err != nil {
		log.Println("SetUserLogStatus", mobile, err)
		return false
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SetUserLogStatus", mobile, err)
		return false
	}
	//sounds good ;-)
	return true
}
