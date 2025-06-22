package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Event struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	VendorID     int       `json:"vendor_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Picture      string    `json:"picture"`
	MaxAttendees int       `json:"maxattendees"`
	Location     string    `json:"location"`
	Attendees    int       `json:"attendees"`
}

type User struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	WalletAddress string    `json:"wallet_address"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Vendor struct {
	ID            int       `json:"id"`
	VendorName    string    `json:"vendor_name"`
	Email         string    `json:"email"`
	ContactInfo   string    `json:"contact_info"`
	WalletAddress string    `json:"wallet_address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var db *sql.DB

func registerUserHandler(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if u.Name == "" || u.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and wallet_address are required"})
	}

	registered, err := isWalletRegistered(u.WalletAddress)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server error"})
	}
	if registered {
		return c.JSON(http.StatusConflict, map[string]string{"error": "wallet address already registered with another account"})
	}

	query := `INSERT INTO users (email, wallet_address, name) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err = db.QueryRow(query, u.Email, u.WalletAddress, u.Name).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, u)
}

func registerVendorHandler(c echo.Context) error {
	var v Vendor
	if err := c.Bind(&v); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if v.VendorName == "" || v.Email == "" || v.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "vendor_name, email, and wallet_address are required"})
	}

	registered, err := isWalletRegistered(v.WalletAddress)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server error"})
	}
	if registered {
		return c.JSON(http.StatusConflict, map[string]string{"error": "wallet address already registered with another account"})
	}

	query := `INSERT INTO vendors (vendor_name, email, contact_info, wallet_address) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err = db.QueryRow(query, v.VendorName, v.Email, v.ContactInfo, v.WalletAddress).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, v)
}

func loginHandler(c echo.Context) error {
	var body struct {
		WalletAddress string `json:"wallet_address"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if body.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
	}

	var role string

	// Cek apakah wallet_address ada di tabel users
	var userExists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE wallet_address = $1)`, body.WalletAddress).Scan(&userExists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if userExists {
		role = "users"
	} else {
		// Jika tidak di users, cek di vendors
		var vendorExists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM vendors WHERE wallet_address = $1)`, body.WalletAddress).Scan(&vendorExists)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if vendorExists {
			role = "vendors"
		}
	}

	if role != "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"isNewUser":     false,
			"wallet_address": body.WalletAddress,
			"role":           role,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"isNewUser": true,
	})
}

func isWalletRegistered(wallet string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE wallet_address = $1)`, wallet).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM vendors WHERE wallet_address = $1)`, wallet).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func main() {
	var err error

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		panic("DATABASE_DSN is not set in environment")
	}

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("Error opening database: %v", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	}))
	e.GET("/api/events/all", getEventsHandler)
	e.POST("/api/events/create", createEventHandler)
	e.POST("/api/users/register", registerUserHandler)
	e.POST("/api/vendors/register", registerVendorHandler)
	e.POST("/api/auth/login", loginHandler)
	e.GET("/api/users/:id/events", getEventsByUserIdHandler)
	e.GET("/api/users/:id/certificate", getCertificatesByUserIdHandler)

	e.Logger.Fatal(e.Start(":4002"))
}

func getEventsHandler(c echo.Context) error {
	rows, err := db.Query(`
    SELECT 
      e.id, e.title, e.description, e.vendor_id, e.start_date, e.end_date, 
      e.status, e.created_at, e.updated_at, e.picture, e.maxattendees, e.location,
      COUNT(a.id) as attendees
    FROM events e
    LEFT JOIN attendance a ON a.event_id = e.id AND a.attendance_status = 'present'
    GROUP BY e.id
    ORDER BY e.start_date ASC
  `)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VendorID, &e.StartDate, &e.EndDate, &e.Status, &e.CreatedAt, &e.UpdatedAt, &e.Picture, &e.MaxAttendees, &e.Location, &e.Attendees)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		events = append(events, e)
	}

	return c.JSON(http.StatusOK, events)
}

