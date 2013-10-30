package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	// "reflect"
)

func main() {

	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	sql := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' ORDER BY table_name;"

	rows, err := db.Query(sql)
	var name string
	for rows.Next() {
		rows.Scan(&name)
		fmt.Println(name)
	}

}
