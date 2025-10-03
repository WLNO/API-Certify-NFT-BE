package model

import "time"

type CertificateWithEvent struct {
	ID                  int                    `json:"id"`
	EventID             int                    `json:"event_id"`
	UserID              int                    `json:"user_id"`
	CertificateData     map[string]interface{} `json:"certificate_data"`
	MintStatus          *string                `json:"mint_status"`
	MintTransactionHash *string                `json:"mint_transaction_hash"`
	CreatedAt           string                 `json:"created_at"`
	UpdatedAt           string                 `json:"updated_at"`
	UrlMetadata         *string                `json:"url_metadata"`
	UrlCertificate      *string                `json:"url_certificate"`
	CertificateType     *string                `json:"certificate_type"`
	// Event summary
	EventTitle       string    `json:"event_title,omitempty"`
	EventDescription string    `json:"event_description,omitempty"`
	EventStartDate   time.Time `json:"event_start_date,omitempty"`
	EventLocation    string    `json:"event_location,omitempty"`
	EventPicture     string    `json:"event_picture,omitempty"`
}
