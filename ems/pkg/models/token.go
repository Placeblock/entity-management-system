package models

import "time"

type Token struct {
	EntityID  uint      `gorm:"primaryKey" json:"entityId"`
	Entity    Entity    `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Pin       string    `gorm:"unique" json:"pin"`
	CreatedAt time.Time `json:"-"`
}
