package models

type TeamEntity struct {
	TeamID   uint   `json:"team_id"`
	Team     Team   `json:"team" gorm:"constraint:OnDelete:CASCADE;"`
	EntityID uint   `json:"entity_id" gorm:"primaryKey;autoIncrement:false"`
	Entity   Entity `json:"entity" gorm:"constraint:OnDelete:CASCADE;"`
}
