package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type EventListRow struct {
	ID, VendorID, MaxAttendees, Whitelisted                        int
	Title, Description, DBStatus, PicturePath, Location, Organizer string
	StartDate, EndDate, CreatedAt, UpdatedAt                       time.Time
}

type EventDetailRow struct {
	ID, VendorID, MaxAttendees, Attendees, Whitelisted, Minted int
	Title, Description, DBStatus, PicturePath, Location, Token string
	StartDate, EndDate, CreatedAt, UpdatedAt                   time.Time
	Organizer                                                  *string
	Requirements, Agenda                                       json.RawMessage
	UrlCertificate                                             *string
}

type AttendanceUserRow struct {
	UserID                            int
	Name, WalletAddress, AttendStatus string
	AttendedAt                        *time.Time
}

type WhitelistUserRow struct {
	UserID                             int
	Name, Email, WalletAddress, Status string
	CreatedAt                          time.Time
}

type EventRepo interface {
	List(ctx context.Context) ([]EventListRow, error)
	Detail(ctx context.Context, id int) (*EventDetailRow, error)
	Insert(ctx context.Context, vendorID int, title, desc string, start, end time.Time, picture string, max int, location string, req, agenda json.RawMessage, token string) (int, time.Time, time.Time, error)
	FindVendorIDByWallet(wallet string) (int, error)
	Cancel(eventID int) error
	GetEventOwnerStatusDates(eventID int) (vendorID int, dbStatus string, start, end time.Time, err error)
	UpdateStatus(eventID int, status string) error
	GetEventIDByToken(token string) (id int, status string, start, end time.Time, err error)
	GetUserIDByWallet(wallet string) (int, error)
	GetWhitelistStatus(eventID, userID int) (string, error)
	UpsertAttendance(eventID, userID int, token string) error
	CountApprovedWhitelist(eventID int) (int, error)
	GetMaxAttendees(eventID int) (int, error)
	InsertWhitelist(eventID, userID int, wallet, status string) (int, time.Time, time.Time, error)
	DeleteWhitelist(eventID, userID int) (int64, error)
	GetAttendanceByEvent(eventID int) ([]AttendanceUserRow, error)
	GetWhitelistByEvent(eventID int) ([]WhitelistUserRow, error)
	GetUserByWallet(wallet string) (id int, name, email string, createdAt, updatedAt time.Time, err error)
	GetEventsByUser(userID int) (*sql.Rows, error)
	GetCertificatesByUser(userID int) (*sql.Rows, error)
	GetVendorIDByWallet(wallet string) (int, error)
	GetEventsByVendor(vendorID int) (*sql.Rows, error)
	GetAllCertificates() (*sql.Rows, error)
}

type eventRepo struct{ db *sql.DB }

func NewEventRepo(db *sql.DB) EventRepo { return &eventRepo{db: db} }

func (r *eventRepo) List(ctx context.Context) ([]EventListRow, error) {
	rows, err := r.db.QueryContext(ctx, `
    SELECT e.id,e.title,e.description,e.vendor_id,e.start_date,e.end_date,
           e.status,e.created_at,e.updated_at,e.picture,e.maxattendees,e.location,
           v.vendor_name,
           (SELECT COUNT(*) FROM whitelist w WHERE w.event_id=e.id AND w.status='approved') as whitelisted
    FROM events e
    LEFT JOIN vendors v ON v.id=e.vendor_id
    GROUP BY e.id, v.vendor_name
    ORDER BY e.start_date ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EventListRow
	for rows.Next() {
		var x EventListRow
		if err := rows.Scan(&x.ID, &x.Title, &x.Description, &x.VendorID, &x.StartDate, &x.EndDate, &x.DBStatus, &x.CreatedAt, &x.UpdatedAt, &x.PicturePath, &x.MaxAttendees, &x.Location, &x.Organizer, &x.Whitelisted); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *eventRepo) Detail(ctx context.Context, id int) (*EventDetailRow, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT e.id,v.vendor_name,e.title,e.description,e.vendor_id,e.start_date,e.end_date,e.status,
		       e.created_at,e.updated_at,e.picture,e.maxattendees,e.location,
		       COALESCE(att.attendees,0), COALESCE(wl.whitelisted,0), COALESCE(mi.minted,0),
		       e.requirements,e.agenda,e.token,COALESCE(ec.url_certificate,'')
		FROM events e
		LEFT JOIN vendors v ON e.vendor_id=v.id
		LEFT JOIN (SELECT event_id, COUNT(*) attendees FROM attendance WHERE attendance_status='present' GROUP BY event_id) att ON e.id=att.event_id
		LEFT JOIN (SELECT event_id, COUNT(*) whitelisted FROM whitelist WHERE status='approved' GROUP BY event_id) wl ON e.id=wl.event_id
		LEFT JOIN (SELECT event_id, COUNT(*) minted FROM certificates WHERE mint_status='success' GROUP BY event_id) mi ON e.id=mi.event_id
		LEFT JOIN event_certificates ec ON e.id=ec.event_id
		WHERE e.id=$1`, id)
	var x EventDetailRow
	if err := row.Scan(&x.ID, &x.Organizer, &x.Title, &x.Description, &x.VendorID, &x.StartDate, &x.EndDate, &x.DBStatus, &x.CreatedAt, &x.UpdatedAt, &x.PicturePath, &x.MaxAttendees, &x.Location, &x.Attendees, &x.Whitelisted, &x.Minted, &x.Requirements, &x.Agenda, &x.Token, &x.UrlCertificate); err != nil {
		return nil, err
	}
	return &x, nil
}

