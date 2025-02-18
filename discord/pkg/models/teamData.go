package models

type TeamData struct {
	TeamID    uint   `gorm:"primaryKey"`
	RoleID    string `gorm:"unique"`
	ChannelID string `gorm:"unique"`
}
