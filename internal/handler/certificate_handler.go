package handler

import (
	"api-certify-nft-be/internal/service"
	"api-certify-nft-be/pkg/resp"

	"github.com/labstack/echo/v4"
)

type CertificateHandler struct{ svc service.EventService }

func NewCertificateHandler(s service.EventService) *CertificateHandler {
	return &CertificateHandler{svc: s}
}

func (h *CertificateHandler) GetAll(c echo.Context) error {
	data, err := h.svc.GetCertificateAll()
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}
