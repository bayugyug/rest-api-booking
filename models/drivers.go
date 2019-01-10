package models

import (
	"context"
	"database/sql"
	"errors"
)

type Driver struct{}

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
