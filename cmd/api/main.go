package main

import (
	"log"

	"api-certify-nft-be/config"
	"api-certify-nft-be/internal/app"
	"api-certify-nft-be/pkg/db"
)

func main() {
	cfg := config.Load()
	pg, err := db.Open(cfg.DBDsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	e := app.NewServer(cfg, pg)
	log.Printf("running on :%s ...", cfg.AppPort)
	if err := e.Start(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
