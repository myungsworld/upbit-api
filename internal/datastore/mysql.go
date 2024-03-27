package datastore

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"upbit-api/internal/models"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@(%s)/%s?parseTime=True",
		"root",
		"",
		"127.0.0.1",
		"upbit",
	)

	dbConfig := &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	}

	//fmt.Println(dsn)

	DB, err = gorm.Open(mysql.Open(dsn), dbConfig)

	if err != nil {
		panic(err)
	}

	updateDB()

}

func updateDB() {
	tables := []interface{}{
		(*models.AutoTrading2)(nil),
	}

	if err := DB.AutoMigrate(tables...); err != nil {
		panic(err)
	}
}
