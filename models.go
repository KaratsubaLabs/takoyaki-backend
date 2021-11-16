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
	Username       string      `gorm:"not null"`
	Password       string      `gorm:"not null"`
	Email          string      // maybe remove this?
}

type VPS struct {
	gorm.Model
	InternalName   string      `gorm:"not null"` // randomly generated string that is used to by libvirt
	UserID         uint        `gorm:"not null"`
    User           User        `gorm:"foreignKey:UserID;preload:false"`
	// VPSConfig      VPSConfig
	CreationTime   time.Time
}

type Request struct {
	gorm.Model
    UserID         uint        `gorm:"foreignKey:UserID;preload:false"`
	RequestTime    time.Time   `gorm:"not null"`
	RequestPurpose int         `gorm:"not null"` // this could be enum instead
	RequestData    string      `gorm:"not null;default:'{}'::JSONB"`
    Message        string
}

// stored as json in db
type VPSConfig struct {
    DisplayName   string
    Hostname      string
    Username      string
    Password      string
    SSHKey        string
    RAM           int // make this 'enum' or sm
    CPU           int
    Disk          int // in gb
    OS            string
}

