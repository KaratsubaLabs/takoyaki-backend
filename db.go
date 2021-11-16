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

func DBUserRegister(db *gorm.DB, user User) (uint, error) {

	err := db.Select("id").Create(&user).Error
	if err != nil { return 0, err }

	return user.ID, nil
}

// returns user id on successful auth
func DBUserCheckCreds(db *gorm.DB, username string, password string) (uint, error) {

	loginUser := User{}
	err := db.
		Select("id").
		Where("username = ? AND password = ?", username, password).
		First(&loginUser).
		Error
	if err != nil { return 0, err }

	return loginUser.ID, nil
}

func DBVPSGetInfo(db *gorm.DB, userID uint) ([]*VPS, error) {

	allVPS := []*VPS{}
	err := db.Where("user_id = ?", userID).Find(&allVPS).Error
	if err != nil { return nil, err }

	return allVPS, nil
}

func DBVPSCreate(db *gorm.DB) {

}

func DBVPSDestroy(db *gorm.DB) {

}

func DBRequestCreate(db *gorm.DB, newRequest Request) error {
	return db.Create(&newRequest).Error
}

func DBRequestDelete(db *gorm.DB) error {
	return nil
}

