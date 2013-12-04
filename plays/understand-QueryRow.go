package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {

	db, _ := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	data := ""
	row := db.QueryRow("SELECT data FROM json_contacts where id = $1;", 10)
	err := row.Scan(&data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[data] %s\n", data)

}
