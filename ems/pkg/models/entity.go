package models

type Entity struct {
	ID     uint    `json:"id"`
	Name   string  `gorm:"unique" json:"name"`
	Member *Member `json:"-"`
}
