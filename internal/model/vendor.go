package model

import "time"

type Vendor struct {
	ID            int       `json:"id"`
	VendorName    string    `json:"vendor_name"`
	Email         string    `json:"email"`
	ContactInfo   string    `json:"contact_info"`
	WalletAddress string    `json:"wallet_address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
