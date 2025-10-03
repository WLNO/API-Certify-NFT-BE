package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"strings"
	"time"

	"api-certify-nft-be/internal/model"
	"api-certify-nft-be/internal/repository"
	"api-certify-nft-be/pkg/timeutil"
)

type EventService interface {
	List(ctx context.Context, baseURL string) ([]model.Event, error)
	Detail(ctx context.Context, id int, baseURL string) (*model.Event, error)
	Create(ctx context.Context, baseURL, title, desc, wallet string, start, end time.Time, picturePath string, max int, location string, req, agenda string) (*model.Event, error)
	Cancel(eventID int, vendorWallet string) error
	UpdateStatus(eventID int, status, vendorWallet string) error
	MarkAttendance(token, wallet string) error
	CreateWhitelist(eventID int, wallet string) (status, message string, id int, created, updated time.Time, err error)
	CancelWhitelist(eventID int, wallet string) error
	AttendanceByEvent(eventID int) ([]map[string]interface{}, error)
	WhitelistByEvent(eventID int) ([]map[string]interface{}, error)
	GetUserByWallet(wallet string) (*model.User, error)
	GetUserEvents(wallet, baseURL string) ([]model.Event, error)
	GetUserCertificates(wallet, baseURL string) ([]model.CertificateWithEvent, error)
	GetVendorEvents(wallet, baseURL string) ([]model.Event, error)
	GetCertificateAll() ([]model.CertificateWithEvent, error)
}

type eventSvc struct{ repo repository.EventRepo }

func NewEventService(r repository.EventRepo) EventService { return &eventSvc{repo: r} }

func (s *eventSvc) List(ctx context.Context, baseURL string) ([]model.Event, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.Event, 0, len(rows))
	for _, r := range rows {
		out = append(out, model.Event{
			ID: r.ID, Title: r.Title, Description: r.Description, VendorID: r.VendorID,
			StartDate: r.StartDate, EndDate: r.EndDate,
			Status:    timeutil.HitungStatus(r.DBStatus, r.StartDate, r.EndDate),
			CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
			Picture:      baseURL + "/" + r.PicturePath,
			MaxAttendees: r.MaxAttendees, Location: r.Location,
			Organizer: &r.Organizer, Whitelisted: r.Whitelisted,
		})
	}
	return out, nil
}

func (s *eventSvc) Detail(ctx context.Context, id int, baseURL string) (*model.Event, error) {
	r, err := s.repo.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	ev := &model.Event{
		ID: r.ID, Organizer: r.Organizer, Title: r.Title, Description: r.Description, VendorID: r.VendorID,
		StartDate: r.StartDate, EndDate: r.EndDate,
		Status:    timeutil.HitungStatus(r.DBStatus, r.StartDate, r.EndDate),
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
		Picture: baseURL + "/" + r.PicturePath, MaxAttendees: r.MaxAttendees, Location: r.Location,
		Attendees: r.Attendees, Whitelisted: r.Whitelisted, Minted: r.Minted,
		Requirements: r.Requirements, Agenda: r.Agenda, Token: r.Token, UrlCertificate: r.UrlCertificate,
	}
	return ev, nil
}

