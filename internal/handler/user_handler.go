package handler

import (
	"strconv"

	"api-certify-nft-be/config"
	"api-certify-nft-be/internal/service"
	"api-certify-nft-be/pkg/resp"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	cfg *config.Config
	svc service.EventService
}

func NewUserHandler(cfg *config.Config, s service.EventService) *UserHandler {
	return &UserHandler{cfg: cfg, svc: s}
}

func (h *UserHandler) GetByWallet(c echo.Context) error {
	w := c.Param("walletAddress")
	if w == "" {
		return resp.Err(c, 400, "wallet_address is required")
	}
	u, err := h.svc.GetUserByWallet(w)
	if err != nil {
		return resp.Err(c, 404, "user not found")
	}
	return resp.OK(c, u)
}

func (h *UserHandler) GetEvents(c echo.Context) error {
	w := c.Param("walletAddress")
	if w == "" {
		return resp.Err(c, 400, "walletAddress is required")
	}
	data, err := h.svc.GetUserEvents(w, h.cfg.BaseURL)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}

func (h *UserHandler) GetCertificates(c echo.Context) error {
	w := c.Param("walletAddress")
	if w == "" {
		return resp.Err(c, 400, "walletAddress is required")
	}
	data, err := h.svc.GetUserCertificates(w, h.cfg.BaseURL)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}

func (h *UserHandler) GetAttendanceStatus(c echo.Context) error {
	w := c.Param("walletAddress")
	eidStr := c.Param("eventId")
	if w == "" || eidStr == "" {
		return resp.Err(c, 400, "walletAddress and eventId are required")
	}
	eid, err := strconv.Atoi(eidStr)
	if err != nil {
		return resp.Err(c, 400, "eventId must be a number")
	}
	// hit DB langsung via repo? untuk ringkas gunakan service.GetUserEvents dan cek present?
	// Agar efisien, langsung query kecil → gunakan service method kecil? Di sini sederhana:
	data, err := h.svc.AttendanceByEvent(eid)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	att := false
	for _, row := range data {
		if row["wallet_address"] == w && row["attend_status"] == "present" {
			att = true
			break
		}
	}
	return resp.OK(c, map[string]bool{"attended": att})
}

func (h *UserHandler) GetWhitelistStatus(c echo.Context) error {
	w := c.Param("walletAddress")
	eidStr := c.Param("eventId")
	if w == "" || eidStr == "" {
		return resp.Err(c, 400, "walletAddress and eventId are required")
	}
	eid, err := strconv.Atoi(eidStr)
	if err != nil {
		return resp.Err(c, 400, "eventId must be a number")
	}
	rows, err := h.svc.WhitelistByEvent(eid)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	ok := false
	for _, r := range rows {
		if r["wallet_address"] == w && r["status"] == "approved" {
			ok = true
			break
		}
	}
	return resp.OK(c, map[string]bool{"whitelisted": ok})
}
