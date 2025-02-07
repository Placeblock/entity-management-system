package models

import "time"

type Token struct {
	EntityID  uint   `gorm:"primaryKey"`
	Pin       string `gorm:"unique"`
	CreatedAt time.Time
}
