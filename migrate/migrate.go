package main

import (
	"2024_2_FIGHT-CLUB/ds"
	"2024_2_FIGHT-CLUB/dsn"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func migrate() (err error) {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&ds.User{}, &ds.Ad{}, &ds.Request{}, &ds.Review{})
	if err != nil {
		return err
	}
	fmt.Println("Database migrated")
	return nil
}

func main() {
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}
