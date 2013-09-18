package main

import (
	. "allochi/gohst/plays/models"
	"encoding/json"
	_ "github.com/lib/pq"
	// "github.com/davecgh/go-spew/spew"
	"bytes"
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

	stmt, err := db.Prepare("INSERT INTO json_bogs (data, created_at, updated_at) VALUES ($1,NOW(),NOW()) RETURNING id")

	for i := 0; i < 10000; i++ {
		var bog Bog
		bog.Name = "Allochi"
		bog.Messages = []string{"This is the first bog", "If another bog is created then it will be bogbog", "I don''t know if bogs are OK, but later we will have complex object tested"}

		data, err := json.Marshal(bog)
		if err != nil {
			log.Fatalln(err)
		}
		data = bytes.Replace(data, []byte("'"), []byte("''"), -1)

		// data := []byte("{}")

		var id int64
		// err = db.QueryRow("INSERT INTO json_bogs (data, created_at, updated_at) VALUES ($1,NOW(),NOW()) RETURNING id", data).Scan(&id)
		err = stmt.QueryRow(data).Scan(&id)
		if err != nil {
			log.Fatalln(err)
		}

	}
}
