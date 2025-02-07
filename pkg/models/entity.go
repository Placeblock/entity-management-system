package models

type Entity struct {
	ID   uint
	Name string `gorm:"unique"`
}
