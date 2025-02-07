package models

type Hue float32

type Team struct {
	ID      uint   `json:"id"`
	OwnerId uint   `json:"owner_id"`
	Name    string `json:"name"`
	Hue     Hue    `json:"hue"`
}
