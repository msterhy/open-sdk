package dbx

import (
	"github.com/trancecho/open-sdk/database"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	DB = database.GetDb("MainMysql")
	if DB == nil {
		log.Fatalln("failed to connect database")
	}
	err := AutoMigrate(DB)
	if err != nil {
		log.Fatalln(err)
	}
}
