package infrastructor

import (
	"fmt"
	"log"

	"onthemat/internal/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMariaDB() *gorm.DB {
	config := config.Info

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MariaDB.User, config.MariaDB.Password, config.MariaDB.Host, config.MariaDB.Port, config.MariaDB.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database :%v", err)
	}
	return db
}
