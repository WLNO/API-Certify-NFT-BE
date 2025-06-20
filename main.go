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
}

var db *sql.DB

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

	e.Logger.Fatal(e.Start(":4002"))
}
func getEventsHandler(c echo.Context) error {
	rows, err := db.Query(`SELECT id, title, description, vendor_id, start_date, end_date, status, created_at, updated_at, picture, maxattendees FROM events`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VendorID, &e.StartDate, &e.EndDate, &e.Status, &e.CreatedAt, &e.UpdatedAt, &e.Picture, &e.MaxAttendees)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		events = append(events, e)
	}

	return c.JSON(http.StatusOK, events)
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
