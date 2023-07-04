package db

import (
	"fmt"

	"github.com/flunks-nft/discord-bot/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID   uint
	Name string
	Age  int
}

func init() {
	utils.LoadEnv()

}

func main() {
	// Replace the connection string with your actual PostgreSQL connection details
	dsn := "host=localhost user=yourusername password=yourpassword dbname=yourdbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate the User struct to create the table if it doesn't exist
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to migrate table")
	}

	// Create a new user
	user := User{Name: "John Doe", Age: 30}
	db.Create(&user)

	// Query all users
	var users []User
	db.Find(&users)
	fmt.Println(users)

	// Update the user
	db.Model(&user).Update("Age", 31)

	// Delete the user
	db.Delete(&user)
}
