package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func main() {

	db, _ := sql.Open("postgres", "user=allochi dbname=allochi_contactizer sslmode=disable")
	timer := time.Now()
	sql := fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s ORDER BY %s) row_data;", "json_contacts", "(data->>'country_id')::int DESC")
	var result string
	db.QueryRow(sql).Scan(&result)
	duration := time.Since(timer).Nanoseconds()
	fmt.Printf("Query in %vs!\n", float64(duration)/float64(1000000000))

}
