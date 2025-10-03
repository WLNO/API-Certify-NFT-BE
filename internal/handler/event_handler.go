package handler

import (
	"net/http"
	"strconv"
	"time"

	"api-certify-nft-be/config"
	"api-certify-nft-be/internal/service"
	"api-certify-nft-be/pkg/files"
	"api-certify-nft-be/pkg/resp"

	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	cfg *config.Config
	svc service.EventService
}

func NewEventHandler(cfg *config.Config, s service.EventService) *EventHandler {
	return &EventHandler{cfg: cfg, svc: s}
}

func (h *EventHandler) GetAll(c echo.Context) error {
	data, err := h.svc.List(c.Request().Context(), h.cfg.BaseURL)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}

func (h *EventHandler) GetDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return resp.Err(c, 400, "Invalid event ID")
	}
	data, err := h.svc.Detail(c.Request().Context(), id, h.cfg.BaseURL)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}

func (h *EventHandler) Create(c echo.Context) error {
	title := c.FormValue("title")
	desc := c.FormValue("description")
	wallet := c.FormValue("wallet_address")
	startStr := c.FormValue("start_date")
	endStr := c.FormValue("end_date")
	maxStr := c.FormValue("maxattendees")
	location := c.FormValue("location")
	req := c.FormValue("requirements")
	agenda := c.FormValue("agenda")

	if title == "" || desc == "" || wallet == "" || startStr == "" || endStr == "" || maxStr == "" || location == "" {
		return resp.Err(c, 400, "missing required field")
	}
	layout := "2006-01-02T15:04"
	start, err := time.Parse(layout, startStr)
	if err != nil {
		return resp.Err(c, 400, "Invalid start_date format. Use YYYY-MM-DDTHH:MM")
	}
	end, err := time.Parse(layout, endStr)
	if err != nil {
		return resp.Err(c, 400, "Invalid end_date format. Use YYYY-MM-DDTHH:MM")
	}
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return resp.Err(c, 400, "Invalid maxattendees format")
	}

	fh, err := c.FormFile("picture")
	if err != nil {
		return resp.Err(c, 400, "picture is required")
	}
	path, err := files.SimpanUpload("uploads", fh)
	if err != nil {
		return resp.Err(c, 500, "Failed to save picture")
	}

	ev, err := h.svc.Create(c.Request().Context(), h.cfg.BaseURL, title, desc, wallet, start, end, path, max, location, req, agenda)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.Created(c, ev)
}

func (h *EventHandler) Cancel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return resp.Err(c, 400, "Invalid event ID")
	}
	var p struct {
		Wallet string `json:"wallet_address"`
	}
	if err := c.Bind(&p); err != nil {
		return resp.Err(c, 400, "Invalid payload")
	}
	if p.Wallet == "" {
		return resp.Err(c, 400, "wallet_address is required")
	}
	if err := h.svc.Cancel(id, p.Wallet); err != nil {
		return resp.Err(c, 403, err.Error())
	}
	return resp.Msg(c, "Event canceled successfully")
}

func (h *EventHandler) UpdateStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return resp.Err(c, 400, "Invalid event ID")
	}
	var p struct {
		Status, Wallet string `json:"status","wallet_address"`
	}
	if err := c.Bind(&p); err != nil {
		return resp.Err(c, 400, "Invalid request payload")
	}
	if p.Wallet == "" {
		return resp.Err(c, 400, "wallet_address is required")
	}
	if err := h.svc.UpdateStatus(id, p.Status, p.Wallet); err != nil {
		return resp.Err(c, 403, err.Error())
	}
	return resp.Msg(c, "Event status updated successfully to "+p.Status)
}

func (h *EventHandler) AttendanceByEvent(c echo.Context) error {
	eid, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		return resp.Err(c, 400, "Event ID must be number")
	}
	data, err := h.svc.AttendanceByEvent(eid)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}

func (h *EventHandler) WhitelistByEvent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return resp.Err(c, 400, "Invalid event id")
	}
	data, err := h.svc.WhitelistByEvent(id)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.OK(c, data)
}

func (h *EventHandler) MarkAttendance(c echo.Context) error {
	var p struct {
		Token, Wallet string `json:"event_token","wallet_address"`
	}
	if err := c.Bind(&p); err != nil {
		return resp.Err(c, 400, "invalid request")
	}
	if p.Token == "" || p.Wallet == "" {
		return resp.Err(c, 400, "event_token and wallet_address required")
	}
	if err := h.svc.MarkAttendance(p.Token, p.Wallet); err != nil {
		return resp.Err(c, 403, err.Error())
	}
	return resp.OK(c, map[string]string{"message": "Attendance marked successfully"})
}

func (h *EventHandler) CreateWhitelist(c echo.Context) error {
	var p struct {
		EventID int    `json:"event_id"`
		Wallet  string `json:"wallet_address"`
	}
	if err := c.Bind(&p); err != nil {
		return resp.Err(c, 400, "Invalid body")
	}
	if p.EventID == 0 || p.Wallet == "" {
		return resp.Err(c, 400, "event_id and wallet_address required")
	}
	status, msg, id, cAt, uAt, err := h.svc.CreateWhitelist(p.EventID, p.Wallet)
	if err != nil {
		return resp.Err(c, errCode(err), err.Error())
	}
	return resp.Created(c, map[string]interface{}{
		"message": msg,
		"data": map[string]interface{}{
			"id": id, "event_id": p.EventID, "wallet_address": p.Wallet, "status": status, "created_at": cAt, "updated_at": uAt,
		},
	})
}

func errCode(err error) int {
	if err.Error() == "already registered" {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func (h *EventHandler) CancelWhitelist(c echo.Context) error {
	var p struct {
		EventID int    `json:"event_id"`
		Wallet  string `json:"wallet_address"`
	}
	if err := c.Bind(&p); err != nil {
		return resp.Err(c, 400, "invalid payload")
	}
	if p.EventID == 0 || p.Wallet == "" {
		return resp.Err(c, 400, "event_id and wallet_address required")
	}
	if err := h.svc.CancelWhitelist(p.EventID, p.Wallet); err != nil {
		return resp.Err(c, 404, err.Error())
	}
	return resp.Msg(c, "user whitelist canceled successfully")
}
