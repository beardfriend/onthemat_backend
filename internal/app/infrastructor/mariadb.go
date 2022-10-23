package infrastructor

import (
	"context"
	"log"

	"onthemat/pkg/ent"

	_ "github.com/go-sql-driver/mysql"
)

func NewMariaDB() {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	c.MariaDB.User, c.MariaDB.Password, c.MariaDB.Host, c.MariaDB.Port, c.MariaDB.Database)

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("Could not connect to database :%v", err)
	// }
	// return db

	client, err := ent.Open("mysql", "root:password@tcp(localhost:3306)/db?parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
