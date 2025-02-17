package models

import "github.com/Placeblock/nostalgicraft-ems/pkg/models"

type UserEntity struct {
	EntityID uint          `gorm:"primaryKey"`
	Entity   models.Entity `gorm:"constraint:OnDelete:CASCADE;"`
	UserID   string        `gorm:"unique"`
}
