package models

type Hue float32

type Team struct {
	ID      uint     `json:"id"`
	Name    string   `json:"name" gorm:"unique"`
	Hue     *Hue     `json:"hue"`
	Members []Member `json:"-"`
}
