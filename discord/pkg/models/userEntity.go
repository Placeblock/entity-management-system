package models

type UserEntity struct {
	UserID   string ``
	EntityID uint   ``
	Entity   Entity `gorm:"constraint:OnDelete:CASCADE;"`
}
