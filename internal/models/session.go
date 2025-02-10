package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	ClientID string `gorm:"uniqueIndex"`
	Session  []byte // Session data
}
