package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID             int             `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	VendorID       int             `json:"vendor_id"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	Status         string          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Picture        string          `json:"picture"`
	MaxAttendees   int             `json:"maxattendees"`
	Location       string          `json:"location"`
	Attendees      int             `json:"attendees"`
	Requirements   json.RawMessage `json:"requirements,omitempty"`
	Agenda         json.RawMessage `json:"agenda,omitempty"`
	Token          string          `json:"token,omitempty"`
	Organizer      *string         `json:"organizer,omitempty"`
	Whitelisted    int             `json:"whitelisted,omitempty"`
	Minted         int             `json:"minted,omitempty"`
	UrlCertificate *string         `json:"url_certificate,omitempty"`
}