func randToken(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

func (s *eventSvc) Create(ctx context.Context, baseURL, title, desc, wallet string, start, end time.Time, picturePath string, max int, location string, req, agenda string) (*model.Event, error) {
	vendorID, err := s.repo.FindVendorIDByWallet(wallet)
	if err != nil {
		return nil, err
	}
	token, err := randToken(5)
	if err != nil {
		return nil, err
	}

	var reqJSON, agendaJSON json.RawMessage
	if req != "" {
		reqJSON = json.RawMessage(req)
	}
	if agenda != "" {
		agendaJSON = json.RawMessage(agenda)
	}

	id, created, updated, err := s.repo.Insert(ctx, vendorID, title, desc, start, end, picturePath, max, location, reqJSON, agendaJSON, token)
	if err != nil {
		return nil, err
	}

	return &model.Event{
		ID: id, Title: title, Description: desc, VendorID: vendorID, StartDate: start, EndDate: end,
		Status: "upcoming", CreatedAt: created, UpdatedAt: updated,
		Picture: baseURL + "/" + picturePath, MaxAttendees: max, Location: location,
		Requirements: reqJSON, Agenda: agendaJSON, Token: token,
	}, nil
}

func (s *eventSvc) Cancel(eventID int, vendorWallet string) error {
	vendorID, err := s.repo.FindVendorIDByWallet(vendorWallet)
	if err != nil {
		return err
	}
	ownerID, _, _, _, err := s.repo.GetEventOwnerStatusDates(eventID)
	if err != nil {
		return err
	}
	if ownerID != vendorID {
		return ErrForbiddenOwner
	}
	return s.repo.Cancel(eventID)
}

var (
	ErrForbiddenOwner = fmtErr("you are not the owner of this event")
	ErrInvalidStatus  = fmtErr("invalid status")
)

type fmtErr string

func (e fmtErr) Error() string { return string(e) }

// Error bisnis dengan pesan dinamis
func ErrBusiness(msg string) error { return fmtErr(msg) }

func (s *eventSvc) UpdateStatus(eventID int, status, vendorWallet string) error {
	if status != "ended" && status != "minting" {
		return ErrInvalidStatus
	}
	vendorID, err := s.repo.FindVendorIDByWallet(vendorWallet)
	if err != nil {
		return err
	}
	ownerID, dbStatus, st, en, err := s.repo.GetEventOwnerStatusDates(eventID)
	if err != nil {
		return err
	}
	if ownerID != vendorID {
		return ErrForbiddenOwner
	}

	calculated := timeutil.HitungStatus(dbStatus, st, en)
	if status == "ended" && calculated != "minting" {
		return ErrBusiness("Event cannot be marked as 'ended' yet. It is currently " + calculated)
	}
	if status == "minting" && dbStatus != "ended" {
		return ErrBusiness("Event is not 'ended', so it cannot be changed back to 'minting'")
	}
	return s.repo.UpdateStatus(eventID, status)
}

func (s *eventSvc) MarkAttendance(token, wallet string) error {
	id, sdb, st, en, err := s.repo.GetEventIDByToken(token)
	if err != nil {
		return err
	}
	if timeutil.HitungStatus(sdb, st, en) != "ongoing" {
		return ErrBusiness("event not ongoing")
	}
	uid, err := s.repo.GetUserIDByWallet(wallet)
	if err != nil {
		return ErrBusiness("user not found")
	}
	wstat, err := s.repo.GetWhitelistStatus(id, uid)
	if err != nil {
		return ErrBusiness("not whitelisted")
	}
	if wstat != "approved" {
		return ErrBusiness("whitelist not approved")
	}
	return s.repo.UpsertAttendance(id, uid, token)
}

func (s *eventSvc) CreateWhitelist(eventID int, wallet string) (string, string, int, time.Time, time.Time, error) {
	uid, err := s.repo.GetUserIDByWallet(wallet)
	if err != nil {
		return "", "", 0, time.Time{}, time.Time{}, ErrBusiness("wallet not registered")
	}
	max, err := s.repo.GetMaxAttendees(eventID)
	if err != nil {
		return "", "", 0, time.Time{}, time.Time{}, err
	}
	appCount, err := s.repo.CountApprovedWhitelist(eventID)
	if err != nil {
		return "", "", 0, time.Time{}, time.Time{}, err
	}

	status := "pending"
	msg := "Whitelist successful, but quota is full. You are on the waiting list."
	if appCount < max {
		status = "approved"
		msg = "Whitelist successful! You are registered as an event participant."
	}
	id, c, u, err := s.repo.InsertWhitelist(eventID, uid, wallet, status)
	if err != nil {
		// unique / duplicate key
		if err.Error() != "" && (strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate")) {
			return "", "", 0, time.Time{}, time.Time{}, ErrBusiness("already registered")
		}
		return "", "", 0, time.Time{}, time.Time{}, err
	}

	return status, msg, id, c, u, nil
}

func (s *eventSvc) CancelWhitelist(eventID int, wallet string) error {
	uid, err := s.repo.GetUserIDByWallet(wallet)
	if err != nil {
		return ErrBusiness("user not found")
	}
	aff, err := s.repo.DeleteWhitelist(eventID, uid)
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrBusiness("whitelist entry not found")
	}
	return nil
}

func (s *eventSvc) AttendanceByEvent(eventID int) ([]map[string]interface{}, error) {
	rows, err := s.repo.GetAttendanceByEvent(eventID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		m := map[string]interface{}{"user_id": r.UserID, "name": r.Name, "wallet_address": r.WalletAddress, "attend_status": r.AttendStatus}
		if r.AttendedAt != nil {
			m["attended_at"] = r.AttendedAt
		}
		out = append(out, m)
	}
	return out, nil
}

