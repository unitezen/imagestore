package core

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDatabase() {
	var err error
	dsn := "host=database user=postgres password=testdatabase dbname=imagestore port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
