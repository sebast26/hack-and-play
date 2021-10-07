package cmd

import (
	"github.com/jmoiron/sqlx"
	"goclass/internal/client"
	"goclass/internal/repository"
	"goclass/internal/service"
	"log"
)

func getService() service.Service {
	db := getDB()

	c := client.New()
	r := repository.New(db)
	s := service.New(c, r)

	return s
}

func getDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("error opening the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	return db
}