func getEventsByUserIdHandler(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	rows, err := db.Query(`
		SELECT 
			e.id, e.title, e.description, e.vendor_id, e.start_date, e.end_date, 
			e.status, e.created_at, e.updated_at, e.picture, e.maxattendees, e.location,
			COUNT(a2.id) as attendees
		FROM events e
		INNER JOIN attendance a ON a.event_id = e.id
		LEFT JOIN attendance a2 ON a2.event_id = e.id AND a2.attendance_status = 'present'
		WHERE a.user_id = $1
		GROUP BY e.id
		ORDER BY e.start_date ASC
	`, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VendorID, &e.StartDate, &e.EndDate, &e.Status, &e.CreatedAt, &e.UpdatedAt, &e.Picture, &e.MaxAttendees, &e.Location, &e.Attendees)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		events = append(events, e)
	}

	return c.JSON(http.StatusOK, events)
}

// CertificateWithEvent combines certificate data with event fields for richer output
type CertificateWithEvent struct {
	ID                  int       `json:"id"`
	EventID             int       `json:"event_id"`
	UserID              int       `json:"user_id"`
	CertificateData     string    `json:"certificate_data"`
	MintStatus          string    `json:"mint_status"`
	MintTransactionHash string    `json:"mint_transaction_hash"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	EventTitle          string    `json:"event_title"`
	EventDescription    string    `json:"event_description"`
	EventStartDate      time.Time `json:"event_start_date"`
	EventLocation       string    `json:"event_location"`
	EventPicture        string    `json:"event_picture"`
}

func getCertificatesByUserIdHandler(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	rows, err := db.Query(`
		SELECT 
			c.id, c.event_id, c.user_id, c.certificate_data, c.mint_status, c.mint_transaction_hash, c.created_at, c.updated_at,
			e.title, e.description, e.start_date, e.location, e.picture
		FROM certificates c
		JOIN events e ON c.event_id = e.id
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC
	`, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var certificates []CertificateWithEvent
	for rows.Next() {
		var cert CertificateWithEvent
		var pictureRaw string
		err := rows.Scan(
			&cert.ID, &cert.EventID, &cert.UserID, &cert.CertificateData, &cert.MintStatus, &cert.MintTransactionHash, &cert.CreatedAt, &cert.UpdatedAt,
			&cert.EventTitle, &cert.EventDescription, &cert.EventStartDate, &cert.EventLocation, &pictureRaw,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		// Convert picture field to URL path
		if pictureRaw != "" {
			cert.EventPicture = "https://api.gpadaka.com/" + pictureRaw
		} else {
			cert.EventPicture = ""
		}
		certificates = append(certificates, cert)
	}

	return c.JSON(http.StatusOK, certificates)
}

func createEventHandler(c echo.Context) error {
	title := c.FormValue("title")
	description := c.FormValue("description")
	vendorIDStr := c.FormValue("vendor_id")
	startDateStr := c.FormValue("start_date")
	endDateStr := c.FormValue("end_date")
	status := c.FormValue("status")
	maxAttendeesStr := c.FormValue("maxattendees")

	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
	}
	if vendorIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "vendor_id is required"})
	}
	if startDateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_date is required"})
	}
	if endDateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "end_date is required"})
	}
	if status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "status is required"})
	}
	if maxAttendeesStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "maxattendees is required"})
	}

	file, err := c.FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "picture file is required"})
	}
	filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to open uploaded file"})
	}
	defer src.Close()

	dst, err := os.Create(filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create file on server"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save picture"})
	}

	vendorID, err := strconv.Atoi(vendorIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "vendor_id must be a number"})
	}	
	maxAttendees, err := strconv.Atoi(maxAttendeesStr)
	if err != nil || maxAttendees <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "maxattendees must be a positive number"})
	}
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid start_date format"})
	}
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid end_date format"})
	}
	if startDate.After(endDate) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_date must be before end_date"})
	}

	query := `
		INSERT INTO events (title, description, vendor_id, start_date, end_date, status, picture, maxattendees)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	var event Event
	event.Title = title
	event.Description = description
	event.VendorID = vendorID
	event.StartDate = startDate
	event.EndDate = endDate
	event.Status = status
	event.Picture = filename
	event.MaxAttendees = maxAttendees

	err = db.QueryRow(
		query,
		event.Title,
		event.Description,
		event.VendorID,
		event.StartDate,
		event.EndDate,
		event.Status,
		event.Picture,
		event.MaxAttendees,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, event)
}



