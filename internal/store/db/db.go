package db

import (
	"errors"
	"fmt"
	"goth/internal/store"
	"os"

	//"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func open(dbName string) (*gorm.DB, error) {

	password :=  os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	server := os.Getenv("DB")
	//port := os.Getenv("DB_PORT")

	if dbName == "" {
		dbName = "goth.db"
	}

	
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, server, dbName)

	config := &gorm.Config{}

	if os.Getenv("env") == "production" {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to open database"))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	// make the temp directory if it doesn't exist
	err = os.MkdirAll("/tmp", 0755)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MustOpen(dbName string) *gorm.DB {
	db, err := open(dbName)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&store.User{}, &store.Session{}, &store.Poll{}, &store.PollOption{}, &store.PollVote{})

	if err != nil {
		panic(err)
	}

	return db
}
