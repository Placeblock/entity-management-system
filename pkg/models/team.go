package models

type Hue float32

type Team struct {
	ID      uint     `json:"id"`
	OwnerID uint     `json:"owner_id,omitempty" gorm:"unique"`
	Owner   Entity   `json:"owner,omitempty"`
	Name    string   `json:"name,omitempty" gorm:"unique"`
	Hue     *Hue     `json:"hue,omitempty"`
	Members []Member `json:"members,omitempty"`
}
