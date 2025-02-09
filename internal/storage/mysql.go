package storage

import (
	"fmt"
	"log"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	user := "nostalgicraft"
	password := "GVIdC4CDvg49GD8h"
	fmt.Println(user, password)
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/nostalgicraft?charset=utf8mb4&parseTime=True", user, password)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database!")
	db.AutoMigrate(&models.Entity{}, &models.Token{}, &models.Team{}, &models.Member{}, &models.MemberInvite{})
	return db
}
