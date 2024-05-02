package database

import (
	"fmt"
	"gin/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbName string) {
	// Connect to MySQL server (without specifying a database)
	dsn := ":@tcp(localhost:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to MySQL server")
	}

	// Create the database if it doesn't exist
	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))

	// Close the initial connection
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get SQL database from GORM")
	}
	sqlDB.Close()

	// Connect to the new database
	dsn = ":@tcp(localhost:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// AutoMigrate will ONLY create tables, missing columns and missing indexes
	err = db.AutoMigrate(&entities.PlayerStats{})
	if err != nil {
		panic("Failed to migrate database")
	}

	DB = db
}
