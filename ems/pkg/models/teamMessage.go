package models

type TeamMessage struct {
	ID       uint   `json:"id"`
	MemberID uint   `json:"member_id" gorm:"primaryKey;autoIncrement:false"`
	Member   Member `json:"member"`
	Message  string `json:"message"`
}
