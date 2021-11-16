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
	Username       string      `gorm:"not null;unique"`
	Email          string      `gorm:"not null;unique"`
	Password       string      `gorm:"not null"`
}

type VPS struct {
	gorm.Model
	DisplayName    string      `gorm:"not null"` // the name the user assigned
	InternalName   string      `gorm:"not null"` // randomly generated string that is used to by libvirt
	UserID         uint        `gorm:"not null"`
    User           User        `gorm:"foreignKey:UserID;preload:false"`
	CreationTime   time.Time
	RAM            int         `gorm:"not null"`
	CPU            int         `gorm:"not null"`
	Disk           int         `gorm:"not null"`
	OS             string      `gorm:"not null"`
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
type VPSCreateRequestData struct {
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

