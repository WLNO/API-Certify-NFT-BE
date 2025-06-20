package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VendorID    int       `json:"vendor_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	e.GET("/api/events/all", getEventsHandler)
	e.POST("/api/events/create", createEventHandler)

	e.Logger.Fatal(e.Start(":4002"))
}

func getEventsHandler(c echo.Context) error {
	rows, err := db.Query(`SELECT id, title, description, vendor_id, start_date, end_date, status, created_at, updated_at FROM events`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VendorID, &e.StartDate, &e.EndDate, &e.Status, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		events = append(events, e)
	}

	return c.JSON(http.StatusOK, events)
}

func createEventHandler(c echo.Context) error {
	var e Event
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Validasi input
	if e.Title == "" || e.VendorID == 0 || e.StartDate.IsZero() || e.EndDate.IsZero() || e.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	if e.StartDate.After(e.EndDate) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_date must be before end_date"})
	}

	// Masukkan data ke database
	query := `
		INSERT INTO events (title, description, vendor_id, start_date, end_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(
		query,
		e.Title,
		e.Description,
		e.VendorID,
		e.StartDate,
		e.EndDate,
		e.Status,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
