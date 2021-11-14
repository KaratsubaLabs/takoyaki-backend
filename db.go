package main

import (
    "fmt"
    "os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func DBConnection() (*gorm.DB, error) {

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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
	return db.Transaction(func(tsx *gorm.DB) error {

		err := tsx.AutoMigrate(
			&User{},
			&VPS{},
			&Request{},
		)
		if err != nil { return err }

		return nil
	})
}

func DBUserRegister(db *gorm.DB, user User) error {
	return db.Transaction(func(tsx *gorm.DB) error {

		err := tsx.Create(&user).Error
		if err != nil { return err }

		return nil
	})
}

func DBUserCheckCreds(db *gorm.DB) {

}

func DBVPSGetInfo(db *gorm.DB) {

}

func DBVPSCreate(db *gorm.DB) {

}

func DBVPSDestroy(db *gorm.DB) {

}

