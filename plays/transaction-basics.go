package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Film struct {
	Id          int64
	Title       string
	ReleaseDate time.Time
}

var postgres *sql.DB

func init() {

	postgres, _ = sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")

	createTable := `
	DROP TABLE IF EXISTS films;
	CREATE TABLE IF NOT EXISTS films (
		id        			serial PRIMARY KEY,
		title       		text NOT NULL,
		release_date   	date
		);
	`

	postgres.Exec(createTable)
}

func main() {

	film_01 := Film{}
	film_01.Title = "Star Trek"
	film_01.ReleaseDate, _ = time.Parse("2006-01-02 15:04:05", "2009-01-01 00:00:00")

	tx, _ := postgres.Begin()
	var i int
	for i = 0; i < 10; i++ {
		_, err := tx.Exec("INSERT INTO films (title, release_date) VALUES ($1,$2)", film_01.Title, "2009-01-01 00:00:00")
		if err != nil {
			panic(err)
		}
	}

	if i%2 == 0 {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	rows, err := postgres.Query("SELECT * FROM films;")
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var film Film
		rows.Scan(&film.Id, &film.Title, &film.ReleaseDate)
		fmt.Printf("%#v\n", film)
	}

}
