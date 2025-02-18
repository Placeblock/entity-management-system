package models

type TeamMessage struct {
	ID       uint   `json:"id"`
	MemberID uint   `json:"member_id" gorm:"autoIncrement:false;"`
	Member   Member `json:"member" gorm:"references:id;constraint:OnDelete:CASCADE;"`
	Message  string `json:"message"`
}
