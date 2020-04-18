package db

import (
	"database/sql"
	"fmt"
	"log"
)

// test connection
func Test() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		// do something here
	}

	fmt.Println("end of func db test")

	defer db.Close()
}
