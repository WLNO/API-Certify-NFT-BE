package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
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

//==============================================
// TYPE DEFINITIONS
//==============================================
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
	Token        string          `json:"token,omitempty"`
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

var db *sql.DB
var appStartTime time.Time

//==============================================
// HELPER FUNCTIONS
//==============================================
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

// Helper function to determine the dynamic status of an event
func calculateStatus(dbStatus string, startDate, endDate time.Time) string {
	// Jika status sudah di-set manual oleh vendor (final), langsung kembalikan.
	if dbStatus == "ended" || dbStatus == "canceled" {
		return dbStatus
	}

	now := time.Now()

	// Jika waktu sekarang sudah melewati tanggal selesai event
	if now.After(endDate) {
		return "minting"
	}

	// Jika waktu sekarang berada di antara tanggal mulai dan selesai
	if now.After(startDate) && now.Before(endDate) {
		return "ongoing"
	}

	// Jika tidak, berarti event belum dimulai
	return "upcoming"
}

//==============================================
// MAIN FUNCTION
//==============================================
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
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve static files from the "uploads" directory
	e.Static("/uploads", "uploads")

	// Group endpoints
	api := e.Group("/api")

	// Auth endpoints
	api.POST("/auth/login", loginHandler)

	// Event endpoints
	api.GET("/events/all", getEventsHandler)
	api.POST("/events/create", createEventHandler)
	api.GET("/events/:id", getEventDetailHandler)
	api.POST("/events/cancel/:id", cancelEventHandler)
	api.POST("/events/:id/update", updateEventStatusHandler)
	api.GET("/attendance/event/:event_id", getAttendanceByEventHandler)
	api.GET("/events/:id/whitelist", getUserByWhitelist)

	// User endpoints
	api.POST("/users/register", registerUserHandler)
	api.POST("/users/attend", markAttendanceHandler)
	api.GET("/users/:walletAddress", getUserByWalletAddressHandler)
	api.GET("/users/:walletAddress/events", getEventsByWalletAddressHandler)
	api.GET("/users/:walletAddress/certificate", getCertificatesByWalletAddressHandler)
	api.POST("/users/whitelist", createWhitelistHandler)
	api.POST("/users/whitelist/cancel", cancelWhitelistHandler)

	// Vendor endpoints
	api.POST("/vendors/register", registerVendorHandler)
	api.GET("/vendors/:walletAddress/events", getEventsByVendorWalletAddressHandler)

	// Health check endpoint
	api.GET("/health", healthHandler)

	// Endpoint: cek status absen user
	api.GET("/users/:walletAddress/events/:eventId/attendance-status", getUserAttendanceStatusHandler)
	// Endpoint: cek status whitelist user
	api.GET("/users/:walletAddress/events/:eventId/whitelist-status", getUserWhitelistStatusHandler)

	// Set waktu mulai aplikasi untuk uptime
	appStartTime = time.Now()

	e.Logger.Fatal(e.Start(":4002"))
}

//==============================================
// AUTH HANDLERS
//==============================================
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
			"isNewUser":      false,
			"wallet_address": body.WalletAddress,
			"role":           role,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"isNewUser": true,
	})
}

