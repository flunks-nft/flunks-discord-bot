package db

import (
	"fmt"
	"os"

	"github.com/flunks-nft/discord-bot/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string
	db          *gorm.DB // Declare the db variable at the package level
)

func init() {
	utils.LoadEnv()

	DB_HOST = os.Getenv("DB_HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	DB_PORT = os.Getenv("DB_PORT")
}

func InitDB() {
	// PostgreSQL connection details
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Singapore", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	fmt.Println(dsn)

	dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db = dbConnection

	// Auto Migrate the structs
	err = db.AutoMigrate(&User{}, &Raid{})
	if err != nil {
		panic("Failed to migrate table")
	}
}
