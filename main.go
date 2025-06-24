package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Event struct {
	ID           int             `json:"id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	VendorID     int             `json:"vendor_id"`
	StartDate    time.Time       `json:"start_date"`
	EndDate      time.Time       `json:"end_date"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Picture      string          `json:"picture"`
	MaxAttendees int             `json:"maxattendees"`
	Location     string          `json:"location"`
	Attendees    int             `json:"attendees"`
	Requirements json.RawMessage `json:"requirements,omitempty"`
	Agenda       json.RawMessage `json:"agenda,omitempty"`
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
	e.GET("/api/users/:walletAddress/events", getEventsByWalletAddressHandler)
	e.GET("/api/users/:walletAddress/certificate", getCertificatesByWalletAddressHandler)
	e.GET("/api/vendors/:walletAddress/events", getEventsByVendorWalletAddressHandler)
	e.GET("/api/events/:id", getEventDetailHandler)
	e.POST("/api/users/whitelist", createWhitelistHandler)
	e.POST("/api/users/whitelist/cancel", cancelWhitelistHandler)

	e.Logger.Fatal(e.Start(":4002"))
}

// Handler to cancel whitelist entry
func cancelWhitelistHandler(c echo.Context) error {
	var payload struct {
		EventID       int    `json:"event_id"`
		WalletAddress string `json:"wallet_address"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}

	if payload.EventID == 0 || payload.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event_id and wallet_address are required"})
	}

	// Get user_id from wallet_address
	var userID int
	err := db.QueryRow(`SELECT id FROM users WHERE wallet_address = $1`, payload.WalletAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Delete the whitelist entry
	result, err := db.Exec(`DELETE FROM whitelist WHERE event_id = $1 AND user_id = $2`, payload.EventID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to cancel whitelist"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "whitelist entry not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Whitelist cancelled successfully"})
}

func getEventsHandler(c echo.Context) error {
	rows, err := db.Query(`
    SELECT 
      e.id, e.title, e.description, e.vendor_id, e.start_date, e.end_date, 
      e.status, e.created_at, e.updated_at, e.picture, e.maxattendees, e.location,
      v.vendor_name,
      (SELECT COUNT(*) FROM whitelist w WHERE w.event_id = e.id AND w.status = 'approved') as whitelisted
    FROM events e
    LEFT JOIN vendors v ON v.id = e.vendor_id
    GROUP BY e.id, v.vendor_name
    ORDER BY e.start_date ASC
  `)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	type EventWithOrganizer struct {
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
		Organizer    string    `json:"organizer"`
		Whitelisted  int       `json:"whitelisted"`
		Location     string    `json:"location"`
	}

	var events []EventWithOrganizer
	for rows.Next() {
		var e EventWithOrganizer
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.VendorID, &e.StartDate, &e.EndDate,
			&e.Status, &e.CreatedAt, &e.UpdatedAt, &e.Picture, &e.MaxAttendees, &e.Location,
			&e.Organizer, &e.Whitelisted,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		events = append(events, e)
	}

	return c.JSON(http.StatusOK, events)
}

func getEventsByWalletAddressHandler(c echo.Context) error {
	walletAddress := c.Param("walletAddress")
	if walletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "walletAddress is required"})
	}

	var userID int
	err := db.QueryRow(`SELECT id FROM users WHERE wallet_address = $1`, walletAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	query := `
		SELECT
			e.id, e.title, e.description, e.vendor_id, e.start_date, e.end_date,
			e.status, e.created_at, e.updated_at, e.picture, e.maxattendees, e.location,
			e.requirements, e.agenda,
			COALESCE(attend_count.attendees, 0) as attendees,
			CASE
				WHEN att.attendance_status = 'present' THEN 'present'
				ELSE wl.status
			END as user_status
		FROM events e
		LEFT JOIN attendance att ON e.id = att.event_id AND att.user_id = $1
		LEFT JOIN whitelist wl ON e.id = wl.event_id AND wl.user_id = $1
		LEFT JOIN (
			SELECT event_id, COUNT(*) as attendees
			FROM attendance
			WHERE attendance_status = 'present'
			GROUP BY event_id
		) AS attend_count ON e.id = attend_count.event_id
		WHERE att.user_id IS NOT NULL OR wl.user_id IS NOT NULL
		ORDER BY e.start_date ASC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	type EventResponse struct {
		ID           int             `json:"id"`
		Title        string          `json:"title"`
		Description  string          `json:"description"`
		VendorID     int             `json:"vendor_id"`
		StartDate    time.Time       `json:"start_date"`
		EndDate      time.Time       `json:"end_date"`
		Status       string          `json:"status"`
		CreatedAt    time.Time       `json:"created_at"`
		UpdatedAt    time.Time       `json:"updated_at"`
		Picture      string          `json:"picture"`
		MaxAttendees int             `json:"maxattendees"`
		Location     string          `json:"location"`
		Attendees    int             `json:"attendees"`
		UserStatus   string          `json:"user_status"`
		Requirements json.RawMessage `json:"requirements,omitempty"`
		Agenda       json.RawMessage `json:"agenda,omitempty"`
	}

	var events []EventResponse
	for rows.Next() {
		var e EventResponse
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.VendorID, &e.StartDate, &e.EndDate, &e.Status,
			&e.CreatedAt, &e.UpdatedAt, &e.Picture, &e.MaxAttendees, &e.Location,
			&e.Requirements, &e.Agenda, &e.Attendees,
			&e.UserStatus,
		)
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

// Handler to get events by vendor wallet address
func getEventsByVendorWalletAddressHandler(c echo.Context) error {
	walletAddress := c.Param("walletAddress")
	if walletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "walletAddress is required"})
	}

	var vendorID int
	err := db.QueryRow(`SELECT id FROM vendors WHERE wallet_address = $1`, walletAddress).Scan(&vendorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "vendor not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	query := `
		SELECT 
			e.id, e.title, e.description, e.vendor_id, e.start_date, e.end_date, 
			e.status, e.created_at, e.updated_at, e.picture, e.maxattendees, e.location,
			e.requirements, e.agenda,
			COALESCE(present_count.attendees, 0) as attendees
		FROM events e
		LEFT JOIN (
			SELECT event_id, COUNT(id) as attendees
			FROM attendance
			WHERE attendance_status = 'present'
			GROUP BY event_id
		) present_count ON e.id = present_count.event_id
		WHERE e.vendor_id = $1
		ORDER BY e.start_date ASC
	`

	rows, err := db.Query(query, vendorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.VendorID, &e.StartDate, &e.EndDate, &e.Status, &e.CreatedAt, &e.UpdatedAt,
			&e.Picture, &e.MaxAttendees, &e.Location, &e.Requirements, &e.Agenda, &e.Attendees,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		events = append(events, e)
	}

	return c.JSON(http.StatusOK, events)
}

func getCertificatesByWalletAddressHandler(c echo.Context) error {
	walletAddress := c.Param("walletAddress")
	if walletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "walletAddress is required"})
	}

	var userID int
	err := db.QueryRow(`SELECT id FROM users WHERE wallet_address = $1`, walletAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
	locationStr := c.FormValue("location")

	// Parse requirements and agenda as string from form
	requirementsStr := c.FormValue("requirements")
	agendaStr := c.FormValue("agenda")

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
	fmt.Println("Received create event request:")
	fmt.Println("Title:", title)
	fmt.Println("Description:", description)
	fmt.Println("VendorIDStr:", vendorIDStr)
	fmt.Println("StartDateStr:", startDateStr)
	fmt.Println("EndDateStr:", endDateStr)
	fmt.Println("Status:", status)
	fmt.Println("MaxAttendeesStr:", maxAttendeesStr)
	fmt.Println("RequirementsStr:", requirementsStr)
	fmt.Println("AgendaStr:", agendaStr)
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

	// Parse requirementsStr and agendaStr as JSON
	var requirementsJSON json.RawMessage
	if requirementsStr != "" {
		if err := json.Unmarshal([]byte(requirementsStr), &requirementsJSON); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid requirements JSON"})
		}
	}
	var agendaJSON json.RawMessage
	if agendaStr != "" {
		if err := json.Unmarshal([]byte(agendaStr), &agendaJSON); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid agenda JSON"})
		}
	}

	// Use Event struct and fill all properties including Location, Requirements, Agenda
	event := Event{
		Title:        title,
		Description:  description,
		VendorID:     vendorID,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		Picture:      filename,
		MaxAttendees: maxAttendees,
		Location:     locationStr,
		Requirements: requirementsJSON,
		Agenda:       agendaJSON,
	}

	// Updated query with requirements, agenda, and location
	query := `
	    INSERT INTO events (
			title, description, vendor_id, start_date, end_date,
			status, picture, maxattendees, location, requirements, agenda
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	    RETURNING id, created_at, updated_at
	`

	row := db.QueryRow(query,
		event.Title,
		event.Description,
		event.VendorID,
		event.StartDate,
		event.EndDate,
		event.Status,
		event.Picture,
		event.MaxAttendees,
		event.Location,
		event.Requirements,
		event.Agenda,
	)

	fmt.Println("QueryRow executed, now scanning result...")

	err = row.Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		fmt.Println("DB Scan failed:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB scan failed", "details": err.Error()})
	}

	return c.JSON(http.StatusCreated, event)
}

// Handler to get event detail by ID
func getEventDetailHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event ID is required"})
	}

	query := `
      SELECT 
        e.id, v.vendor_name, e.title, e.description, e.vendor_id,
        e.start_date, e.end_date, e.status, e.created_at, e.updated_at,
        e.picture, e.maxattendees, e.location, e.requirements, e.agenda,
        COUNT(a.id) as attendees
      FROM events e
      LEFT JOIN vendors v ON e.vendor_id = v.id
      LEFT JOIN attendance a ON a.event_id = e.id AND a.attendance_status = 'present'
      WHERE e.id = $1
      GROUP BY e.id, v.vendor_name
    `
	row := db.QueryRow(query, id)

	var dbEvent struct {
		ID           int
		Organizer    sql.NullString
		Title        string
		Description  string
		VendorID     int
		StartDate    time.Time
		EndDate      time.Time
		Status       string
		CreatedAt    time.Time
		UpdatedAt    time.Time
		Picture      string
		MaxAttendees int
		Location     string
		Requirements json.RawMessage
		Agenda       json.RawMessage
		Attendees    int
	}

	err := row.Scan(
		&dbEvent.ID, &dbEvent.Organizer, &dbEvent.Title, &dbEvent.Description, &dbEvent.VendorID,
		&dbEvent.StartDate, &dbEvent.EndDate, &dbEvent.Status, &dbEvent.CreatedAt, &dbEvent.UpdatedAt,
		&dbEvent.Picture, &dbEvent.MaxAttendees, &dbEvent.Location, &dbEvent.Requirements, &dbEvent.Agenda,
		&dbEvent.Attendees,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "event not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 🔍 Count whitelisted entries for this event
	var whitelistedCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM whitelist WHERE event_id = $1`, dbEvent.ID).Scan(&whitelistedCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count whitelist entries"})
	}

	type EventDetailResponse struct {
		ID           int             `json:"id"`
		Organizer    *string         `json:"organizer"`
		Title        string          `json:"title"`
		Description  string          `json:"description"`
		VendorID     int             `json:"vendor_id"`
		StartDate    time.Time       `json:"start_date"`
		EndDate      time.Time       `json:"end_date"`
		Status       string          `json:"status"`
		CreatedAt    time.Time       `json:"created_at"`
		UpdatedAt    time.Time       `json:"updated_at"`
		Picture      string          `json:"picture"`
		MaxAttendees int             `json:"maxattendees"`
		Location     string          `json:"location"`
		Attendees    int             `json:"attendees"`
		Whitelisted  int             `json:"whitelisted"`
		Requirements json.RawMessage `json:"requirements,omitempty"`
		Agenda       json.RawMessage `json:"agenda,omitempty"`
	}

	response := EventDetailResponse{
		ID:           dbEvent.ID,
		Title:        dbEvent.Title,
		Description:  dbEvent.Description,
		VendorID:     dbEvent.VendorID,
		StartDate:    dbEvent.StartDate,
		EndDate:      dbEvent.EndDate,
		Status:       dbEvent.Status,
		CreatedAt:    dbEvent.CreatedAt,
		UpdatedAt:    dbEvent.UpdatedAt,
		Picture:      dbEvent.Picture,
		MaxAttendees: dbEvent.MaxAttendees,
		Location:     dbEvent.Location,
		Attendees:    dbEvent.Attendees,
		Whitelisted:  whitelistedCount,
		Requirements: dbEvent.Requirements,
		Agenda:       dbEvent.Agenda,
	}

	if dbEvent.Organizer.Valid {
		response.Organizer = &dbEvent.Organizer.String
	}

	return c.JSON(http.StatusOK, response)
}

