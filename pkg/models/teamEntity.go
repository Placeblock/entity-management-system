package models

type TeamEntity struct {
	TeamID   uint   `json:"team_id"`
	Team     Team   `json:"team" gorm:"constraint:OnDelete:CASCADE;"`
	EntityID uint   `json:"entity_id" gorm:"primaryKey;autoIncrement:false"`
	Entity   Entity `json:"entity" gorm:"constraint:OnDelete:CASCADE;"`
}

type TeamEntityInvite struct {
	InvitedID uint   `json:"invited_id" gorm:"primaryKey;autoIncrement:false"`
	Invited   Entity `json:"invited" gorm:"constraint:OnDelete:CASCADE;"`
	InviterID uint   `json:"inviter_id"`
	Inviter   Entity `json:"inviter" gorm:"constraint:OnDelete:CASCADE;"`
	TeamID    uint   `json:"team_id" gorm:"primaryKey;autoIncrement:false"`
	Team      Team   `json:"team" gorm:"constraint:OnDelete:CASCADE;"`
}
