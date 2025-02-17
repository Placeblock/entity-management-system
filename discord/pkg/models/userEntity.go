package models

import "github.com/Placeblock/nostalgicraft-ems/pkg/models"

type UserEntity struct {
	UserID   string        ``
	EntityID uint          ``
	Entity   models.Entity `gorm:"constraint:OnDelete:CASCADE;"`
}
