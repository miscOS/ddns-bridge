package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/miscOS/ddns-bridge/models"
)

var db *gorm.DB

func Init() {

	db_user := "db_user"
	db_password := "db_pass"
	db_name := "test_db"
	db_host := "10.0.11.171"
	db_port := "3306"

	dsn := db_user + ":" + db_password + "@tcp(" + db_host + ":" + db_port + ")/" + db_name + "?parseTime=true"

	/*TODO:
	- Use sqlite if no db_host is provided
	*/

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})
}

func GetDB() *gorm.DB {
	return db
}
