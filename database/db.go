package database

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/miscOS/ddns-bridge/models"
)

var db *gorm.DB

func Init() {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	var err error

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" {
		db, err = gorm.Open(sqlite.Open("./data/db.sqlite"), &gorm.Config{})
	} else {
		dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ")/" + dbName + "?parseTime=true"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	}

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Webhook{}, &models.Task{})
}

func GetDB() *gorm.DB {
	return db
}
