package models

import "time"

type Token struct {
	EntityID  uint      `gorm:"primaryKey" json:"entityId"` //TODO: ENTITY REFERENCE
	Entity    Entity    `json:"entity" gorm:"constraint:OnDelete:CASCADE;"`
	Pin       string    `gorm:"unique" json:"pin"`
	CreatedAt time.Time `json:"-"`
}
