package models

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
)

type Driver User

func (u *Driver) Bind(r *http.Request) error {
	//sanity check
	if u == nil {
		return errors.New("Missing required parameter")
	}
	// just a post-process after a decode..
	if u.Type != UserTypeDriver {
		return errors.New("Missing required parameter: " + u.Type)
	}
	return nil
}

func (u *Driver) GetDriver(ctx context.Context, db *sql.DB, mobile string) (*User, error) {
	return nil, errors.New("Not implemented")
}

func (u *Driver) CreateDriver(ctx context.Context, db *sql.DB, data *User) (int, error) {
	return 0, errors.New("Not implemented")
}

func (u *Driver) UpdateDriver(ctx context.Context, db *sql.DB, data *User) (bool, error) {
	return true, errors.New("Not implemented")
}

func (u *Driver) DeleteDriver(ctx context.Context, db *sql.DB, mobile string) (bool, error) {
	return true, errors.New("Not implemented")
}
