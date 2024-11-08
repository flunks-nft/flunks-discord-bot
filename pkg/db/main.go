package db

import (
	"fmt"
	"os"

	"github.com/flunks-nft/discord-bot/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB_HOST                  string
	DB_USER                  string
	DB_PASSWORD              string
	DB_NAME                  string
	DB_PORT                  string
	CLOUD_DB_CONNECTION_NAME string
	db                       *gorm.DB // Declare the db variable at the package level
)

func init() {
	utils.LoadEnv()

	DB_HOST = os.Getenv("DB_HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	DB_PORT = os.Getenv("DB_PORT")
	CLOUD_DB_CONNECTION_NAME = os.Getenv("CLOUD_DB_CONNECTION_NAME")
}

func GetDB() *gorm.DB {
	return db
}

func InitDB() {
	// PostgreSQL connection details
	var dsn string
	if os.Getenv("ENV") == "production" {
		dsn = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Singapore", CLOUD_DB_CONNECTION_NAME, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Singapore", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	}

	dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db = dbConnection

	// Auto Migrate the structs
	err = db.AutoMigrate(&User{}, &Raid{}, &Nft{}, &Trait{})
	if err != nil {
		panic("Failed to migrate table")
	}
}
