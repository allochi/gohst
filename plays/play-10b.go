package main

import (
	_ "github.com/lib/pq"
	// "github.com/davecgh/go-spew/spew"
	"database/sql"
	"log"
)

// Benchmark 10000
// 7.281s play-09
// 7.154s play-10 JSON (98.25%)
// 6.890s play-10 FLAT (94.62%)
// 4.262s play-10 PREPARED JSON (58.53%)

func main() {

	db, err := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	if err != nil {
		return
	}
	defer db.Close()

	// select * from json_bogs where id in (select unnest(string_to_array('1,2,3,4', ',')::integer[]))
	stmt, err := db.Prepare("select id, data from json_bogs where id in (select unnest(string_to_array($1, ',')::integer[]))")
	ids := "1,2,3,4"
	// ids := []int64{1, 2, 3, 4}
	rows, err := stmt.Query(ids)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var id int
		var data string
		err = rows.Scan(&id, &data)
		log.Printf("%d: %s", id, data)
	}
}
