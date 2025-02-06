package dao

import (
	"task_manager/logger"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User DB Schema
type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"not null"`
	MobileNo string `gorm:"not null"`
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

// Avatar DB schema
type Avatar struct {
	ID     int64  `gorm:"primaryKey;autoIncrement"`
	Data   []byte `gorm:"type:blob;not null"`
	Name   string `gorm:"not null"`
	UserID int64
	User   User `gorm:"foreignKey:UserID"`
}

// Task DB schema
type Task struct {
	ID          int64  `json:"id"`
	Title       string ` json:"title"`
	Description string `json:"description"`
	Completed   string `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      int64 ` json:"userId"`
}

func InitDB() {
	var err error

	// DB, err = gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})
	// if err != nil {
	//
	// 	logger.Error("requestID", "could not connect to database", err.Error())
	// }

	DB, err = gorm.Open(sqlite.Open("task.db"), &gorm.Config{})
	if err != nil {
		logger.Error("requestid", "Could not connect to database", err.Error())
	}

	createTables()
}

func createTables() {
	err := DB.AutoMigrate(&User{}, &Login{}, &Token{}, &Avatar{}, &Task{})
	if err != nil {
		logger.Error("requestID", "could not migrate tables", err.Error())
	}
}
