package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	database := os.Getenv("DB_NAME")
	fmt.Println(user, password)
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True", user, password, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database!")
	db.AutoMigrate(&models.Entity{}, &models.Token{}, &models.Team{}, &models.Member{}, &models.MemberInvite{})
	return db
}
