package main

import (
	dbModule "github.com/notegio/openrelay/db"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"io/ioutil"
	"fmt"
)

func main() {
	pgHost := os.Args[1]
	pgUser := os.Args[2]
	pgPassword := ""
	for _, arg := range os.Args[3:] {
		pgPasswordFile := arg
		pgPasswordBytes, err := ioutil.ReadFile(pgPasswordFile)
		if err != nil {
			log.Fatalf("Could not read password file: %v", err.Error());
		}
		pgPassword = string(pgPasswordBytes)
	}
	if pgPassword == "" {
		pgPassword = os.Getenv("POSTGRES_PASSWORD")
	}
	connectionString := fmt.Sprintf(
		"host=%v sslmode=disable user=%v password=%v",
		pgHost,
		pgUser,
		pgPassword,
	)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Could not open postgres connection: %v", err.Error())
	}
	if err := db.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		log.Fatalf("Error migrating database: %v", err.Error())
	}
}