//==============================================
// EVENT HANDLERS
//==============================================
func createEventHandler(c echo.Context) error {
	title := c.FormValue("title")
	description := c.FormValue("description")
	walletAddress := c.FormValue("wallet_address")
	startDateStr := c.FormValue("start_date")
	endDateStr := c.FormValue("end_date")
	maxAttendeesStr := c.FormValue("maxattendees")
	location := c.FormValue("location")
	requirementsStr := c.FormValue("requirements")
	agendaStr := c.FormValue("agenda")

	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
	}
	if description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "description is required"})
	}
	if walletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
	}
	if startDateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_date is required"})
	}
	if endDateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "end_date is required"})
	}
	if maxAttendeesStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "maxattendees is required"})
	}
	if location == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "location is required"})
	}

	var vendorID int
	err := db.QueryRow("SELECT id FROM vendors WHERE wallet_address = $1", walletAddress).Scan(&vendorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Vendor not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	layout := "2006-01-02T15:04"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format. Use YYYY-MM-DDTHH:MM"})
	}
	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format. Use YYYY-MM-DDTHH:MM"})
	}
	maxAttendees, err := strconv.Atoi(maxAttendeesStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid maxattendees format"})
	}

	file, err := c.FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "picture is required"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open picture file"})
	}
	defer src.Close()

	picturePath := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
	dst, err := os.Create(picturePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save picture file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to copy picture file"})
	}

	// Generate a random 5-character alphanumeric token.
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const length = 5
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}
		ret[i] = letters[num.Int64()]
	}
	token := string(ret)

	query := `INSERT INTO events (title, description, vendor_id, start_date, end_date, picture, maxattendees, location, requirements, agenda, token) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at, updated_at`
	var eventID int
	var createdAt, updatedAt time.Time
	err = db.QueryRow(query, title, description, vendorID, startDate, endDate, picturePath, maxAttendees, location, json.RawMessage(requirementsStr), json.RawMessage(agendaStr), token).Scan(&eventID, &createdAt, &updatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	type EventCreateResponse struct {
		ID            int             `json:"id"`
		Title         string          `json:"title"`
		Description   string          `json:"description"`
		VendorID      int             `json:"vendor_id"`
		WalletAddress string          `json:"wallet_address"`
		StartDate     time.Time       `json:"start_date"`
		EndDate       time.Time       `json:"end_date"`
		Status        string          `json:"status"`
		CreatedAt     time.Time       `json:"created_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
		Picture       string          `json:"picture"`
		MaxAttendees  int             `json:"maxattendees"`
		Location      string          `json:"location"`
		Requirements  json.RawMessage `json:"requirements,omitempty"`
		Agenda        json.RawMessage `json:"agenda,omitempty"`
		Token         string          `json:"token"`
	}

	eventResponse := EventCreateResponse{
		ID:            eventID,
		Title:         title,
		Description:   description,
		VendorID:      vendorID,
		WalletAddress: walletAddress,
		StartDate:     startDate,
		EndDate:       endDate,
		Status:        "upcoming",
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Picture:       "https://api.gpadaka.com/" + picturePath,
		MaxAttendees:  maxAttendees,
		Location:      location,
		Requirements:  json.RawMessage(requirementsStr),
		Agenda:        json.RawMessage(agendaStr),
		Token:         token,
	}

	return c.JSON(http.StatusCreated, eventResponse)
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
		var event EventWithOrganizer
		var dbStatus string
		var startDate, endDate time.Time
		var picturePath string

		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.VendorID, &startDate, &endDate,
			&dbStatus, &event.CreatedAt, &event.UpdatedAt, &picturePath, &event.MaxAttendees, &event.Location,
			&event.Organizer, &event.Whitelisted,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		event.StartDate = startDate
		event.EndDate = endDate
		event.Status = calculateStatus(dbStatus, startDate, endDate)
		event.Picture = "https://api.gpadaka.com/" + picturePath
		events = append(events, event)
	}

	return c.JSON(http.StatusOK, events)
}

func getEventDetailHandler(c echo.Context) error {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event ID"})
	}

	query := `
		SELECT
			e.id,
			v.vendor_name,
			e.title,
			e.description,
			e.vendor_id,
			e.start_date,
			e.end_date,
			e.status,
			e.created_at,
			e.updated_at,
			e.picture,
			e.maxattendees,
			e.location,
			COALESCE(att_counts.attendees, 0) as attendees,
			COALESCE(wl_counts.whitelisted, 0) as whitelisted,
			COALESCE(mint_counts.minted, 0) as minted,
			e.requirements,
			e.agenda,
			e.token,
			COALESCE(ec.url_certificate, '') as url_certificate
		FROM events e
		LEFT JOIN vendors v ON e.vendor_id = v.id
		LEFT JOIN (SELECT event_id, COUNT(*) as attendees FROM attendance WHERE attendance_status = 'present' GROUP BY event_id) att_counts ON e.id = att_counts.event_id
		LEFT JOIN (SELECT event_id, COUNT(*) as whitelisted FROM whitelist WHERE status = 'approved' GROUP BY event_id) wl_counts ON e.id = wl_counts.event_id
		LEFT JOIN (SELECT event_id, COUNT(*) as minted FROM certificates WHERE mint_status = 'success' GROUP BY event_id) mint_counts ON e.id = mint_counts.event_id
		LEFT JOIN event_certificates ec ON e.id = ec.event_id
		WHERE e.id = $1
	`
	var event struct {
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
		Minted       int             `json:"minted"`
		Requirements json.RawMessage `json:"requirements,omitempty"`
		Agenda       json.RawMessage `json:"agenda,omitempty"`
		Token        string          `json:"token"`
		UrlCertificate *string       `json:"url_certificate,omitempty"`
	}
	var dbStatus string
	var startDate, endDate time.Time
	var picturePath string

	err = db.QueryRow(query, eventID).Scan(
		&event.ID, &event.Organizer, &event.Title, &event.Description,
		&event.VendorID, &startDate, &endDate, &dbStatus,
		&event.CreatedAt, &event.UpdatedAt, &picturePath, &event.MaxAttendees,
		&event.Location, &event.Attendees, &event.Whitelisted, &event.Minted,
		&event.Requirements, &event.Agenda, &event.Token, &event.UrlCertificate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	event.StartDate = startDate
	event.EndDate = endDate
	event.Status = calculateStatus(dbStatus, startDate, endDate)
	event.Picture = "https://api.gpadaka.com/" + picturePath

	return c.JSON(http.StatusOK, event)
}

func updateEventStatusHandler(c echo.Context) error {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil || eventID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event ID"})
	}

	var payload struct {
		Status        string `json:"status"`
		WalletAddress string `json:"wallet_address"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if payload.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
	}

	// Validate status
	if payload.Status != "ended" && payload.Status != "minting" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid status. Must be 'ended' or 'minting'"})
	}

	// 1. Get vendor_id from wallet_address
	var vendorID int
	err = db.QueryRow("SELECT id FROM vendors WHERE wallet_address = $1", payload.WalletAddress).Scan(&vendorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Vendor not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify vendor"})
	}

	// 2. Verify vendor owns the event
	var dbVendorID int
	var dbStatus string
	var startDate, endDate time.Time
	err = db.QueryRow("SELECT vendor_id, status, start_date, end_date FROM events WHERE id = $1", eventID).Scan(&dbVendorID, &dbStatus, &startDate, &endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get event details"})
	}

	if dbVendorID != vendorID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not the owner of this event"})
	}

	// 3. Check logic for status change
	calculatedStatus := calculateStatus(dbStatus, startDate, endDate)

	if payload.Status == "ended" && calculatedStatus != "minting" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Event cannot be marked as 'ended' yet. It is currently " + calculatedStatus})
	}

	if payload.Status == "minting" && dbStatus != "ended" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Event is not 'ended', so it cannot be changed back to 'minting'"})
	}

	// 4. Update the status in the database
	_, err = db.Exec("UPDATE events SET status = $1 WHERE id = $2", payload.Status, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update event status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Event status updated successfully to " + payload.Status})
}

