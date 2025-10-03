package http

import (
	"github.com/labstack/echo/v4"

	"api-certify-nft-be/internal/handler"
)

type Deps struct {
	Auth   *handler.AuthHandler
	Event  *handler.EventHandler
	User   *handler.UserHandler
	Vendor *handler.VendorHandler
	Cert   *handler.CertificateHandler
	Health *handler.HealthHandler
}

func Register(e *echo.Echo, d Deps) {
	e.Static("/uploads", "uploads")
	api := e.Group("/api")

	// Auth
	api.POST("/auth/login", d.Auth.Login)
	api.POST("/users/register", d.Auth.RegisterUser)
	api.POST("/vendors/register", d.Auth.RegisterVendor)

	// Events
	api.GET("/events/all", d.Event.GetAll)
	api.POST("/events/create", d.Event.Create)
	api.GET("/events/:id", d.Event.GetDetail)
	api.POST("/events/cancel/:id", d.Event.Cancel)
	api.POST("/events/:id/update", d.Event.UpdateStatus)
	api.GET("/attendance/event/:event_id", d.Event.AttendanceByEvent)
	api.GET("/events/:id/whitelist", d.Event.WhitelistByEvent)

	// User
	api.POST("/users/attend", d.Event.MarkAttendance)
	api.GET("/users/:walletAddress", d.User.GetByWallet)
	api.GET("/users/:walletAddress/events", d.User.GetEvents)
	api.GET("/users/:walletAddress/certificate", d.User.GetCertificates)
	api.POST("/users/whitelist", d.Event.CreateWhitelist)
	api.POST("/users/whitelist/cancel", d.Event.CancelWhitelist)
	api.GET("/users/:walletAddress/events/:eventId/attendance-status", d.User.GetAttendanceStatus)
	api.GET("/users/:walletAddress/events/:eventId/whitelist-status", d.User.GetWhitelistStatus)

	// Vendor
	api.GET("/vendors/:walletAddress/events", d.Vendor.GetEvents)

	// Certificate
	api.GET("/certificate/all", d.Cert.GetAll)

	// Health
	api.GET("/health", d.Health.Health)
}
