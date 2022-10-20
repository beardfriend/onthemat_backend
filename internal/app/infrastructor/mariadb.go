package infrastructor

import (
	"fmt"
	"log"

	"onthemat/internal/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMariaDB(c *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MariaDB.User, c.MariaDB.Password, c.MariaDB.Host, c.MariaDB.Port, c.MariaDB.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database :%v", err)
	}
	return db
}
