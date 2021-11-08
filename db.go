package main

import (
    "fmt"
    "os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func DBConnection() (*gorm.DB, error) {

	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
    if err != nil { return nil, err }

    return db, nil

}

func DBMigrate(db *gorm.DB) error {

    tsx := db.Begin()
    defer tsx.RollbackUnlessCommitted()

    err := tsx.AutoMigrate(
		&User{},
		&VPS{},
		&Request{},
	)
    if err != nil { return err }

    return tsx.Commit().Error

}

func RegisterUser(db *gorm.DB) {

}

func GetVPS(db *gorm.DB) {

}

func CreateVPS(db *gorm.DB) {

}

func DestroyVPS(db *gorm.DB) {

}

