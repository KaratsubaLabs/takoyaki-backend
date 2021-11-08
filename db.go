package main

import (
    "fmt"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type PostgresConnection struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
}

func StartDB() {

	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		ConnectionSecrets.Host,
		ConnectionSecrets.Port,
		ConnectionSecrets.User,
		ConnectionSecrets.Password,
		ConnectionSecrets.DBName,
	)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
    _ = err

	db.AutoMigrate(
		&User{},
		&VPS{},
		&Request{},
	)

}

func RegisterUser(tsx *gorm.DB) {

}

func GetVPS(tsx *gorm.DB) {

}

func CreateVPS(tsx *gorm.DB) {

}

func DestroyVPS(tsx *gorm.DB) {

}

