package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	REQUEST_PURPOSE_REGISTER = 0
	REQUEST_PURPOSE_VPS_CREATE = 1
	REQUEST_PURPOSE_VPS_UPGRADE = 2
)

type User struct {
	gorm.Model
	Username       string
	Password       string
    Email          string
}

type VPS struct {
	gorm.Model
	Name           string
    UserID         uint
    User           User        `gorm:"foreignKey:UserID;preload:false"`
	// VPSConfig      VPSConfig
	CreationTime   time.Time
}

type Request struct {
	gorm.Model
    UserID         uint        `gorm:"foreignKey:UserID;preload:false"`
	RequestTime    time.Time
	// RequestPurpose
    Message        string
}

