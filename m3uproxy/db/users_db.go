package db

import (
	"database/sql"
	"log"
)

//test connection
func test() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		// do something here
	}

	defer db.Close()
}
