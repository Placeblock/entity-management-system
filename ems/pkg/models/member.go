package models

type Member struct {
	ID       uint   `json:"id" gorm:"unique;autoIncrement:true"`
	TeamID   uint   `json:"team_id"`
	Team     Team   `json:"team,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	EntityID uint   `json:"entity_id" gorm:"primaryKey;autoIncrement:false"`
	Entity   Entity `json:"entity,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
}

type MemberInvite struct {
	ID        uint   `json:"id" gorm:"unique;autoIncrement:true"`
	InvitedID uint   `json:"invited_id" gorm:"primaryKey;autoIncrement:false"`
	Invited   Entity `json:"invited,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	InviterID uint   `json:"inviter_id" gorm:"unique;autoIncrement:false;"`
	Inviter   Member `json:"inviter,omitempty" gorm:"constraint:OnDelete:CASCADE;references:id"`
}
