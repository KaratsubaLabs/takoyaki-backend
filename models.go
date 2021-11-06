package main

import (
	"time"

	"github.com/jinzhu/gorm"
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