func createWhitelistHandler(c echo.Context) error {
	var input struct {
		EventID       int    `json:"event_id"`
		WalletAddress string `json:"wallet_address"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if input.EventID == 0 || input.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event_id and wallet_address are required"})
	}

	// Get user_id based on wallet_address
	var userID int
	err := db.QueryRow(`SELECT id FROM users WHERE wallet_address = $1`, input.WalletAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Wallet address is not registered as a user"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find user"})
	}

	// Get maxAttendees quota from event
	var maxAttendees int
	err = db.QueryRow(`SELECT maxattendees FROM events WHERE id = $1`, input.EventID).Scan(&maxAttendees)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read event data"})
	}

	// Count approved whitelist entries
	var approvedCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM whitelist WHERE event_id = $1 AND status = 'approved'`, input.EventID).Scan(&approvedCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to count whitelist quota"})
	}

	status := "pending"
	message := "Whitelist successful, but quota is full. You are on the waiting list."
	if approvedCount < maxAttendees {
		status = "approved"
		message = "Whitelist successful! You are registered as an event participant."
	}

	// Insert into whitelist
	var wlID int
	var createdAt, updatedAt time.Time
	err = db.QueryRow(`
		INSERT INTO whitelist (event_id, user_id, wallet_address, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, input.EventID, userID, input.WalletAddress, status).Scan(&wlID, &createdAt, &updatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "You have already registered for this event"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register whitelist"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": message,
		"data": map[string]interface{}{
			"id":             wlID,
			"event_id":       input.EventID,
			"user_id":        userID,
			"wallet_address": input.WalletAddress,
			"status":         status,
			"created_at":     createdAt,
			"updated_at":     updatedAt,
		},
	})
}