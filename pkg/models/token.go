package models

import "time"

type Token struct {
	EntityID  uint      `gorm:"primaryKey" json:"entityId"` //TODO: ENTITY REFERENCE
	Pin       string    `gorm:"unique" json:"pin"`
	CreatedAt time.Time `json:"-"`
}
