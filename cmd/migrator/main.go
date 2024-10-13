package main

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/module/dsn"
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
	err = db.AutoMigrate(&domain.User{}, &domain.Ad{}, &domain.Request{}, &domain.Review{})
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
