package models

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/webcrawler?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	DB = database

	err = DB.AutoMigrate(&URL{}, &BrokenLink{})
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	log.Println("Database connection successful!")
}
