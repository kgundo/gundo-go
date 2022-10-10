package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() *gorm.DB {
	if Db == nil {
		Db = connectDB()
		mysqlDbInstance, err := Db.DB()
		if mysqlDbInstance != nil && err == nil {
			fmt.Println("Database connection setting")
			mysqlDbInstance.SetMaxIdleConns(35)
			mysqlDbInstance.SetMaxOpenConns(50)
			mysqlDbInstance.SetConnMaxLifetime(15 * time.Minute)
		}
	}
	return Db
}

func connectDB() *gorm.DB {
	dbHost := viper.Get("DB_HOST").(string)
	dbPort := viper.Get("DB_PORT").(string)
	dbName := viper.Get("DB_NAME").(string)
	dbUser := viper.Get("DB_USER").(string)
	dbPassword := viper.Get("DB_PASSWORD").(string)
	dsn := dbUser + ":" + dbPassword + "@tcp" + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil
	}
	return db
}