func cancelEventHandler(c echo.Context) error {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event ID"})
	}

	var payload struct {
		WalletAddress string `json:"wallet_address"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	var vendorID int
	err = db.QueryRow("SELECT id FROM vendors WHERE wallet_address = $1", payload.WalletAddress).Scan(&vendorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Vendor not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify vendor"})
	}

	var dbVendorID int
	err = db.QueryRow("SELECT vendor_id FROM events WHERE id = $1", eventID).Scan(&dbVendorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get event details"})
	}

	if dbVendorID != vendorID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not the owner of this event"})
	}

	_, err = db.Exec("UPDATE events SET status = 'canceled' WHERE id = $1", eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to cancel event"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Event canceled successfully"})
}

func getAttendanceByEventHandler(c echo.Context) error {
	eventIDStr := c.Param("event_id")
	if eventIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Event ID parameter is required in the URL."})
	}
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil || eventID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Event ID must be a valid positive number."})
	}

	query := `
		SELECT
			u.id AS user_id,
			u.name,
			u.wallet_address,
			COALESCE(a.attendance_status, 'absent') AS attend_status,
			a.created_at AS attended_at
		FROM whitelist w
		JOIN users u ON w.user_id = u.id
		LEFT JOIN attendance a ON a.event_id = w.event_id AND a.user_id = w.user_id AND a.attendance_status = 'present'
		WHERE w.event_id = $1
		ORDER BY u.name ASC
	`

	rows, err := db.Query(query, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to retrieve attendance data for this event. Please try again later."})
	}
	defer rows.Close()

	type AttendanceUser struct {
		UserID        int        `json:"user_id"`
		Name          string     `json:"name"`
		WalletAddress string     `json:"wallet_address"`
		AttendStatus  string     `json:"attend_status"`
		AttendedAt    *time.Time `json:"attended_at"`
	}

	var result []AttendanceUser
	for rows.Next() {
		var u AttendanceUser
		var attendedAt sql.NullTime
		if err := rows.Scan(&u.UserID, &u.Name, &u.WalletAddress, &u.AttendStatus, &attendedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process attendance data. Please contact support if this continues."})
		}
		if attendedAt.Valid {
			u.AttendedAt = &attendedAt.Time
		} else {
			u.AttendedAt = nil
		}
		result = append(result, u)
	}

	// If no whitelist entries, return empty array
	return c.JSON(http.StatusOK, result)
}

func getUserByWhitelist(c echo.Context) error {
	eventIDStr := c.Param("id")
	if eventIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Event ID parameter is required in the URL."})
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event ID"})
	}

	query := `SELECT u.id, u.name, u.email, u.wallet_address, w.status, w.created_at FROM users u JOIN whitelist w ON u.id = w.user_id WHERE w.event_id = $1`

	rows, err := db.Query(query, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve whitelisted users"})
	}
	defer rows.Close()

	type WhitelistedUser struct {
		ID            int       `json:"id"`
		UserID        int       `json:"user_id"`
		Name          string    `json:"name"`
		Email         string    `json:"email"`
		WalletAddress string    `json:"wallet_address"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
	}

	var result []WhitelistedUser
	for rows.Next() {
		var u WhitelistedUser
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.WalletAddress, &u.Status, &u.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to scan user data"})
		}
		result = append(result, u)
	}

	// If no whitelist entries, return empty array
	return c.JSON(http.StatusOK, result)
}

