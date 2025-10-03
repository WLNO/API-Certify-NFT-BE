package app

import (
	"database/sql"

	"api-certify-nft-be/config"
	"api-certify-nft-be/internal/handler"
	routes "api-certify-nft-be/internal/http"
	"api-certify-nft-be/internal/middleware"
	"api-certify-nft-be/internal/repository"
	"api-certify-nft-be/internal/service"

	"github.com/labstack/echo/v4"
)

func NewServer(cfg *config.Config, db *sql.DB) *echo.Echo {
	e := echo.New()
	middleware.Use(e)

	// repo
	authRepo := repository.NewAuthRepo(db)
	evRepo := repository.NewEventRepo(db)

	// service
	authSvc := service.NewAuthService(authRepo)
	evSvc := service.NewEventService(evRepo)

	// handler
	authH := handler.NewAuthHandler(authSvc)
	evH := handler.NewEventHandler(cfg, evSvc)
	userH := handler.NewUserHandler(cfg, evSvc)
	vH := handler.NewVendorHandler(cfg, evSvc)
	cH := handler.NewCertificateHandler(evSvc)
	hH := handler.NewHealthHandler(cfg, db)

	routes.Register(e, routes.Deps{
		Auth: authH, Event: evH, User: userH, Vendor: vH, Cert: cH, Health: hH,
	})
	return e
}
