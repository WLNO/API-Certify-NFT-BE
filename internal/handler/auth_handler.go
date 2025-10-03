package handler

import (
	"net/http"

	"api-certify-nft-be/internal/service"
	"api-certify-nft-be/pkg/resp"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct{ svc service.AuthService }

func NewAuthHandler(s service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

func (h *AuthHandler) Login(c echo.Context) error {
	var body struct {
		Wallet string `json:"wallet_address"`
	}
	if err := c.Bind(&body); err != nil || body.Wallet == "" {
		return resp.Err(c, http.StatusBadRequest, "wallet_address is required")
	}

	isNew, role, err := h.svc.Login(body.Wallet)
	if err != nil {
		return resp.Err(c, http.StatusInternalServerError, err.Error())
	}

	// Selalu kirim isNewUser + role (role bisa kosong jika new user)
	return resp.OK(c, map[string]interface{}{
		"isNewUser":      isNew,
		"wallet_address": body.Wallet,
		"role":           role,
	})
}

func (h *AuthHandler) RegisterUser(c echo.Context) error {
	var in struct {
		Name, Email, Wallet string `json:"name","email","wallet_address"`
	}
	if err := c.Bind(&in); err != nil {
		return resp.Err(c, 400, "invalid request")
	}
	if in.Name == "" || in.Email == "" || in.Wallet == "" {
		return resp.Err(c, 400, "name,email,wallet_address required")
	}
	// cek duplikat
	reg, err := h.svc.IsWalletRegistered(in.Wallet)
	if err != nil {
		return resp.Err(c, 500, "server error")
	}
	if reg {
		return resp.Err(c, 409, "wallet address already registered with another account")
	}
	id, created, updated, err := h.svc.RegisterUser(in.Name, in.Email, in.Wallet)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.Created(c, map[string]interface{}{"id": id, "name": in.Name, "email": in.Email, "wallet_address": in.Wallet, "created_at": created, "updated_at": updated})
}

func (h *AuthHandler) RegisterVendor(c echo.Context) error {
	var in struct {
		VendorName, Email, Contact, Wallet string `json:"vendor_name","email","contact_info","wallet_address"`
	}
	if err := c.Bind(&in); err != nil {
		return resp.Err(c, 400, "invalid request")
	}
	if in.VendorName == "" || in.Email == "" || in.Wallet == "" {
		return resp.Err(c, 400, "vendor_name,email,wallet_address required")
	}
	reg, err := h.svc.IsWalletRegistered(in.Wallet)
	if err != nil {
		return resp.Err(c, 500, "server error")
	}
	if reg {
		return resp.Err(c, 409, "wallet address already registered with another account")
	}
	id, err := h.svc.RegisterVendor(in.VendorName, in.Email, in.Contact, in.Wallet)
	if err != nil {
		return resp.Err(c, 500, err.Error())
	}
	return resp.Created(c, map[string]interface{}{"id": id, "vendor_name": in.VendorName, "email": in.Email, "contact_info": in.Contact, "wallet_address": in.Wallet})
}
