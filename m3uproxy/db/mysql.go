package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var db *sql.DB
var dbGorm *gorm.DB

type User struct {
	ID             int
	Username       string
	Password       string
	Status         string
	ExpDate        time.Time
	IsTrial        bool
	CreatedAt      time.Time
	MaxConnections int
}

// init database
func Init() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		fmt.Println(err.Error())
	} else {

	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Database opened successfully")

	createDatabase()
	createTable()
}

func createDatabase() {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS m3uproxy")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created database..")
	}

	_, err = db.Exec("USE m3uproxy")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("DB selected successfully..")
	}
}

func createTable() {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS `users` ( `id` INT NOT NULL AUTO_INCREMENT , `username` VARCHAR(50) NOT NULL , `password` VARCHAR(50) NOT NULL , `status` VARCHAR(20) NOT NULL , `exp_date` TIMESTAMP NOT NULL , `is_trial` BOOLEAN NOT NULL , `created_at` TIMESTAMP NOT NULL , `max_connections` INT NOT NULL , PRIMARY KEY (`id`), INDEX (`username`), INDEX (`password`), INDEX (`exp_date`)) ENGINE = InnoDB;")

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}

	dbGorm, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/m3uproxy?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
}

func CreateUser(u User) {
	dbGorm.NewRecord(u)
	dbGorm.Create(&u)
}

func GetUser(username string, password string) (user User) {
	var u User
	dbGorm.Where("username = ? AND password = ?", username, password).Find(&u)
	//fmt.Println(u)
	return u
}

func Close() {
	defer db.Close()
}
