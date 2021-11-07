package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	REQUEST_PURPOSE_ACCOUNT = 0
	REQUEST_PURPOSE_VPSCREATE = 1
)

type User struct {
	gorm.Model
	Username       string
	Password       string
}

type VPS struct {
	gorm.Model
	Name           string
	UserID         uint
	// VPSConfig      VPSConfig
	CreationTime   time.Time
}

type Request struct {
	gorm.Model
	RequestTime    time.Time
	// RequestPurpose
}