func (r *eventRepo) Insert(ctx context.Context, vendorID int, title, desc string, start, end time.Time, picture string, max int, location string, req, agenda json.RawMessage, token string) (int, time.Time, time.Time, error) {
	var id int
	var created, updated time.Time
	err := r.db.QueryRowContext(ctx, `
	  INSERT INTO events (title,description,vendor_id,start_date,end_date,picture,maxattendees,location,requirements,agenda,token)
	  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	  RETURNING id, created_at, updated_at`,
		title, desc, vendorID, start, end, picture, max, location, req, agenda, token).
		Scan(&id, &created, &updated)
	return id, created, updated, err
}

func (r *eventRepo) FindVendorIDByWallet(wallet string) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM vendors WHERE wallet_address=$1`, wallet).Scan(&id)
	return id, err
}

func (r *eventRepo) Cancel(eventID int) error {
	_, err := r.db.Exec(`UPDATE events SET status='canceled' WHERE id=$1`, eventID)
	return err
}

func (r *eventRepo) GetEventOwnerStatusDates(eventID int) (int, string, time.Time, time.Time, error) {
	var vid int
	var s string
	var st, en time.Time
	err := r.db.QueryRow(`SELECT vendor_id,status,start_date,end_date FROM events WHERE id=$1`, eventID).Scan(&vid, &s, &st, &en)
	return vid, s, st, en, err
}

func (r *eventRepo) UpdateStatus(eventID int, status string) error {
	_, err := r.db.Exec(`UPDATE events SET status=$1 WHERE id=$2`, status, eventID)
	return err
}

func (r *eventRepo) GetEventIDByToken(token string) (int, string, time.Time, time.Time, error) {
	var id int
	var s string
	var st, en time.Time
	err := r.db.QueryRow(`SELECT id,status,start_date,end_date FROM events WHERE token=$1`, token).Scan(&id, &s, &st, &en)
	return id, s, st, en, err
}

func (r *eventRepo) GetUserIDByWallet(wallet string) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM users WHERE wallet_address=$1`, wallet).Scan(&id)
	return id, err
}

func (r *eventRepo) GetWhitelistStatus(eventID, userID int) (string, error) {
	var status string
	err := r.db.QueryRow(`SELECT status FROM whitelist WHERE event_id=$1 AND user_id=$2`, eventID, userID).Scan(&status)
	return status, err
}

func (r *eventRepo) UpsertAttendance(eventID, userID int, token string) error {
	_, err := r.db.Exec(`
        INSERT INTO attendance (event_id,user_id,attendance_status,token_input)
        VALUES ($1,$2,'present',$3)
        ON CONFLICT (event_id,user_id)
        DO UPDATE SET attendance_status='present', token_input=$3, updated_at=NOW()`,
		eventID, userID, token)
	return err
}

