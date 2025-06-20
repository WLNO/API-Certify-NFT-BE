package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

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
	dsn := "postgres://admin:10062004Dk-@103.175.219.68:5432/certify_nft?sslmode=disable"
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
	e.GET("/events", getEventsHandler)

	e.Logger.Fatal(e.Start(":8080"))
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
