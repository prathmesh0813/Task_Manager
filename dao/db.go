package dao

import (
	"os"
	"task_manager/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User DB Schema
type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"not null"`
	MobileNo int64  `gorm:"not null"`
	Gender   string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
}

// Token DB Schema
type Token struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	RefreshToken string    `gorm:"not null;unique"`
	UserToken    string    `gorm:"not null;unique"`
	Timestamp    time.Time `gorm:"not null"`
	UserID       int64
	User         User `gorm:"foreignKey:UserID"`
}

// Login DB schema
type Login struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	UserID   int64
	User     User `gorm:"foreignKey:UserID"`
}

func InitDB() {
	var err error



	DB, err = gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		utils.Logger.Fatal("Could not connect to database", zap.Error(err))
	}

	createTables()
}

func createTables() {
	err := DB.AutoMigrate(&User{}, &Login{}, &Token{})
	if err != nil {
		utils.Logger.Fatal("could not migrate tables", zap.Error(err))
	}
}
