package db

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() (*gorm.DB, error) {

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	conn, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Migrate(conn *gorm.DB) error {
	return conn.Transaction(func(tsx *gorm.DB) error {

		err := tsx.AutoMigrate(
			&User{},
			&VPS{},
			&Request{},
		)
		if err != nil {
			return err
		}

		return nil
	})
}

/* a lot of these methods are very trivial - could just call db methods
 * directly in the client code
 */

func UserRegister(conn *gorm.DB, user *User) (uint, error) {

	err := conn.Create(user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// returns user id on successful auth
func UserCheckCreds(conn *gorm.DB, email string, password string) (uint, error) {

	loginUser := User{}
	err := conn.
		Select("id", "password").
		Where("email = ?", email).
		First(&loginUser).
		Error
	if err != nil {
		return 0, err
	}

	// maybe move the bcrypt stuff into it's own function
	err = bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(password))
	if err != nil {
		return 0, err
	}

	return loginUser.ID, nil
}

// check if username or email are already taken (true if not avaliable - possibly bad design)
func UserCheckRegistered(conn *gorm.DB, email string) (bool, error) {

	matches := []User{}
	err := conn.
		Where("email = ?", email).
		Find(&matches).
		Error
	if err != nil {
		return true, err
	}

	return len(matches) != 0, nil
}

func UserOwnsVPS(conn *gorm.DB, userID uint, vpsID uint) (bool, error) {

	matches := []VPS{}
	err := conn.
		Where("id = ? AND user_id = ?", vpsID, userID).
		Find(&matches).
		Error
	if err != nil {
		return false, err
	}

	return len(matches) != 0, nil
}

func VPSGetUserAll(conn *gorm.DB, userID uint) ([]VPS, error) {

	allVPS := []VPS{}
	err := conn.Where("user_id = ?", userID).Find(&allVPS).Error
	if err != nil {
		return nil, err
	}

	return allVPS, nil
}

func VPSGet(conn *gorm.DB, vpsID uint) (VPS, error) {

	vpsInfo := VPS{}
	err := conn.Where("id = ?", vpsID).Find(&vpsInfo).Error
	if err != nil {
		return VPS{}, err
	}

	return vpsInfo, nil
}

func VPSCreate(conn *gorm.DB, newVPS VPS) error {

	err := conn.Create(&newVPS).Error
	return err
}

func VPSDestroy(conn *gorm.DB, vpsID uint) error {

	err := conn.Delete(&VPS{}, vpsID).Error
	return err
}

func RequestListWithPurpose(conn *gorm.DB, purpose uint) ([]Request, error) {

	request := []Request{}
	err := conn.Preload("User").Where("request_purpose = ?", purpose).Find(&request).Error

	return request, err
}

func RequestListUser(conn *gorm.DB, userID uint) ([]Request, error) {

	requests := []Request{}
	err := conn.Where("user_id = ?", userID).Find(&requests).Error

	return requests, err
}

func RequestByID(conn *gorm.DB, requestID uint) (Request, error) {
	request := Request{}
	err := conn.Where("id = ?", requestID).First(&request).Error

	return request, err
}

// did the given user create the request
func RequestUserOwns(conn *gorm.DB, userID uint, requestID uint) (bool, error) {

	request := []Request{}
	err := conn.Where("id = ? AND user_id = ?", requestID, userID).Error
	if err != nil {
		return false, err
	}

	return len(request) != 0, nil
}

func RequestCreate(conn *gorm.DB, newRequest Request) error {
	return conn.Create(&newRequest).Error
}

func RequestDelete(conn *gorm.DB, requestID uint) error {
	return conn.Delete(&Request{}, requestID).Error
}

func RequestTruncate(conn *gorm.DB) error {
	// gorm will not execute batch delete without a condition
	// TODO this is retarded
	return conn.Where("1 = 1").Delete(&Request{}).Error
}
