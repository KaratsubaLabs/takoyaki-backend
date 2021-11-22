package main

import (
    "fmt"
    "os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"golang.org/x/crypto/bcrypt"
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

/* a lot of these methods are very trivial - could just call db methods
 * directly in the client code
 */

func DBUserRegister(db *gorm.DB, user User) (uint, error) {

	err := db.Select("id").Create(&user).Error
	if err != nil { return 0, err }

	return user.ID, nil
}

// returns user id on successful auth
func DBUserCheckCreds(db *gorm.DB, username string, password string) (uint, error) {

	loginUser := User{}
	err := db.
		Select("id", "password").
		Where("username = ?", username).
		First(&loginUser).
		Error
	if err != nil { return 0, err }

	// maybe move the bcrypt stuff into it's own function
	err = bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(password))
	if err != nil { return 0, err }

	return loginUser.ID, nil
}

// check if username or email are already taken (true if not avaliable - possibly bad design)
func DBUserCheckRegistered(db *gorm.DB, username string, email string) (bool, error) {

	matches := []User{}
	err := db.
		Where("username = ?", username).
		Or("email = ?", email).
		Find(&matches).
		Error
	if err != nil { return true, err }

	return len(matches) == 0, nil
}

func DBUserOwnsVPS(db *gorm.DB, userID uint, vpsID uint) (bool, error) {

	matches := []VPS{}
	err := db.
		Where("id = ? AND user_id = ?", vpsID, userID).
		Find(&matches).
		Error
	if err != nil { return false, err }

	return len(matches) != 0, nil
}

func DBVPSGetInfo(db *gorm.DB, userID uint) ([]VPS, error) {

	allVPS := []VPS{}
	err := db.Where("user_id = ?", userID).Find(&allVPS).Error
	if err != nil { return nil, err }

	return allVPS, nil
}

func DBVPSCreate(db *gorm.DB) {

}

func DBVPSDestroy(db *gorm.DB) {

}

func DBRequestListWithPurpose(db *gorm.DB, purpose uint) ([]Request, error) {

	request := []Request{}
	err := db.Preload("User").Where("request_purpose = ?", purpose).Find(&request).Error

	return request, err
}

func DBRequestListUser(db *gorm.DB, userID uint) ([]Request, error) {

	requests := []Request{}
	err := db.Where("user_id = ?", userID).Find(&requests).Error

	return requests, err
}

func DBRequestByID(db *gorm.DB, requestID uint) (Request, error) {
	request := Request{}
	err := db.Where("id = ?", requestID).First(&request).Error

	return request, err
}

// did the given user create the request
func DBRequestUserOwns(db *gorm.DB, userID uint, requestID uint) (bool, error) {

	request := []Request{}
	err := db.Where("id = ? AND user_id = ?", requestID, userID).Error
	if err != nil {
		return false, err
	}

	return len(request) != 0, nil
}

func DBRequestCreate(db *gorm.DB, newRequest Request) error {
	return db.Create(&newRequest).Error
}

func DBRequestDelete(db *gorm.DB, requestID uint) error {
	return db.Delete(&Request{}, requestID).Error
}

func DBRequestTruncate(db *gorm.DB) error {
	// gorm will not execute batch delete without a condition
	return db.Where("1 = 1").Delete(&Request{}).Error
}

