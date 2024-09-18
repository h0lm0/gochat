package database

import (
	dbmodels "gochat/database/models"
	"gochat/keydb"
	"gochat/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Connect() *gorm.DB {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := "host=db user=" + username + " password=" + password + " dbname=" + dbName + " port=5432 sslmode=disable TimeZone=Europe/Paris"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil
	}
	log.Println("Successfully connected to database")
	return db
}

func Init() {
	Db = Connect()
	Db.AutoMigrate(&dbmodels.User{}, &dbmodels.Ip{}, &dbmodels.Message{}, &dbmodels.Room{})
	CreateAdminUser()
	CreateInitialRooms()
}

// Creation functions
func CreateAdminUser() {
	systemUser := dbmodels.User{
		Username: "System",
		Role:     models.RoleSystem,
		Color:    31,
	}
	createResult := Db.Create(&systemUser)
	if createResult.Error != nil {
		panic(createResult.Error.Error())
	}
	keydb.AddEncryptionKey(systemUser)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), 14)
	Db.Create(&dbmodels.User{
		Username: os.Getenv("ADMIN_USERNAME"),
		Password: string(hashedPassword),
		Role:     models.RoleAdmin,
	})
}

func CreateInitialRooms() {
	Db.Create(&dbmodels.Room{
		Name: "public",
		Type: 0,
	})
	Db.Create(&dbmodels.Room{
		Name: "private",
		Type: 1,
	})
}