func (s *eventSvc) WhitelistByEvent(eventID int) ([]map[string]interface{}, error) {
	rows, err := s.repo.GetWhitelistByEvent(eventID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]interface{}{
			"user_id": r.UserID, "name": r.Name, "email": r.Email, "wallet_address": r.WalletAddress, "status": r.Status, "created_at": r.CreatedAt,
		})
	}
	return out, nil
}

func (s *eventSvc) GetUserByWallet(wallet string) (*model.User, error) {
	id, name, email, c, u, err := s.repo.GetUserByWallet(wallet)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Name: name, Email: email, WalletAddress: wallet, CreatedAt: c, UpdatedAt: u}, nil
}

func (s *eventSvc) GetUserEvents(wallet, baseURL string) ([]model.Event, error) {
	uid, err := s.repo.GetUserIDByWallet(wallet)
	if err != nil {
		return nil, ErrBusiness("user not found")
	}
	rows, err := s.repo.GetEventsByUser(uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Event
	for rows.Next() {
		var e model.Event
		var dbStatus, picturePath string
		var start, end time.Time
		var userStatus *string
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VendorID, &start, &end, &dbStatus, &e.CreatedAt, &e.UpdatedAt, &picturePath, &e.MaxAttendees, &e.Location, &e.Requirements, &e.Agenda, &e.Attendees, &userStatus); err != nil {
			return nil, err
		}
		e.StartDate = start
		e.EndDate = end
		e.Status = timeutil.HitungStatus(dbStatus, start, end)
		e.Picture = baseURL + "/" + picturePath
		// userStatus bisa dipakai kalau ingin expose status user khusus
		out = append(out, e)
	}
	return out, nil
}

func (s *eventSvc) GetUserCertificates(wallet, baseURL string) ([]model.CertificateWithEvent, error) {
	uid, err := s.repo.GetUserIDByWallet(wallet)
	if err != nil {
		return nil, ErrBusiness("user not found")
	}
	rows, err := s.repo.GetCertificatesByUser(uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.CertificateWithEvent
	for rows.Next() {
		var c model.CertificateWithEvent
		var pictureRaw string
		if err := rows.Scan(&c.ID, &c.EventID, &c.UserID, &c.CertificateData, &c.MintStatus, &c.MintTransactionHash, &c.CreatedAt, &c.UpdatedAt, &c.EventTitle, &c.EventDescription, &c.EventStartDate, &c.EventLocation, &pictureRaw); err != nil {
			return nil, err
		}
		if pictureRaw != "" {
			c.EventPicture = baseURL + "/" + pictureRaw
		}
		out = append(out, c)
	}
	return out, nil
}

func (s *eventSvc) GetVendorEvents(wallet, baseURL string) ([]model.Event, error) {
	vid, err := s.repo.GetVendorIDByWallet(wallet)
	if err != nil {
		return nil, ErrBusiness("vendor not found")
	}
	rows, err := s.repo.GetEventsByVendor(vid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Event
	for rows.Next() {
		var e model.Event
		var dbStatus, picturePath string
		var start, end time.Time
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VendorID, &start, &end, &dbStatus, &e.CreatedAt, &e.UpdatedAt, &picturePath, &e.MaxAttendees, &e.Location, &e.Requirements, &e.Agenda, &e.Token, &e.Attendees); err != nil {
			return nil, err
		}
		e.StartDate = start
		e.EndDate = end
		e.Status = timeutil.HitungStatus(dbStatus, start, end)
		e.Picture = baseURL + "/" + picturePath
		out = append(out, e)
	}
	return out, nil
}

func (s *eventSvc) GetCertificateAll() ([]model.CertificateWithEvent, error) {
	rows, err := s.repo.GetAllCertificates()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.CertificateWithEvent
	for rows.Next() {
		var c model.CertificateWithEvent
		if err := rows.Scan(&c.ID, &c.EventID, &c.UserID, &c.CertificateData, &c.MintStatus, &c.MintTransactionHash, &c.CreatedAt, &c.UpdatedAt, &c.UrlMetadata, &c.UrlCertificate, &c.CertificateType); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}
