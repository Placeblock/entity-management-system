package models

import "time"

type Token struct {
	EntityID  uint      `gorm:"primaryKey" json:"entity_id"`
	Entity    Entity    `json:"entity,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Pin       string    `gorm:"unique" json:"pin"`
	CreatedAt time.Time `json:"-"`
}
