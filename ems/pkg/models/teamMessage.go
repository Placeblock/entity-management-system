package models

type TeamMessage struct {
	ID       uint   `json:"id"`
	MemberID uint   `json:"member_id" gorm:"primaryKey;autoIncrement:false"`
	Member   Member `json:"member"`
	TeamID   uint   `json:"team_id"`
	Team     Team   `json:"team,omitempty"`
	Message  string `json:"message"`
}
