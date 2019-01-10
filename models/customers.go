package models

import (
	"context"
	"database/sql"
	"errors"
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

func (u *Customer) CreateCheck(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (u *Customer) GetCustomer(ctx context.Context, db *sql.DB, mobile string) (*Customer, error) {
	return nil, errors.New("Not implemented")
}

func (u *Customer) CreateCustomer(ctx context.Context, db *sql.DB, data *Customer) (int, error) {
	return 0, errors.New("Not implemented")
}

func (u *Customer) UpdateCustomer(ctx context.Context, db *sql.DB, data *Customer) (bool, error) {
	return true, errors.New("Not implemented")
}

func (u *Customer) DeleteCustomer(ctx context.Context, db *sql.DB, mobile string) (bool, error) {
	return true, errors.New("Not implemented")
}
