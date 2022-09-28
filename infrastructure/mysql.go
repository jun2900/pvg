package infrastructure

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	_  = godotenv.Load(".pvg.test")
	DB *gorm.DB
)

func OpenDbConnection() *gorm.DB {

	//dialect := os.Getenv("DB_DIALECT")   //utility.KVGet("DB_DIALECT")
	username := os.Getenv("DB_USER")     //utility.KVGet("DB_USER")
	password := os.Getenv("DB_PASSWORD") //utility.KVGet("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")       //utility.KVGet("DB_NAME")
	host := os.Getenv("DB_HOST")         //utility.KVGet("DB_HOST")
	port := os.Getenv("DB_PORT")         //utility.KVGet("DB_PORT")
	var db *gorm.DB
	var err error

	// db, err := gorm.Open("mysql", "root:root@localhost/go_api_shop_gonc?charset=utf8")
	//databaseUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable ", host, username, password, dbName)
	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)

	//fmt.Println(databaseURL)
	db, err = gorm.Open(mysql.Open(databaseURL), &gorm.Config{})

	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
		//fmt.Println("db err: ", err)
		os.Exit(-1)
	}

	//db.DB().SetMaxIdleConns(10)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	// db.LogMode(true)
	DB = db

	return DB
}
