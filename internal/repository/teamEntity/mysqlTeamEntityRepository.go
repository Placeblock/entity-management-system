package teamentity

import "gorm.io/gorm"

type MysqlTeamEntityRepository struct {
	db *gorm.DB
}

func NewMysqlTeamEntityRepository(db *gorm.DB) *MysqlTeamEntityRepository {
	return &MysqlTeamEntityRepository{db}
}
