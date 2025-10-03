package handler

import (
	"api-certify-nft-be/config"
	"api-certify-nft-be/internal/service"
	"api-certify-nft-be/pkg/resp"

	"github.com/labstack/echo/v4"
)

type VendorHandler struct {
	cfg *config.Config
	svc service.EventService
}

func NewVendorHandler(cfg *config.Config, s service.EventService) *VendorHandler {
	return &VendorHandler{cfg: cfg, svc: s}
}

func (h *VendorHandler) GetEvents(c echo.Context) error {
	w := c.Param("walletAddress")
	if w == "" {
		return resp.Err(c, 400, "walletAddress is required")
	}
	data, err := h.svc.GetVendorEvents(w, h.cfg.BaseURL)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}
