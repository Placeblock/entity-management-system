package models

type Hue float32

type Team struct {
	ID      uint   `json:"id"`
	OwnerID uint   `json:"owner_id" gorm:"unique"`
	Owner   Entity `json:"owner"`
	Name    string `json:"name" gorm:"unique"`
	Hue     Hue    `json:"hue"`
}
