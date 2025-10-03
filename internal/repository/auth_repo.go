package repository

import "database/sql"

type AuthRepo interface {
	IsUserWallet(wallet string) (bool, error)
	IsVendorWallet(wallet string) (bool, error)
	CreateUser(name, email, wallet string) (int, string, string, error)
	CreateVendor(name, email, contact, wallet string) (int, error)
}

type authRepo struct{ db *sql.DB }

func NewAuthRepo(db *sql.DB) AuthRepo { return &authRepo{db: db} }

func (r *authRepo) IsUserWallet(w string) (bool, error) {
	var e bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE wallet_address=$1)`, w).Scan(&e)
	return e, err
}

func (r *authRepo) IsVendorWallet(w string) (bool, error) {
	var e bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM vendors WHERE wallet_address=$1)`, w).Scan(&e)
	return e, err
}

func (r *authRepo) CreateUser(name, email, wallet string) (int, string, string, error) {
	var id int
	var created, updated string
	err := r.db.QueryRow(
		`INSERT INTO users (email, wallet_address, name) VALUES ($1,$2,$3) RETURNING id, created_at, updated_at`,
		email, wallet, name,
	).Scan(&id, &created, &updated)
	return id, created, updated, err
}

func (r *authRepo) CreateVendor(name, email, contact, wallet string) (int, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO vendors (vendor_name, email, contact_info, wallet_address) VALUES ($1,$2,$3,$4) RETURNING id`,
		name, email, contact, wallet,
	).Scan(&id)
	return id, err
}
