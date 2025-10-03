package handler

import (
	"database/sql"
	"time"

	"api-certify-nft-be/config"
	"api-certify-nft-be/pkg/resp"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	cfg *config.Config
	db  *sql.DB
}

func NewHealthHandler(cfg *config.Config, db *sql.DB) *HealthHandler {
	return &HealthHandler{cfg: cfg, db: db}
}

func (h *HealthHandler) Health(c echo.Context) error {
	dbStatus := "ok"
	if err := h.db.Ping(); err != nil {
		dbStatus = err.Error()
	}
	return resp.OK(c, map[string]interface{}{
		"status": "ok", "database": dbStatus, "server_time": time.Now().Format(time.RFC3339),
		"app_version": h.cfg.AppVersion, "environment": h.cfg.AppEnv,
		"hostname": c.Request().Host,
	})
}