func (r *eventRepo) CountApprovedWhitelist(eventID int) (int, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM whitelist WHERE event_id=$1 AND status='approved'`, eventID).Scan(&n)
	return n, err
}

func (r *eventRepo) GetMaxAttendees(eventID int) (int, error) {
	var m int
	err := r.db.QueryRow(`SELECT maxattendees FROM events WHERE id=$1`, eventID).Scan(&m)
	return m, err
}

func (r *eventRepo) InsertWhitelist(eventID, userID int, wallet, status string) (int, time.Time, time.Time, error) {
	var id int
	var c, u time.Time
	err := r.db.QueryRow(`
		INSERT INTO whitelist (event_id,user_id,wallet_address,status)
		VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at`,
		eventID, userID, wallet, status).Scan(&id, &c, &u)
	return id, c, u, err
}

func (r *eventRepo) DeleteWhitelist(eventID, userID int) (int64, error) {
	res, err := r.db.Exec(`DELETE FROM whitelist WHERE event_id=$1 AND user_id=$2`, eventID, userID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *eventRepo) GetAttendanceByEvent(eventID int) ([]AttendanceUserRow, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.name, u.wallet_address,
		       COALESCE(a.attendance_status,'absent') as attend_status,
		       a.created_at
		FROM whitelist w
		JOIN users u ON w.user_id=u.id
		LEFT JOIN attendance a ON a.event_id=w.event_id AND a.user_id=w.user_id AND a.attendance_status='present'
		WHERE w.event_id=$1
		ORDER BY u.name ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AttendanceUserRow
	for rows.Next() {
		var x AttendanceUserRow
		var at sql.NullTime
		if err := rows.Scan(&x.UserID, &x.Name, &x.WalletAddress, &x.AttendStatus, &at); err != nil {
			return nil, err
		}
		if at.Valid {
			t := at.Time
			x.AttendedAt = &t
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *eventRepo) GetWhitelistByEvent(eventID int) ([]WhitelistUserRow, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.name, u.email, u.wallet_address, w.status, w.created_at
		FROM users u JOIN whitelist w ON u.id=w.user_id
		WHERE w.event_id=$1`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []WhitelistUserRow
	for rows.Next() {
		var x WhitelistUserRow
		if err := rows.Scan(&x.UserID, &x.Name, &x.Email, &x.WalletAddress, &x.Status, &x.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *eventRepo) GetUserByWallet(wallet string) (int, string, string, time.Time, time.Time, error) {
	var id int
	var n, e string
	var c, u time.Time
	err := r.db.QueryRow(`SELECT id,name,email,created_at,updated_at FROM users WHERE wallet_address=$1`, wallet).Scan(&id, &n, &e, &c, &u)
	return id, n, e, c, u, err
}

func (r *eventRepo) GetEventsByUser(userID int) (*sql.Rows, error) {
	return r.db.Query(`
		SELECT e.id,e.title,e.description,e.vendor_id,e.start_date,e.end_date,e.status,e.created_at,e.updated_at,
		       e.picture,e.maxattendees,e.location,e.requirements,e.agenda,
		       COALESCE(attend_count.attendees,0) as attendees,
		       CASE WHEN att.attendance_status='present' THEN 'present' ELSE wl.status END as user_status
		FROM events e
		LEFT JOIN attendance att ON e.id=att.event_id AND att.user_id=$1
		LEFT JOIN whitelist wl ON e.id=wl.event_id AND wl.user_id=$1
		LEFT JOIN (SELECT event_id, COUNT(*) attendees FROM attendance WHERE attendance_status='present' GROUP BY event_id) attend_count
			ON e.id=attend_count.event_id
		WHERE att.user_id IS NOT NULL OR wl.user_id IS NOT NULL
		ORDER BY e.start_date ASC`, userID)
}

func (r *eventRepo) GetCertificatesByUser(userID int) (*sql.Rows, error) {
	return r.db.Query(`
		SELECT c.id,c.event_id,c.user_id,c.certificate_data,c.mint_status,c.mint_transaction_hash,c.created_at,c.updated_at,
		       e.title,e.description,e.start_date,e.location,e.picture
		FROM certificates c
		JOIN events e ON c.event_id=e.id
		WHERE c.user_id=$1
		ORDER BY c.created_at DESC`, userID)
}

func (r *eventRepo) GetVendorIDByWallet(wallet string) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM vendors WHERE wallet_address=$1`, wallet).Scan(&id)
	return id, err
}

func (r *eventRepo) GetEventsByVendor(vendorID int) (*sql.Rows, error) {
	return r.db.Query(`
		SELECT e.id,e.title,e.description,e.vendor_id,e.start_date,e.end_date,e.status,e.created_at,e.updated_at,
		       e.picture,e.maxattendees,e.location,e.requirements,e.agenda,e.token,
		       COALESCE(present_count.attendees,0) as attendees
		FROM events e
		LEFT JOIN (SELECT event_id, COUNT(id) attendees FROM attendance WHERE attendance_status='present' GROUP BY event_id) present_count
			ON e.id=present_count.event_id
		WHERE e.vendor_id=$1
		ORDER BY e.start_date ASC`, vendorID)
}

func (r *eventRepo) GetAllCertificates() (*sql.Rows, error) {
	return r.db.Query(`SELECT id,event_id,user_id,certificate_data,mint_status,mint_transaction_hash,created_at,updated_at,url_metadata,url_certificate,certificate_type FROM certificates`)
}
