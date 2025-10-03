package timeutil

import "time"

func HitungStatus(dbStatus string, mulai, selesai time.Time) string {
	if dbStatus == "ended" || dbStatus == "canceled" {
		return dbStatus
	}
	now := time.Now()
	if now.After(selesai) {
		return "minting"
	}
	if now.After(mulai) && now.Before(selesai) {
		return "ongoing"
	}
	return "upcoming"
}