//==============================================
// USER HANDLERS
//==============================================
func registerUserHandler(c echo.Context) error {
	var u User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if u.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}
	if u.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email is required"})
	}
	if u.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
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

func getUserByWalletAddressHandler(c echo.Context) error {
	walletAddress := c.Param("walletAddress")
	if walletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
	}

	var user User
	query := `SELECT id, name, email, wallet_address, created_at, updated_at FROM users WHERE wallet_address = $1`
	err := db.QueryRow(query, walletAddress).Scan(&user.ID, &user.Name, &user.Email, &user.WalletAddress, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
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
		var event EventResponse
		var dbStatus string
		var startDate, endDate time.Time
		var picturePath string
		var userStatus sql.NullString

		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.VendorID, &startDate, &endDate,
			&dbStatus, &event.CreatedAt, &event.UpdatedAt, &picturePath, &event.MaxAttendees, &event.Location,
			&event.Requirements, &event.Agenda, &event.Attendees,
			&userStatus,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		event.StartDate = startDate
		event.EndDate = endDate
		event.Status = calculateStatus(dbStatus, startDate, endDate)
		event.Picture = "https://api.gpadaka.com/" + picturePath
		if userStatus.Valid {
			event.UserStatus = userStatus.String
		} else {
			event.UserStatus = "not_whitelisted"
		}
		events = append(events, event)
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

func markAttendanceHandler(c echo.Context) error {
	var payload struct {
		EventToken    string `json:"event_token"`
		WalletAddress string `json:"wallet_address"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if payload.EventToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event_token is required"})
	}
	if payload.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
	}

	var eventID int
	var eventStatus string
	var startDate, endDate time.Time
	err := db.QueryRow("SELECT id, status, start_date, end_date FROM events WHERE token = $1", payload.EventToken).Scan(&eventID, &eventStatus, &startDate, &endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid event token"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate event token"})
	}

	// Recalculate status to ensure it's current
	liveStatus := calculateStatus(eventStatus, startDate, endDate)
	if liveStatus != "ongoing" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "This event is not currently ongoing. Attendance cannot be marked."})
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE wallet_address = $1", payload.WalletAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User with this wallet address not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find user"})
	}

	var whitelistStatus string
	err = db.QueryRow("SELECT status FROM whitelist WHERE event_id = $1 AND user_id = $2", eventID, userID).Scan(&whitelistStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not whitelisted for this event"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check whitelist status"})
	}

	if whitelistStatus != "approved" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Your whitelist status is not approved"})
	}

	tx, err := db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}

	_, err = tx.Exec(`
        INSERT INTO attendance (event_id, user_id, attendance_status, token_input)
        VALUES ($1, $2, 'present', $3)
        ON CONFLICT (event_id, user_id) 
        DO UPDATE SET attendance_status = 'present', token_input = $3, updated_at = NOW()`,
		eventID, userID, payload.EventToken)

	if err != nil {
		log.Println("Failed to mark attendance:", err)
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to mark attendance"})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Attendance marked successfully",
	})
}

func createWhitelistHandler(c echo.Context) error {
	var input struct {
		EventID       int    `json:"event_id"`
		WalletAddress string `json:"wallet_address"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if input.EventID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event_id is required"})
	}
	if input.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
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

func cancelWhitelistHandler(c echo.Context) error {
	var payload struct {
		EventID       int    `json:"event_id"`
		WalletAddress string `json:"wallet_address"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}

	if payload.EventID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event_id is required"})
	}
	if payload.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
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

	return c.JSON(http.StatusOK, map[string]string{"message": "user whitelist canceled successfully"})
}

//==============================================
// VENDOR HANDLERS
//==============================================
func registerVendorHandler(c echo.Context) error {
	var v Vendor
	if err := c.Bind(&v); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if v.VendorName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "vendor_name is required"})
	}
	if v.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email is required"})
	}
	if v.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wallet_address is required"})
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
			e.requirements, e.agenda, e.token,
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

	type EventWithCert struct {
		Event
		CertificateUploaded bool `json:"certificate_uploaded"`
	}
	var events []EventWithCert
	for rows.Next() {
		var event EventWithCert
		var dbStatus string
		var startDate, endDate time.Time
		var picturePath string
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.VendorID, &startDate, &endDate,
			&dbStatus, &event.CreatedAt, &event.UpdatedAt, &picturePath, &event.MaxAttendees, &event.Location,
			&event.Requirements, &event.Agenda, &event.Token, &event.Attendees,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning event"})
		}
		event.StartDate = startDate
		event.EndDate = endDate
		event.Status = calculateStatus(dbStatus, startDate, endDate)
		event.Picture = "https://api.gpadaka.com/" + picturePath

		// Tambahkan pengecekan certificate_uploaded
		var certCount int
		db.QueryRow("SELECT COUNT(*) FROM event_certificates WHERE event_id = $1 AND url_certificate IS NOT NULL AND url_certificate <> ''", event.ID).Scan(&certCount)
		event.CertificateUploaded = certCount > 0
		events = append(events, event)
	}
	return c.JSON(http.StatusOK, events)
}

//==============================================
// HEALTH CHECK HANDLER
//==============================================
func healthHandler(c echo.Context) error {
	// Cek koneksi database
	dbStatus := "ok"
	dbErr := db.Ping()
	if dbErr != nil {
		dbStatus = dbErr.Error()
	}

	// Info versi aplikasi (bisa di-set manual atau dari env)
	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "dev"
	}

	// Info environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Info waktu server
	serverTime := time.Now().Format(time.RFC3339)

	// Info host
	hostname, _ := os.Hostname()

	// Info uptime
	uptime := time.Since(appStartTime).String()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"database": dbStatus,
		"server_time": serverTime,
		"app_version": appVersion,
		"environment": env,
		"hostname": hostname,
		"uptime": uptime,
	})
}

// Handler: cek status absen user
func getUserAttendanceStatusHandler(c echo.Context) error {
	walletAddress := c.Param("walletAddress")
	eventIdStr := c.Param("eventId")
	if walletAddress == "" || eventIdStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "walletAddress and eventId are required"})
	}
	eventID, err := strconv.Atoi(eventIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "eventId must be a number"})
	}
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE wallet_address = $1", walletAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, map[string]bool{"attended": false})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find user"})
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM attendance WHERE event_id = $1 AND user_id = $2 AND attendance_status = 'present'", eventID, userID).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check attendance"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"attended": count > 0})
}

// Handler: cek status whitelist user
func getUserWhitelistStatusHandler(c echo.Context) error {
	walletAddress := c.Param("walletAddress")
	eventIdStr := c.Param("eventId")
	if walletAddress == "" || eventIdStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "walletAddress and eventId are required"})
	}
	eventID, err := strconv.Atoi(eventIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "eventId must be a number"})
	}
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE wallet_address = $1", walletAddress).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, map[string]bool{"whitelisted": false})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find user"})
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM whitelist WHERE event_id = $1 AND user_id = $2 AND status = 'approved'", eventID, userID).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check whitelist"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"whitelisted": count > 0})
}